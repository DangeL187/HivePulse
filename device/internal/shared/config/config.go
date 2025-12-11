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
		AuthServerURL:          "http://localhost:30080", // 8000 - local, 30080 - k8s
		MqttBrokerURL:          "tcp://localhost:31883",  // 1883 - local, 31883 - k8s
		DeviceID:               "",
		DevicePassword:         "secret",
		ConnectRetryInterval:   5 * time.Second,
		MaxReconnectInterval:   30 * time.Second,
		PublishMetricsInterval: time.Millisecond * 10,
	}
}
