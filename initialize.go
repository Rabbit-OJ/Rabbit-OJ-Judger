package judger

import (
	"Rabbit-OJ-Backend/services/channel"
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

	channel.MQPublishMessageChannel = make(chan *channel.MQMessage)
	if os.Getenv("Role") == "Judge" {
		channel.JudgeRequestDeliveryChan = make(chan []byte)
		channel.JudgeRequestBridgeChan = make(chan *channel.JudgeRequestBridgeMessage)

		mq.CreateJudgeRequestConsumer([]string{config.JudgeRequestTopicName}, "req1")
		go JudgeRequestHandler()
		go MachineJudgeRequestBridge()
	}

	if os.Getenv("Role") == "Server" {
		channel.JudgeResponseDeliveryChan = make(chan []byte)

		mq.CreateJudgeResponseConsumer([]string{config.JudgeResponseTopicName}, "res1")
		go JudgeResultHandler()
	}
	go mq.PublishService()
}

