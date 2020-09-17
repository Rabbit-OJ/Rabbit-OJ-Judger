package mq

import (
	"context"
	"fmt"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/config"
	"log"

	"github.com/Shopify/sarama"
)

func CreateJudgeRequestConsumer(topics []string, group string) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = Version
	consumer := JudgeRequestConsumer{
		ready: make(chan bool, 0),
	}

	client, err := sarama.NewConsumerGroup(config.Global.Kafka.Brokers, group, saramaConfig)
	go func() {
		select {
		case <-CancelCtx.Done():
			_ = client.Close()
		}
	}()

	if err != nil {
		log.Panicf("Error when creating consumer group: %v", err)
		return
	}

	ctx, _ := context.WithCancel(CancelCtx)
	go func() {
		for {
			fmt.Println("[MQ] topic: request consumer group running", group)

			if err := client.Consume(ctx, topics, &consumer); err != nil {
				fmt.Printf("Error from consumer consume: %+v \n", ctx.Err())
				return
			}

			if ctx.Err() != nil {
				fmt.Printf("Error from ctx: %+v \n", ctx.Err())
				return
			}

			consumer.ready = make(chan bool, 0)
		}
	}()
}

func CreateJudgeResponseConsumer(topics []string, group string) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = Version
	consumer := JudgeResponseConsumer{
		ready: make(chan bool, 0),
	}

	client, err := sarama.NewConsumerGroup(config.Global.Kafka.Brokers, group, saramaConfig)
	if err != nil {
		log.Panicf("Error when creating consumer group: %v", err)
		return
	}

	ctx, _ := context.WithCancel(CancelCtx)
	go func() {
		for {
			fmt.Println("[MQ] topic: response consumer group running", group)

			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			if ctx.Err() != nil {
				log.Panicf("Error from ctx: %v", ctx.Err())
				return
			}

			consumer.ready = make(chan bool, 0)
		}
	}()
}
