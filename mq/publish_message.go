package mq

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func PublishMessage(topic string, key, value []byte, async bool) error {
	mqMessage := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	if async {
		AsyncProducer.Input() <- mqMessage
	} else {
		if _, _, err := SyncProducer.SendMessage(mqMessage); err != nil {
			fmt.Println("[MQ] sync send error ", err)
			return err
		}
	}
	return nil
}
