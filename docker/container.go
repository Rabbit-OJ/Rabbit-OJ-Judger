package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"io"
	"os"
)

func ForceContainerRemove(ID string) {
	fmt.Printf("[Docker] will force remove container %s", ID)
	if err := Client.ContainerRemove(Context, ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	}); err != nil {
		fmt.Printf("[Docker] Error when force removing %s container, %+v \n", ID, err)
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
			fmt.Println(err)
			return
		}
		defer func() { _ = out.Close() }()

		if _, err := io.Copy(os.Stderr, out); err != nil {
			fmt.Println(err)
		}
	}()
}
