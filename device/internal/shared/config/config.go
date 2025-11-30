package config

import "time"

type Config struct {
	AuthServerURL string
	MqttBrokerURL string

	DeviceID       string
	DevicePassword string

	ConnectRetryInterval   time.Duration
	MaxReconnectInterval   time.Duration
	PublishMetricsInterval time.Duration
}

func NewConfig() *Config {
	return &Config{
		AuthServerURL:          "http://localhost:8000",
		MqttBrokerURL:          "tcp://localhost:1883",
		DeviceID:               "",
		DevicePassword:         "secret",
		ConnectRetryInterval:   5 * time.Second,
		MaxReconnectInterval:   30 * time.Second,
		PublishMetricsInterval: time.Millisecond * 10,
	}
}
