package app

import (
	"context"
	"go.uber.org/zap"
	"runtime"

	"github.com/DangeL187/erax"
	"github.com/IBM/sarama"

	consumerInfra "consumer/internal/features/consumer/infra"
	consumerRuntime "consumer/internal/features/consumer/runtime"
	flusherInfra "consumer/internal/features/flusher/infra"
	flusherRuntime "consumer/internal/features/flusher/runtime"
	"consumer/internal/shared/config"
)

type App struct {
	msgChan chan *sarama.ConsumerMessage

	cfg *config.Config

	consumerLoop        *consumerRuntime.ConsumerLoop[sarama.ConsumerMessage]
	messageBatchFlusher *flusherRuntime.MessageBatchFlusher[sarama.ConsumerMessage]

	cancel context.CancelFunc
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	a.messageBatchFlusher.Run(ctx, runtime.NumCPU()*2)

	errChanIn := a.consumerLoop.Run(ctx)
	go func() {
		for err := range errChanIn {
			zap.L().Error("Consumer error", zap.Error(err))
		}
	}()

	zap.L().Info("Consumer started")
}

func (a *App) Stop() {
	a.cancel()

	err := a.consumerLoop.Stop()
	if err != nil {
		err = erax.Wrap(err, "failed to stop consumer loop")
		zap.S().Errorf("\n%f", err)
	}

	close(a.msgChan)

	a.messageBatchFlusher.Stop()
}

func NewApp() (*App, error) {
	app := &App{
		msgChan: make(chan *sarama.ConsumerMessage, 10000),
	}

	var err error
	app.cfg, err = config.NewConfig()
	if err != nil {
		return nil, erax.Wrap(err, "failed to load config")
	}

	flusher, err := flusherInfra.NewKafkaClickHouseFlusher(app.cfg)
	if err != nil {
		return nil, erax.Wrap(err, "failed to initialize kafka-clickhouse flusher")
	}
	app.messageBatchFlusher = flusherRuntime.NewMessageBatchFlusher[sarama.ConsumerMessage](app.cfg, app.msgChan, flusher)

	kafkaConsumer, err := consumerInfra.NewKafkaConsumer(app.cfg)
	if err != nil {
		return nil, erax.Wrap(err, "failed to create kafka consumer")
	}
	app.consumerLoop = consumerRuntime.NewConsumerLoop[sarama.ConsumerMessage](kafkaConsumer, app.msgChan)

	return app, nil
}
