package infra

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DangeL187/erax"
	"github.com/IBM/sarama"

	"consumer/internal/infra/metrics"
	"consumer/internal/shared/config"
)

type KafkaClickHouseFlusher struct {
	conn  clickhouse.Conn
	table string
}

type deviceData struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Battery   float64 `json:"battery"`
	Timestamp int64   `json:"timestamp"`
}

func (f *KafkaClickHouseFlusher) Flush(batch []*sarama.ConsumerMessage) {
	if len(batch) == 0 {
		return
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	chBatch, err := f.conn.PrepareBatch(ctx, "INSERT INTO "+f.table+" (id, latitude, longitude, altitude, battery, timestamp)")
	if err != nil {
		zap.L().Error("ClickHouse PrepareBatch failed", zap.Error(err))
		metrics.FlushErrors.Inc()
		return
	}

	for _, msg := range batch {
		var device deviceData
		if err = json.Unmarshal(msg.Value, &device); err != nil {
			zap.L().Error("JSON unmarshal failed", zap.Error(err))
			metrics.FlushErrors.Inc()
			continue
		}

		ts := time.Unix(device.Timestamp, 0)

		if err = chBatch.Append(
			device.ID,
			device.Latitude,
			device.Longitude,
			device.Altitude,
			device.Battery,
			ts,
		); err != nil {
			zap.L().Error("ClickHouse batch append failed", zap.Error(err))
			metrics.FlushErrors.Inc()
		}
	}

	if err = chBatch.Send(); err != nil {
		zap.L().Error("ClickHouse batch send failed", zap.Error(err))
	}

	metrics.BatchesFlushed.Inc()
	metrics.MessagesFlushed.Add(float64(len(batch)))
	metrics.FlushDuration.Observe(time.Since(start).Seconds())
}

func NewKafkaClickHouseFlusher(cfg *config.Config) (*KafkaClickHouseFlusher, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.ClickHouseDSN},
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouseDB,
			Username: cfg.ClickHouseUsername,
			Password: cfg.ClickHousePassword,
		},
		MaxOpenConns: 50,
		MaxIdleConns: 25,
		DialTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 10,
		Debug:        false,
	})
	if err != nil {
		return nil, erax.Wrap(err, "failed to open clickhouse")
	}

	return &KafkaClickHouseFlusher{
		conn:  conn,
		table: cfg.ClickHouseTable,
	}, nil
}
