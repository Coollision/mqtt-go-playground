package mqtt

import (
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

type Config struct {
	BrokerAddress  string `required:"true"`
	BrokerId       string `required:"true"`
	BrokerUser     string `required:"true"`
	BrokerPassword string `secret:"true" required:"true"`
	CleanSession   bool   `default:"true"`
}

var client MQTT.Client

func Init(cfg *Config) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(cfg.BrokerAddress)
	opts.SetClientID(cfg.BrokerId)
	opts.SetUsername(cfg.BrokerUser)
	opts.SetPassword(cfg.BrokerPassword)
	opts.SetCleanSession(cfg.CleanSession)

	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		logrus.Errorf("received unexpected information on topic: %s, with msg: %s", msg.Topic(), string(msg.Payload()))
	})

	client = MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func Disconnect() {
	client.Disconnect(250)
	logrus.Infof("Mqtt Subscriber Disconnected")
}

func Subscribe(topic string, handle func(topic string, message []byte)) {
	handlefunc := func(c MQTT.Client, m MQTT.Message) {
		logrus.Tracef("Received message on topic: %s, with msg: %s", m.Topic(), string(m.Payload()))
		handle(m.Topic(), m.Payload())
	}

	if token := client.Subscribe(topic, byte(0), handlefunc); token.Error() != nil {
		logrus.Fatal(token.Error())
	}

	logrus.Infof("Mqtt Subscriber Subscribed to %s", topic)
}

func Unsubscribe(topic string) {
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		logrus.Fatal(token.Error())
	}
	logrus.Infof("Mqtt Subscriber Unsubscribed from %s", topic)
}

func DeleteTopic(topic string) {
	if token := client.Publish(topic, byte(0), true, ""); token.Wait() && token.Error() != nil {
		logrus.Fatal(token.Error())
	}
	logrus.Infof("Mqtt Subscriber Deleted Topic %s", topic)
}

func Publish(topic string, payload string) error {
	if token := client.Publish(topic, byte(0), true, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	logrus.Tracef("Mqtt Publisher Published to %s with msg: %s", topic, payload)
	return nil
}

func PublishStruct(topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if token := client.Publish(topic, byte(0), true, data); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	logrus.Tracef("Mqtt Publisher Published to %s", topic)
	return nil
}
