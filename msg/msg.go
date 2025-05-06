package msg

import (
	"log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Msg struct {
	topic       string
	broker      string
	user        string
	password    string
	mqttOptions *mqtt.ClientOptions
	mqtt        mqtt.Client
	mqttChan    chan int
}

func New(topic string, broker string, user string, password string) *Msg {
	m := Msg{
		topic:    strings.TrimSpace(topic),
		broker:   strings.TrimSpace(broker),
		user:     strings.TrimSpace(user),
		password: strings.TrimSpace(password),
		mqttChan: nil,
	}

	m.mqtt = connect(&m)

	return &m
}

func connect(m *Msg) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(m.broker)
	opts.SetClientID("Door_Control")
	opts.SetUsername(m.user)
	opts.SetPassword(m.password)
	opts.SetCleanSession(false)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	// Uncomment to DEBUG MQTT
	//mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	//mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	//mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	//mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	log.Println("MQTT Connected Successfully")

	return client
}

func (m *Msg) Listen(c chan bool) {
	m.mqtt.Subscribe("/frontdoor/unlock", 0, func(client mqtt.Client, msg mqtt.Message) {
		go MqttMessageHandler(client, msg, c)
	})
}

func MqttMessageHandler(client mqtt.Client, msg mqtt.Message, c chan bool) {
	//fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	if string(msg.Payload()) == "1" && msg.Topic() == "/frontdoor/unlock" && msg.Retained() == false {
		c <- true
	}

}

func (m *Msg) Message(msg string) {
	if token := m.mqtt.Publish(m.topic, 1, false, msg); token.Wait() && token.Error() != nil {
		log.Println("MQTT ", token.Error())
	}
}
