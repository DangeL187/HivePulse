package config

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCAddr     string
	KafkaBroker  string
	KafkaTopic   string
	MQTTBroker   string
	MQTTClientID string
	MQTTTopic    string
	MsgChanSize  int
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		zap.L().Info("[env] no .env file found, skipping")
	}

	vars := map[string]string{
		"AUTH_GRPC":      os.Getenv("AUTH_GRPC"),
		"KAFKA_BROKER":   os.Getenv("KAFKA_BROKER"),
		"KAFKA_TOPIC":    os.Getenv("KAFKA_TOPIC"),
		"MQTT_BROKER":    os.Getenv("MQTT_BROKER"),
		"MQTT_CLIENT_ID": os.Getenv("MQTT_CLIENT_ID"),
		"MQTT_TOPIC":     os.Getenv("MQTT_TOPIC"),
		"MSG_CHAN_SIZE":  os.Getenv("MSG_CHAN_SIZE"),
	}

	for name, value := range vars {
		if value == "" {
			return nil, fmt.Errorf("missing required env var: %s", name)
		}
	}

	msgChanSize, err := strconv.Atoi(vars["MSG_CHAN_SIZE"])
	if err != nil {
		return nil, fmt.Errorf("invalid MSG_CHAN_SIZE: %s", vars["MSG_CHAN_SIZE"])
	}

	return &Config{
		MsgChanSize:  msgChanSize,
		GRPCAddr:     vars["AUTH_GRPC"],
		KafkaBroker:  vars["KAFKA_BROKER"],
		KafkaTopic:   vars["KAFKA_TOPIC"],
		MQTTBroker:   vars["MQTT_BROKER"],
		MQTTClientID: vars["MQTT_CLIENT_ID"],
		MQTTTopic:    vars["MQTT_TOPIC"],
	}, nil
}
