package mqtt

import (
	"github.com/DangeL187/erax"
	"github.com/eclipse/paho.mqtt.golang"
)

type Publisher struct {
	mqttClient mqtt.Client
}

func (p *Publisher) Publish(topic string, payload any) error {
	token := p.mqttClient.Publish(topic, 0, false, payload)
	if token.Wait() && token.Error() != nil {
		return erax.Wrap(token.Error(), "failed to publish device data payload")
	}

	return nil
}

func NewPublisher(mqttClient mqtt.Client) *Publisher {
	return &Publisher{
		mqttClient: mqttClient,
	}
}
