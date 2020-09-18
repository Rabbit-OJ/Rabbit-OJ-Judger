package judger

import (
	"context"
	"encoding/json"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/config"
	JudgerConfig "github.com/Rabbit-OJ/Rabbit-OJ-Judger/config"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/docker"
	JudgerModels "github.com/Rabbit-OJ/Rabbit-OJ-Judger/models"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/mq"
)

type StorageGetFuncType = func(tid uint32, version string) ([]JudgerModels.TestCaseType, error)

var (
	StorageGetFunc StorageGetFuncType
)

func InitJudger(ctx context.Context, config *JudgerModels.JudgerConfigType, storageGetFunc StorageGetFuncType, withDocker bool, withKafka bool, role string) {
	JudgerConfig.Global = config
	JudgerConfig.Role = role
	StorageGetFunc = storageGetFunc

	Language()
	if withDocker {
		docker.InitDocker()
	}

	if withKafka {
		MQ(ctx)
	}
}

func MQ(ctx context.Context) {
	mq.InitKafka(ctx)

	if JudgerConfig.Role == "Judge" {
		mq.JudgeRequestDeliveryChan = make(chan []byte)
		mq.CreateJudgeRequestConsumer([]string{config.JudgeRequestTopicName}, "req1")
		go JudgeRequestHandler()
	}

	if JudgerConfig.Role == "Server" {
		mq.JudgeResponseDeliveryChan = make(chan []byte)
		mq.CreateJudgeResponseConsumer([]string{config.JudgeResponseTopicName}, "res1")
		go JudgeResultHandler()
	}
}

func Language() {
	languageCount := 0
	for _, item := range config.Global.Languages {
		if item.Enabled {
			languageCount++
		}
	}

	JudgerConfig.LocalImages = map[string]bool{}
	JudgerConfig.CompileObject = map[string]JudgerModels.CompileInfo{}
	JudgerConfig.SupportLanguage = make([]JudgerModels.SupportLanguage, languageCount)

	for _, item := range config.Global.LocalImages {
		JudgerConfig.LocalImages[item] = true
	}
	for index, item := range config.Global.Languages {
		if !item.Enabled {
			continue
		}

		JudgerConfig.SupportLanguage[index] = JudgerModels.SupportLanguage{
			Name:  item.Name,
			Value: item.ID,
		}

		runArgsJson, err := json.Marshal(item.Args.RunArgs)
		if err != nil {
			panic(err)
		}

		currentCompileObject := item.Args
		currentCompileObject.RunArgsJSON = string(runArgsJson)
		JudgerConfig.CompileObject[item.ID] = currentCompileObject
	}
}
