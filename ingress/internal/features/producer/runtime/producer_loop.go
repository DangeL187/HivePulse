package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"runtime"
	"sync"

	"github.com/DangeL187/erax"

	"ingress/internal/infra/metrics"
)

type producer interface {
	Produce(topic string, payload []byte)
	Close() error
	Errors() <-chan error
}

type authService interface {
	Run(ctx context.Context)
	Stop()
	Auth(deviceID, deviceToken string) error
}

type ProducerLoop struct {
	msgChanIn <-chan []byte

	authService authService
	producer    producer

	bufPool *sync.Pool

	drainWg sync.WaitGroup
	sendWg  sync.WaitGroup
}

func (ps *ProducerLoop) Run(ctx context.Context, kafkaTopic string) {
	ps.authService.Run(ctx)
	ps.runDrainWorkers(ctx, 1)
	ps.runProducerWorkers(ctx, runtime.NumCPU()*2, kafkaTopic)
}

func (ps *ProducerLoop) Stop() {
	ps.sendWg.Wait()
	if err := ps.producer.Close(); err != nil {
		zap.L().Error("failed to close producer", zap.Error(err))
	}
	ps.drainWg.Wait()

	ps.authService.Stop()
}

func (ps *ProducerLoop) runDrainWorkers(ctx context.Context, workerCount int) {
	ps.drainWg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer ps.drainWg.Done()
			for {
				select {
				case err, ok := <-ps.producer.Errors():
					if !ok {
						return
					}
					metrics.MessagesSendErrors.Inc()
					// TODO: DLQ
					zap.L().Error("Kafka send error", zap.Error(err))
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (ps *ProducerLoop) runProducerWorkers(ctx context.Context, workerCount int, kafkaTopic string) {
	ps.sendWg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer ps.sendWg.Done()
			for {
				select {
				case payload, ok := <-ps.msgChanIn:
					if !ok {
						return
					}
					processedPayload, err := ps.processMessage(payload)
					if err != nil {
						zap.S().Errorf("failed to process message:\n%f", err)
					}
					ps.producer.Produce(kafkaTopic, processedPayload)
					metrics.MessagesSent.Inc()
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

type deviceData struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Battery   float64 `json:"battery"`
	Timestamp int64   `json:"timestamp"`
	Token     string  `json:"token"`
}

type kafkaDeviceData struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Battery   float64 `json:"battery"`
	Timestamp int64   `json:"timestamp"`
}

func (ps *ProducerLoop) processMessage(payload []byte) ([]byte, error) {
	var data deviceData
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, erax.Wrap(err, "failed to unmarshal message payload")
	}

	err := ps.authService.Auth(data.ID, data.Token)
	if err != nil {
		metrics.AuthFail.Inc()
		return nil, erax.Wrap(err, "failed to authenticate")
	}
	metrics.AuthSuccess.Inc()

	kafkaData := kafkaDeviceData{
		ID:        data.ID,
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
		Altitude:  data.Altitude,
		Battery:   data.Battery,
		Timestamp: data.Timestamp,
	}

	buf := ps.bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	encoder := json.NewEncoder(buf)
	if err = encoder.Encode(&kafkaData); err != nil {
		ps.bufPool.Put(buf)
		return nil, erax.Wrap(err, "failed to encode data")
	}

	out := make([]byte, buf.Len())
	copy(out, buf.Bytes())
	ps.bufPool.Put(buf)

	return out, nil
}

func NewProducerLoop(msgChanIn <-chan []byte, authService authService, producer producer) *ProducerLoop {
	return &ProducerLoop{
		msgChanIn:   msgChanIn,
		authService: authService,
		producer:    producer,
		bufPool: &sync.Pool{
			New: func() interface{} {
				b := bytes.NewBuffer(make([]byte, 0, 128)) // 128 bytes
				return b
			},
		},
	}
}
