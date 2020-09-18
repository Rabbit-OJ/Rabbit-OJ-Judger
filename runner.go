package judger

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/config"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/docker"
	JudgerModels "github.com/Rabbit-OJ/Rabbit-OJ-Judger/models"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"
)

func Runner(
	sid uint32, codePath string,
	compileInfo *JudgerModels.CompileInfo,
	testCases []*JudgerModels.TestCaseType,
	timeLimit, spaceLimit, outputPath string,
	code []byte, buildProduction []byte,
) (map[string][]byte, error) {
	vmPath := codePath + "vm/"
	fmt.Printf("(%d) [Runner] Compile OK, start run container \n", sid)

	resultFilePathInHost := filepath.Join(codePath, "result.json")
	err := utils.TouchFile(resultFilePathInHost)
	if err != nil {
		fmt.Printf("(%d) %+v \n", sid, err)
		return nil, err
	}
	fmt.Printf("(%d) [Runner] Touched empty result file for build \n", sid)

	containerConfig := &container.Config{
		Image:           compileInfo.RunImage,
		NetworkDisabled: true,
		Env: []string{
			"EXEC_COMMAND=" + compileInfo.RunArgsJSON,
			"CASE_COUNT=" + fmt.Sprintf("%d", len(testCases)),
			"TIME_LIMIT=" + timeLimit,
			"SPACE_LIMIT=" + spaceLimit,
			"Role=Tester",
		},
	}

	containerHostConfig := &container.HostConfig{}
	if config.Global.Extensions.HostBind {
		if len(testCases) > 0 {
			containerHostConfig.Mounts = []mount.Mount{
				{
					Source:   path.Dir(testCases[0].Path),
					Target:   "/case",
					ReadOnly: true,
					Type:     mount.TypeBind,
				},
			}
		}

		containerHostConfig.Binds = []string{
			utils.DockerHostConfigBinds(resultFilePathInHost, "/result/info.json"),
			utils.DockerHostConfigBinds(outputPath, "/output"),
		}
	}

	if !compileInfo.NoBuild {
		containerHostConfig.Binds = append(containerHostConfig.Binds,
			utils.DockerHostConfigBinds(vmPath, path.Dir(compileInfo.BuildTarget)))
	}

	if config.Global.AutoRemove.Containers && config.Global.Extensions.HostBind {
		containerHostConfig.AutoRemove = true
	}

	fmt.Printf("(%d) [Runner] Creating container \n", sid)
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

	if !compileInfo.NoBuild && !config.Global.Extensions.HostBind {
		fmt.Printf("(%d) [Runner] Copying build production to container \n", sid)
		io := bytes.NewReader(buildProduction)

		if err := docker.Client.CopyToContainer(
			docker.Context,
			resp.ID,
			"/",
			io,
			types.CopyToContainerOptions{
				AllowOverwriteDirWithFile: true,
				CopyUIDGID:                false,
			}); err != nil {
			return nil, err
		}
	}

	var tarInfos []utils.TarFileBasicInfo
	if !config.Global.Extensions.HostBind {
		fmt.Printf("(%d) [Runner] Preparing judge cases for container \n", sid)
		caseArchiveInfos, err := utils.TestCasesToTarArchiveInfo(testCases, "/case")
		if err != nil {
			return nil, err
		}
		tarInfos = append(tarInfos, caseArchiveInfos...)
	}

	if compileInfo.NoBuild {
		sourceTarInfo := utils.TarFileBasicInfo{
			Name: compileInfo.Source,
			Body: code,
			Mode: 0644,
		}

		tarInfos = append(tarInfos, sourceTarInfo)
	}

	if len(tarInfos) > 0 {
		fmt.Printf("(%d) [Runner] Copying files to container \n", sid)
		io, err := utils.ConvertToTar(tarInfos)
		if err != nil {
			return nil, err
		}

		if err := docker.Client.CopyToContainer(
			docker.Context,
			resp.ID,
			"/",
			io,
			types.CopyToContainerOptions{
				AllowOverwriteDirWithFile: true,
				CopyUIDGID:                false,
			}); err != nil {
			return nil, err
		}
	}

	fmt.Printf("(%d) [Runner] Running container \n", sid)
	if err := docker.Client.ContainerStart(docker.Context, resp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Printf("(%d) [Runner] %+v \n", sid, err)
		return nil, err
	}

	docker.ContainerErrToStdErr(resp.ID)
	statusCh, errCh := docker.Client.ContainerWait(docker.Context, resp.ID, container.WaitConditionNotRunning)
	fmt.Printf("(%d) [Runner] Waiting for status \n", sid)

	var collectedStdout map[string][]byte
	select {
	case err := <-errCh:
		return nil, err
	case status := <-statusCh:
		if !config.Global.Extensions.HostBind {
			if err := copyResultJsonFile(resp.ID, resultFilePathInHost); err != nil {
				return nil, err
			}
			if collectedStdout, err = copyStdoutFile(resp.ID); err != nil {
				return nil, err
			}
		}
		fmt.Printf("(%d) %+v \n", sid, status)
	case <-time.After(time.Duration(compileInfo.Constraints.RunTimeout) * time.Second):
		go docker.ForceContainerRemove(resp.ID)
		return nil, errors.New("run timeout")
	}

	return collectedStdout, nil
}

func copyResultJsonFile(containerID, resultFilePathInHost string) error {
	files, err := docker.CopyFromContainer(docker.Context, containerID, "/result/info.json")
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New(fmt.Sprintf("files.length = %d, not 1", len(files)))
	}
	configFile := files[0]
	if err := ioutil.WriteFile(resultFilePathInHost, configFile.Body, 0644); err != nil {
		return err
	}

	return nil
}

//func copyStdoutFile(containerID, outputPath string) error {
//	files, err := docker.CopyFromContainer(docker.Context, containerID, "/result/info.json")
//	if err != nil {
//		return err
//	}
//
//	for _, file := range files {
//		casePath := fmt.Sprintf("%s/%s", outputPath, file.Name)
//		if err := ioutil.WriteFile(casePath, file.Body, os.FileMode(file.Mode)); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

func copyStdoutFile(containerID string) (map[string][]byte, error) {
	files, err := docker.CopyFromContainer(docker.Context, containerID, "/output")
	if err != nil {
		return nil, err
	}

	collectedStdout := make(map[string][]byte)
	for _, file := range files {
		collectedStdout[file.Name] = file.Body
	}
	return collectedStdout, nil
}
