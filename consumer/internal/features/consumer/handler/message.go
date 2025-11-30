package handler

import (
	"github.com/IBM/sarama"
)

type MessageHandler struct {
	handler func(*sarama.ConsumerMessage)
}

func (mh MessageHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (mh MessageHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (mh MessageHandler) ConsumeClaim(_ sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		mh.handler(msg)
	}
	return nil
}

func NewMessageHandler(handler func(*sarama.ConsumerMessage)) *MessageHandler {
	return &MessageHandler{handler: handler}
}
