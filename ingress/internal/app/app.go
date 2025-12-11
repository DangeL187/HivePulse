package app

import (
	"context"
	"go.uber.org/zap"

	"github.com/DangeL187/erax"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"

	authInfra "ingress/internal/features/auth/infra"
	authRuntime "ingress/internal/features/auth/runtime"
	consumerRuntime "ingress/internal/features/consumer/runtime"
	producerInfra "ingress/internal/features/producer/infra"
	producerRuntime "ingress/internal/features/producer/runtime"
	infraMqtt "ingress/internal/infra/mqtt"
	"ingress/internal/shared/config"
)

type App struct {
	msgChan chan []byte

	cfg *config.Config

	consumerLoop *consumerRuntime.ConsumerLoop
	producerLoop *producerRuntime.ProducerLoop

	cancel context.CancelFunc
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	a.producerLoop.Run(ctx, a.cfg.KafkaTopic)

	err := a.consumerLoop.Run()
	if err != nil {
		return erax.Wrap(err, "failed to run ingress service")
	}

	zap.L().Info("Ingress started")

	return nil
}

func (a *App) Stop() {
	a.cancel()

	a.consumerLoop.Stop()
	close(a.msgChan)
	a.producerLoop.Stop()
}

func NewApp() (*App, error) {
	app := &App{}

	var err error
	app.cfg, err = config.NewConfig()
	if err != nil {
		return nil, erax.Wrap(err, "failed to load config")
	}

	app.msgChan = make(chan []byte, app.cfg.MsgChanSize)

	// publisher and consumer
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://" + app.cfg.MQTTBroker).
		SetClientID(app.cfg.MQTTClientID + uuid.New().String())
	mqttClient := mqtt.NewClient(opts)

	publisher := infraMqtt.NewPublisher(mqttClient)
	consumer := infraMqtt.NewConsumer(mqttClient, app.cfg.MQTTTopic)

	// ConsumerLoop
	app.consumerLoop, err = consumerRuntime.NewConsumerLoop(app.cfg, app.msgChan, consumer)
	if err != nil {
		return nil, erax.Wrap(err, "failed to create consumer loop")
	}

	// Auth Service
	authenticator, err := authInfra.NewGRPCAuthenticator(app.cfg.GRPCAddr)
	if err != nil {
		return nil, erax.Wrap(err, "failed to create authenticator")
	}

	authService, err := authRuntime.NewAuthService(authenticator, publisher)
	if err != nil {
		return nil, erax.Wrap(err, "failed to create auth service")
	}

	// ProducerLoop
	producer, err := producerInfra.NewKafkaProducer(app.cfg.KafkaBroker)
	if err != nil {
		return nil, erax.Wrap(err, "failed to create producer")
	}

	app.producerLoop = producerRuntime.NewProducerLoop(app.msgChan, authService, producer)

	return app, nil
}
