package runtime

import (
	"context"
	"errors"
	"time"

	"github.com/DangeL187/erax"

	"consumer/internal/infra/metrics"
)

type consumer[T any] interface {
	Run(ctx context.Context, msgHandler func(message *T)) error
	Stop() error
}

type ConsumerLoop[T any] struct {
	consumer   consumer[T]
	msgChanOut chan<- *T
}

func (cl *ConsumerLoop[T]) Run(ctx context.Context) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		err := cl.consumer.Run(ctx, cl.handleIncomingMessage)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}

			select {
			case errChan <- err:
			default:
			}
		}
	}()

	return errChan
}

func (cl *ConsumerLoop[T]) Stop() error {
	err := cl.consumer.Stop()
	if err != nil {
		return erax.Wrap(err, "failed to stop consumer")
	}

	return nil
}

func (cl *ConsumerLoop[T]) handleIncomingMessage(msg *T) {
	start := time.Now()
	cl.msgChanOut <- msg
	duration := time.Since(start).Seconds()
	metrics.MessagesConsumed.Inc()
	metrics.ConsumerLatency.Observe(duration)
}

func NewConsumerLoop[T any](consumer consumer[T], msgChanOut chan<- *T) *ConsumerLoop[T] {
	return &ConsumerLoop[T]{
		consumer:   consumer,
		msgChanOut: msgChanOut,
	}
}
