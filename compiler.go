package judger

import (
	"errors"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/config"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/docker"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/logger"
	JudgerModels "github.com/Rabbit-OJ/Rabbit-OJ-Judger/models"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func Compiler(sid uint32, codePath string, code []byte, compileInfo *JudgerModels.CompileInfo) ([]byte, error) {
	vmPath := codePath + "vm/"
	logger.Printf("(%d) [Compile] Start %s \n", sid, codePath)

	logger.Printf("(%d) [Compile] Touched empty output file for build \n", sid)
	containerConfig := &container.Config{
		Entrypoint:      compileInfo.BuildArgs,
		Tty:             true,
		OpenStdin:       true,
		Image:           compileInfo.BuildImage,
		NetworkDisabled: true,
	}

	containerHostConfig := &container.HostConfig{
		Resources: container.Resources{
			NanoCPUs: compileInfo.Constraints.CPU,
			Memory:   compileInfo.Constraints.Memory,
		},
	}

	if config.Global.Extensions.HostBind {
		containerHostConfig.Binds = []string{
			utils.DockerHostConfigBinds(vmPath, path.Dir(compileInfo.BuildTarget)),
		}
	}

	if config.Global.AutoRemove.Containers && config.Global.Extensions.HostBind {
		containerHostConfig.AutoRemove = true
	}

	logger.Printf("(%d) [Compile] Creating container \n", sid)
	resp, err := docker.Client.ContainerCreate(docker.Context,
		containerConfig,
		containerHostConfig,
		nil,
		"")

	if err != nil {
		return nil, err
	}

	if config.Global.AutoRemove.Containers && !config.Global.Extensions.HostBind {
		defer func() {
			go docker.ForceContainerRemove(resp.ID)
		}()
	}

	logger.Printf("(%d) [Compile] Copying files to container \n", sid)
	io, err := utils.ConvertToTar(
		[]utils.TarFileBasicInfo{
			{path.Base(compileInfo.Source), code, 0644},
		},
	)
	if err != nil {
		return nil, err
	}

	if err := docker.Client.CopyToContainer(
		docker.Context,
		resp.ID,
		path.Dir(compileInfo.Source),
		io,
		types.CopyToContainerOptions{
			AllowOverwriteDirWithFile: true,
			CopyUIDGID:                false,
		}); err != nil {
		return nil, err
	}

	logger.Printf("(%d) [Compile] Running container \n", sid)
	if err := docker.Client.ContainerStart(docker.Context, resp.ID, types.ContainerStartOptions{}); err != nil {
		logger.Printf("(%d) %+v \n", sid, err)
		return nil, err
	}

	docker.ContainerErrToStdErr(resp.ID)
	statusCh, errCh := docker.Client.ContainerWait(docker.Context, resp.ID, container.WaitConditionNotRunning)
	logger.Printf("(%d) [Compile] Waiting for status \n", sid)

	var compileProductionFiles []byte
	select {
	case err := <-errCh:
		return nil, err
	case status := <-statusCh:
		if status.StatusCode != int64(0) {
			return nil, errors.New("CE")
		}

		if !config.Global.Extensions.HostBind {
			logger.Printf("(%d) [Compile] Copying build production \n", sid)
			reader, _, err := docker.Client.CopyFromContainer(docker.Context, resp.ID, path.Dir(compileInfo.BuildTarget))
			if err != nil {
				return nil, err
			}
			defer func() {
				_ = reader.Close()
			}()

			compileProductionFiles, err = ioutil.ReadAll(reader)
			if err != nil {
				return nil, err
			}

			// todo: check if have magic number of build prod
		} else {
			if err := checkBuildResult(vmPath + path.Base(compileInfo.BuildTarget)); err != nil {
				return nil, err
			}
		}

		logger.Printf("(%d) %+v \n", sid, status)
		break
	case <-time.After(time.Duration(compileInfo.Constraints.BuildTimeout) * time.Second):
		go docker.ForceContainerRemove(resp.ID)
		return nil, errors.New("compile timeout")
	}

	return compileProductionFiles, nil
}

func checkBuildResult(path string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}

	if file.Size() <= int64(len(utils.MagicBytes)) {
		return errors.New("compile file invalid")
	}
	return nil
}
