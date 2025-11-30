package infra

import (
	"context"
	"time"

	"github.com/DangeL187/erax"
	"github.com/IBM/sarama"

	"consumer/internal/features/consumer/handler"
	"consumer/internal/shared/config"
)

type KafkaConsumer struct {
	cfg *config.Config

	cg sarama.ConsumerGroup
}

func (kc *KafkaConsumer) Run(ctx context.Context, msgHandler func(message *sarama.ConsumerMessage)) error {
	h := handler.NewMessageHandler(msgHandler)

	err := kc.cg.Consume(ctx, []string{kc.cfg.KafkaTopic}, h)
	if err != nil {
		return erax.Wrap(err, "failed to consume message from Kafka")
	}

	return nil
}

func (kc *KafkaConsumer) Stop() error {
	err := kc.cg.Close()
	if err != nil {
		return erax.Wrap(err, "failed to close kafka consumer group")
	}

	return nil
}

func NewKafkaConsumer(cfg *config.Config) (*KafkaConsumer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.V4_0_0_0
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	kafkaConfig.Consumer.Offsets.AutoCommit.Enable = true
	kafkaConfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	cg, err := sarama.NewConsumerGroup([]string{cfg.KafkaBroker}, cfg.KafkaGroupID, kafkaConfig)
	if err != nil {
		return nil, erax.Wrap(err, "failed to create kafka consumer group")
	}

	return &KafkaConsumer{
		cfg: cfg,
		cg:  cg,
	}, nil
}
