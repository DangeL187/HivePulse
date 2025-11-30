package config

import (
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCAddr    string
	KafkaBroker string
	KafkaTopic  string
	MQTTBroker  string
	MQTTTopic   string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		zap.L().Info("[env] no .env file found, skipping")
	}

	vars := map[string]string{
		"AUTH_GRPC":    os.Getenv("AUTH_GRPC"),
		"KAFKA_BROKER": os.Getenv("KAFKA_BROKER"),
		"KAFKA_TOPIC":  os.Getenv("KAFKA_TOPIC"),
		"MQTT_BROKER":  os.Getenv("MQTT_BROKER"),
		"MQTT_TOPIC":   os.Getenv("MQTT_TOPIC"),
	}

	for name, value := range vars {
		if value == "" {
			return nil, fmt.Errorf("missing required env var: %s", name)
		}
	}

	return &Config{
		GRPCAddr:    vars["AUTH_GRPC"],
		KafkaBroker: vars["KAFKA_BROKER"],
		KafkaTopic:  vars["KAFKA_TOPIC"],
		MQTTBroker:  vars["MQTT_BROKER"],
		MQTTTopic:   vars["MQTT_TOPIC"],
	}, nil
}
