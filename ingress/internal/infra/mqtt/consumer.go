package mqtt

import (
	"github.com/DangeL187/erax"
	"github.com/eclipse/paho.mqtt.golang"
)

type Consumer struct {
	mqttClient mqtt.Client
	mqttTopic  string
}

func (c *Consumer) Run(messageHandler func([]byte)) error {
	if token := c.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return erax.Wrap(token.Error(), "failed to connect to mqtt broker")
	}

	callback := func(_ mqtt.Client, msg mqtt.Message) {
		messageHandler(msg.Payload())
	}
	if token := c.mqttClient.Subscribe(c.mqttTopic, 0, callback); token.Wait() && token.Error() != nil {
		return erax.Wrap(token.Error(), "failed to subscribe to mqtt topic")
	}

	return nil
}

func (c *Consumer) Stop() error {
	var err error

	if token := c.mqttClient.Unsubscribe(c.mqttTopic); token.Wait() && token.Error() != nil {
		err = erax.Wrap(token.Error(), "failed to unsubscribe from mqtt topic")
	}
	c.mqttClient.Disconnect(250)

	return err
}

func NewConsumer(mqttClient mqtt.Client, mqttTopic string) *Consumer {
	return &Consumer{
		mqttClient: mqttClient,
		mqttTopic:  mqttTopic,
	}
}
