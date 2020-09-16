package judger

import (
	"Rabbit-OJ-Judger/config"
	JudgerConfig "Rabbit-OJ-Judger/config"
	"Rabbit-OJ-Judger/docker"
	JudgerModels "Rabbit-OJ-Judger/models"
	"Rabbit-OJ-Judger/mq"
	"context"
	"encoding/json"
	"os"
)

type StorageInitFuncType = func(tid uint32, version string) (uint32, uint32, string, error)
var (
	StorageInitFunc StorageInitFuncType
)

func InitJudger(ctx context.Context, config *JudgerModels.JudgerConfigType, storageInitFunc StorageInitFuncType) {
	JudgerConfig.Global = config
	StorageInitFunc = storageInitFunc
	if os.Getenv("Role") == "Judger" {
		docker.InitDocker()
	}
	Language()
	MQ(ctx)
}

func MQ(ctx context.Context) {
	mq.InitKafka(ctx)

	if os.Getenv("Role") == "Judge" {
		mq.JudgeRequestDeliveryChan = make(chan []byte)
		mq.CreateJudgeRequestConsumer([]string{config.JudgeRequestTopicName}, "req1")
		go JudgeRequestHandler()
	}

	if os.Getenv("Role") == "Server" {
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

