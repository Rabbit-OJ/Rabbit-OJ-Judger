package mq

import (
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/logger"
	"github.com/Shopify/sarama"
)

func PublishMessageAsync(topic string, key, value []byte) {
	mqMessage := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	AsyncProducer.Input() <- mqMessage
}

func PublishMessageSync(topic string, key, value []byte) error {
	mqMessage := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	if _, _, err := SyncProducer.SendMessage(mqMessage); err != nil {
		logger.Println("[MQ] sync send error ", err)
		return err
	}

	return nil
}
