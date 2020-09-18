package judger

import (
	"context"
	"fmt"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/docker"
	JudgerModels "github.com/Rabbit-OJ/Rabbit-OJ-Judger/models"
	"sync"
	"testing"
)

var  (
	alreadyInit = false
	initMu sync.Mutex
)

func initJudger() {
	initMu.Lock()
	defer initMu.Unlock()
	if alreadyInit {
		return
	}

	alreadyInit = true
	ctx, _ := context.WithCancel(context.Background())
	cfg := &JudgerModels.JudgerConfigType{
		Kafka: JudgerModels.KafkaConfig{
			Brokers: []string{
				"localhost:9092",
			},
		},
		Rpc: "",
		AutoRemove: JudgerModels.AutoRemoveType{
			Containers: true,
			Files:      true,
		},
		Concurrent: JudgerModels.ConcurrentType{
			Judge: 2,
		},
		LocalImages: []string{
			"alpine_tester:latest",
		},
		Languages: []JudgerModels.LanguageType{
			{
				ID:      "cpp17",
				Name:    "C++17",
				Enabled: true,
				Args: JudgerModels.CompileInfo{
					BuildArgs: []string{
						"g++",
						"-std=c++17",
						"/home/code.cpp",
						"-Wall",
						"-lm",
						"-fno-asm",
						"--static",
						"-O2",
						"-o",
						"/home/code.o",
					},
					Source:      "/home/code.cpp",
					NoBuild:     false,
					BuildTarget: "/home/code.o",
					BuildImage:  "gcc:10.2.0",
					Constraints: JudgerModels.Constraints{
						CPU:          1000000000,
						Memory:       1073741824,
						BuildTimeout: 120,
						RunTimeout:   120,
					},
					RunArgs:     []string{"/home/code.o"},
					RunArgsJSON: "[\"/home/code.o\"]",
					RunImage:    "alpine_tester:latest",
				},
			},
		},
		Extensions: JudgerModels.ExtensionsType{
			HostBind: false,
			AutoPull: true,
			CheckJudge: JudgerModels.CheckJudgeType{
				Enabled:  false,
				Interval: 0,
				Requeue:  false,
			},
			Expire: JudgerModels.ExpireType{
				Enabled:  false,
				Interval: 0,
			},
		},
	}

	InitJudger(ctx, cfg, func(tid uint32, version string) ([]*JudgerModels.TestCaseType, error) {
		return make([]*JudgerModels.TestCaseType, 0), nil
	}, true, false, "Judge")

	OnJudgeResponse = append(OnJudgeResponse, func(sid uint32, isContest bool, judgeResult []*JudgerModels.JudgeResult) {
		fmt.Println(sid, isContest, judgeResult)
	})
}

func TestInitJudger(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v \n", err)
			t.Fail()
		}
	}()

	initJudger()
	needImages := docker.GetNeedImages()
	for image, need := range needImages {
		if need {
			fmt.Printf("Need to image %s", image)
			t.Fail()
		}
	}
}

//func TestScheduler(t *testing.T) {
//	defer func() {
//		if err := recover(); err != nil {
//			fmt.Printf("%+v \n", err)
//			t.Fail()
//		}
//	}()
//
//	initJudger()
//}