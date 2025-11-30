package config

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	BatchInterval time.Duration
	BatchSize     int

	ClickHouseDB       string
	ClickHouseDSN      string
	ClickHousePassword string
	ClickHouseTable    string
	ClickHouseUsername string
	KafkaBroker        string
	KafkaTopic         string
	KafkaGroupID       string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		zap.L().Info("[env] no .env file found, skipping")
	}

	vars := map[string]string{
		"CLICKHOUSE_DB":       os.Getenv("CLICKHOUSE_DB"),
		"CLICKHOUSE_DSN":      os.Getenv("CLICKHOUSE_DSN"),
		"CLICKHOUSE_TABLE":    os.Getenv("CLICKHOUSE_TABLE"),
		"CLICKHOUSE_USERNAME": os.Getenv("CLICKHOUSE_USERNAME"),
		"KAFKA_BROKER":        os.Getenv("KAFKA_BROKER"),
		"KAFKA_GROUP_ID":      os.Getenv("KAFKA_GROUP_ID"),
		"KAFKA_TOPIC":         os.Getenv("KAFKA_TOPIC"),
	}

	for name, value := range vars {
		if value == "" {
			return nil, fmt.Errorf("missing required env var: %s", name)
		}
	}

	return &Config{
		BatchInterval:      time.Second,
		BatchSize:          10000,
		ClickHouseDB:       vars["CLICKHOUSE_DB"],
		ClickHouseDSN:      vars["CLICKHOUSE_DSN"],
		ClickHousePassword: os.Getenv("CLICKHOUSE_PASSWORD"),
		ClickHouseTable:    vars["CLICKHOUSE_TABLE"],
		ClickHouseUsername: vars["CLICKHOUSE_USERNAME"],
		KafkaBroker:        vars["KAFKA_BROKER"],
		KafkaGroupID:       vars["KAFKA_GROUP_ID"],
		KafkaTopic:         vars["KAFKA_TOPIC"],
	}, nil
}
