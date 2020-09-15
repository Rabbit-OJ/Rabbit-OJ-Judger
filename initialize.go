package judger

import (
	"Rabbit-OJ-Backend/services/config"
	"Rabbit-OJ-Backend/services/judger/docker"
	"Rabbit-OJ-Backend/services/judger/mq"
	"context"
	"os"
)

func InitJudger(ctx context.Context) {
	if os.Getenv("Role") == "Judger" {
		docker.InitDocker()
	}
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

