package docker

import (
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/logger"
	"github.com/docker/docker/api/types"
	"io"
	"os"
)

func ForceContainerRemove(ID string) {
	logger.Printf("[Docker] will force remove container %s \n", ID)
	if err := Client.ContainerRemove(Context, ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         true,
	}); err != nil {
		logger.Printf("[Docker] Error when force removing %s container, %+v \n", ID, err)
	}
}

func ContainerErrToStdErr(ID string) {
	go func() {
		out, err := Client.ContainerLogs(Context, ID, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Follow:     true,
		})
		if err != nil {
			logger.Println(err)
			return
		}
		defer func() { _ = out.Close() }()

		if _, err := io.Copy(os.Stderr, out); err != nil {
			logger.Println(err)
		}
	}()
}
