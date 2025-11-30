package runtime

import (
	"context"
	"sync"
	"time"

	"consumer/internal/features/flusher/domain"
	"consumer/internal/shared/config"
)

type MessageBatchFlusher[T any] struct {
	cfg       *config.Config
	flusher   domain.Flusher[T]
	msgChanIn <-chan *T
	wg        sync.WaitGroup
}

func (mbf *MessageBatchFlusher[T]) Run(ctx context.Context, workerCount int) {
	mbf.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer mbf.wg.Done()

			ticker := time.NewTicker(mbf.cfg.BatchInterval)
			defer ticker.Stop()

			batch := make([]*T, 0, mbf.cfg.BatchSize)

			for {
				select {
				case msg, ok := <-mbf.msgChanIn:
					if !ok {
						mbf.flusher.Flush(batch)
						return
					}

					batch = append(batch, msg)

					if len(batch) >= mbf.cfg.BatchSize {
						mbf.flusher.Flush(batch)
						batch = batch[:0]
					}
				case <-ctx.Done():
					mbf.flusher.Flush(batch)
					return
				case <-ticker.C:
					mbf.flusher.Flush(batch)
					batch = batch[:0]
				}
			}
		}()
	}
}

func (mbf *MessageBatchFlusher[T]) Stop() {
	mbf.wg.Wait()
}

func NewMessageBatchFlusher[T any](cfg *config.Config, msgChanIn <-chan *T, flusher domain.Flusher[T]) *MessageBatchFlusher[T] {
	return &MessageBatchFlusher[T]{
		cfg:       cfg,
		flusher:   flusher,
		msgChanIn: msgChanIn,
	}
}
