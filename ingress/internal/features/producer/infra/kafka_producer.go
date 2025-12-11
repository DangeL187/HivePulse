package infra

import (
	"time"

	"github.com/DangeL187/erax"
	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.AsyncProducer
	errChan  chan error
}

func (kp *KafkaProducer) Produce(topic string, payload []byte) {
	kp.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
	}
}

func (kp *KafkaProducer) Close() error {
	err := kp.producer.Close()
	if err != nil {
		return erax.Wrap(err, "failed to close kafka producer")
	}

	return nil
}

func (kp *KafkaProducer) Errors() <-chan error {
	return kp.errChan
}

func NewKafkaProducer(kafkaBroker string) (*KafkaProducer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	kafkaConfig.Producer.Retry.Max = 3
	kafkaConfig.Producer.Idempotent = false
	kafkaConfig.Producer.Flush.Frequency = 5 * time.Millisecond
	kafkaConfig.Producer.Flush.Bytes = 32 * 1024 // 32 KB
	kafkaConfig.Producer.Compression = sarama.CompressionLZ4
	kafkaConfig.Net.KeepAlive = 30 * time.Second
	kafkaConfig.Producer.Return.Successes = false
	kafkaConfig.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer([]string{kafkaBroker}, kafkaConfig)
	if err != nil {
		return nil, erax.Wrap(err, "failed to create kafka producer")
	}

	kp := &KafkaProducer{
		producer: producer,
		errChan:  make(chan error, 100),
	}

	go func() {
		defer close(kp.errChan)
		for err = range kp.producer.Errors() {
			kp.errChan <- err
		}
	}()

	return kp, nil
}
