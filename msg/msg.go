package msg

import (
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Msg struct {
	topic    string
	broker   string
	user     string
	password string
}

func New(topic string, broker string, user string, password string) *Msg {
	return &Msg{
		topic:    strings.TrimSpace(topic),
		broker:   strings.TrimSpace(broker),
		user:     strings.TrimSpace(user),
		password: strings.TrimSpace(password),
	}
}

func (m *Msg) Message(msg string) {
	options := mqtt.NewClientOptions()
	options.AddBroker(m.broker)
	options.SetClientID("RFID")
	options.SetUsername(m.user)
	options.SetPassword(m.password)
	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("MQTT ", token.Error())
		return
	}

	if token := client.Publish(m.topic, 1, false, msg); token.Wait() && token.Error() != nil {
		log.Println("MQTT ", token.Error())
	}

	client.Disconnect(250)

	return
}
