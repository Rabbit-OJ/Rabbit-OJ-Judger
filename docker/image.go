package docker

import (
	"fmt"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/logger"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/utils"
	"github.com/docker/docker/api/types"
	"io"
	"os"
	"strings"
)

func PullImage(tag string) {
	logger.Println("[Docker] pulling image : " + tag)
	out, err := Client.ImagePull(Context, tag, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(os.Stderr, out); err != nil {
		logger.Println(err)
	}
}

func BuildImage(tag string) {
	logger.Println("[Docker] building image from local Dockerfile : " + tag)

	name := strings.Split(tag, ":")[0]
	dockerFileBytes, err := utils.ReadFileBytes(fmt.Sprintf("./dockerfiles/%s/Dockerfile", name))
	if err != nil {
		panic(err)
	}
	serverFileBytes, err := utils.ReadFileBytes("./tester")
	if err != nil {
		panic(err)
	}
	tarBytes, err := utils.ConvertToTar([]utils.TarFileBasicInfo{
		{
			Name: "Dockerfile",
			Body: dockerFileBytes,
		},
		{
			Name: "tester",
			Body: serverFileBytes,
		},
	})
	if err != nil {
		panic(err)
	}

	resp, err := Client.ImageBuild(Context, tarBytes, types.ImageBuildOptions{
		Tags:   []string{tag},
		Remove: true,
	})
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(os.Stderr, resp.Body); err != nil {
		logger.Println(err)
	}
}
