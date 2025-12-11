package mqtt

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"time"

	"github.com/eclipse/paho.mqtt.golang"

	"device/internal/http"
	"device/internal/shared/config"
	"device/internal/shared/tokens"
	m "device/metrics"
)

type deviceData struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Battery   float64 `json:"battery"`
	Timestamp int64   `json:"timestamp"`
	Token     string  `json:"token"`
}

type MetricsService struct {
	authService *http.AuthService
	cfg         *config.Config
	tokens      *tokens.Tokens
	mqttClient  mqtt.Client

	dataChan chan deviceData
}

func (ms *MetricsService) Run(ctx context.Context) {
	zap.L().Info("MetricsService started")

	if token := ms.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		zap.L().Error("failed to connect to MQTT broker", zap.Error(token.Error()))
		return
	}

	if !ms.authService.Login(ctx) {
		return
	}

	m.AuthCounter.Add(1)

	token := ms.mqttClient.Subscribe(
		"devices/"+ms.cfg.DeviceID+"/auth_response",
		0,
		func(_ mqtt.Client, msg mqtt.Message) {
			ms.handleAuthResponse(ctx, msg)
		},
	)
	if token.Wait() && token.Error() != nil {
		zap.L().Error("failed to subscribe to auth response topic", zap.Error(token.Error()))
		return
	}

	ms.startPublishing(ctx)
	ms.startMetricsPutting(ctx)
}

func (ms *MetricsService) Stop() {
	zap.L().Info("MetricsService stopping...")
	if ms.mqttClient.IsConnected() {
		ms.mqttClient.Disconnect(250)
	}
	zap.L().Info("MetricsService stopped")
}

func (ms *MetricsService) startPublishing(ctx context.Context) {
	go func() {
		for {
			if !ms.authService.WaitForAuth(ctx) {
				return
			}

			select {
			case <-ctx.Done():
				return
			case msg := <-ms.dataChan:
				ms.publish(msg)
			}
		}
	}()
}

func (ms *MetricsService) startMetricsPutting(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(ms.cfg.PublishMetricsInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ms.putInQueue()
			}
		}
	}()
}

func (ms *MetricsService) putInQueue() {
	// [WARNING] Simulated data:
	msg := deviceData{
		ID:        ms.cfg.DeviceID,
		Latitude:  55.7,
		Longitude: 37.6,
		Altitude:  120,
		Battery:   88.0,
		Timestamp: time.Now().Unix(),
	}
	ms.dataChan <- msg
}

func (ms *MetricsService) publish(msg deviceData) {
	msg.Token = ms.tokens.GetAccess()

	payload, err := json.Marshal(msg)
	if err != nil {
		zap.L().Error("failed to marshal device data msg", zap.Error(err))
		return
	}

	token := ms.mqttClient.Publish("devices/telemetry", 0, false, payload)
	if token.Wait() && token.Error() != nil {
		zap.L().Error("failed to publish device data payload", zap.Error(token.Error()))
		return
	}

	m.Counter.Add(1)
}

type authResponse struct {
	Error string `json:"error"`
}

func (ms *MetricsService) handleAuthResponse(ctx context.Context, msg mqtt.Message) {
	var resp authResponse
	if err := json.Unmarshal(msg.Payload(), &resp); err != nil {
		zap.L().Error("failed to unmarshal message payload", zap.Error(err))
		return
	}

	if resp.Error != "" {
		zap.L().Debug("Auth failed", zap.String("Error", resp.Error))
		ms.authService.ResetAuth()
		ms.authService.Refresh(ctx)
		return
	}

	zap.L().Info("Auth successful")
}

func NewMetricsService(cfg *config.Config, tokens *tokens.Tokens, authService *http.AuthService) *MetricsService {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.MqttBrokerURL).
		SetClientID(cfg.DeviceID).
		SetAutoReconnect(true).
		SetConnectRetryInterval(cfg.ConnectRetryInterval).
		SetMaxReconnectInterval(cfg.MaxReconnectInterval)
	mqttClient := mqtt.NewClient(opts)

	ms := &MetricsService{
		authService: authService,
		cfg:         cfg,
		mqttClient:  mqttClient,
		tokens:      tokens,
		dataChan:    make(chan deviceData, 100),
	}

	return ms
}
