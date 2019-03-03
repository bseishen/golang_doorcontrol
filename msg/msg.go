package msg

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type Msg struct {
	topic    string
	broker   string
	user     string
	password string
}

func New(topic string, broker string, user string, password string) *Msg {
	return &Msg{
		topic:    topic,
		broker:   broker,
		user:     user,
		password: password,
	}
}

func (m *Msg) Message(msg string) {
	o := mqtt.NewClientOptions()
	o.AddBroker(m.broker)
	o.SetClientID("RFID")
	o.SetUsername(m.user)
	o.SetPassword(m.password)
	client := mqtt.NewClient(o)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("MQTT ", token.Error())
		return
	}

	if token := client.Publish(m.topic, 1, false, msg); token.Wait() && token.Error() != nil {
		log.Println("MQTT ", token.Error())
		return
	}
}
