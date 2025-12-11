package runtime

import (
	"go.uber.org/zap"
	"time"

	"github.com/DangeL187/erax"

	"ingress/internal/infra/metrics"
	"ingress/internal/shared/config"
)

type consumer interface {
	Run(messageHandler func([]byte)) error
	Stop() error
}

type ConsumerLoop struct {
	cfg *config.Config

	consumer consumer

	msgChanOut chan<- []byte
}

func (cl *ConsumerLoop) Run() error {
	err := cl.consumer.Run(cl.handleIncomingMessage)
	if err != nil {
		return erax.Wrap(err, "failed to run consumer")
	}

	return nil
}

func (cl *ConsumerLoop) Stop() {
	err := cl.consumer.Stop()
	if err != nil {
		zap.L().Error("failed to stop consumer", zap.Error(err))
	}
}

func (cl *ConsumerLoop) handleIncomingMessage(payload []byte) {
	start := time.Now()
	metrics.MessagesReceived.Inc()

	select {
	case cl.msgChanOut <- payload:
	default:
		metrics.MessagesDropped.Inc()
		// TODO: ask for retry
		zap.L().Debug("msgChanOut full: dropping message")
	}

	metrics.ConsumerLatency.Observe(time.Since(start).Seconds())
}

func NewConsumerLoop(cfg *config.Config, msgChanOut chan<- []byte, consumer consumer) (*ConsumerLoop, error) {
	cl := &ConsumerLoop{
		cfg:        cfg,
		consumer:   consumer,
		msgChanOut: msgChanOut,
	}

	return cl, nil
}
