package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bseishen/golang_doorcontrol/api"
	"github.com/bseishen/golang_doorcontrol/msg"
	"github.com/bseishen/golang_doorcontrol/store"
	"github.com/jacobsa/go-serial/serial"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	LogFile           string `default:"./rfid.log"`
	SerialPort        string `default:"/dev/ttyUSB0"`
	BaudRate          int    `default:"9600"`
	DataBits          int    `default:"8"`
	StopBits          int    `default:"1"`
	MinimumReadSize   int    `default:"1"`
	DBFile            string `default:"./rfid.sqlite"`
	ApiUpdateInterval int    `default:"30"`
	ApiUrl            string `default:"https://rfid.midsouthmakers.org/api"`
	ApiKey            string
	MqttServer        string `default:"tcp://192.168.10.5:1883"`
	MqttPassword      string `default:""`
	MqttUsername      string `default:""`
	MqttTopic         string `default:"/frontdoor/notifications"`
	MqttUnlockTopic   string `default:"/frontdoor/unlock"`
}

var (
	config Configuration
	sp     io.ReadWriteCloser
	s      store.Store
	key    int
	pw     string
	m      msg.Msg
	a      api.Api
)

func main() {

	f, err := os.OpenFile("./rfid.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	configure()

	a = *(api.New(config.ApiUrl, config.ApiKey))
	s = *(store.New(config.DBFile))
	m = *(msg.New(config.MqttTopic, config.MqttServer, config.MqttUsername, config.MqttPassword))
	log.Println("MQTT Configuration Complete")

	m.Message("RFID Application Started")
	//Update the database immediatly.
	_, d := a.CheckForUpdates()
	s.UpdateDatabase(d)

	ticker := time.NewTicker(time.Second * time.Duration(config.ApiUpdateInterval))
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				updateRequired, dbData := a.CheckForUpdates()
				if updateRequired {
					s.UpdateDatabase(dbData)
					m.Message("RFID Database Updated")
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	for {
		PollSerial()
	}

	log.Fatal("We somehow have arrived at a point that is never supposed to happen, cleaning up\n")
	sp.Close()

	return
}

func configure() {
	// Collect configurations from env
	err := envconfig.Process("rfid", &config)
	if err != nil {
		log.Fatalf("Unable to process configuration: %v\n", err.Error())
	}

	log.Println("Opening serial port " + config.SerialPort)

	//Create serial port
	sp, err = serial.Open(serial.OpenOptions{
		PortName:        config.SerialPort,
		BaudRate:        uint(config.BaudRate),
		DataBits:        uint(config.DataBits),
		StopBits:        uint(config.StopBits),
		MinimumReadSize: uint(config.MinimumReadSize),
	})

	if err != nil {
		log.Fatalf("Error opening serial port: %v\n", err)
	}

	log.Println("Configuration complete")
}

func PollSerial() {

	//read serial till you hit a new line
	buf := bufio.NewReader(sp)
	b, err := buf.ReadBytes('\n')
	if err != nil {
		log.Printf("Error reading: %v\n", err)
	}

	str := strings.ToLower(strings.TrimSpace(string(b)))

	if str != "" {
		//log.Println("Message Received:", str)
		switch str[0] {
		// RFID Key
		case 'r':
			key, err = strconv.Atoi(strings.TrimPrefix(str, "r"))
			if err != nil {
				log.Printf("Unable to convert [%v] to an integer: %v\n", string(b[1:]), err.Error())
			}
		// Escape button
		case 0x1B:
			key = 0
			pw = ""
			err = nil
		// Password
		default:
			pw = str
		}

		if pw != "" && key != 0 {

			user, err := s.FindUser(key, pw)
			if err != nil {
				log.Printf("Error: %v\n", err.Error())
				m.Message(fmt.Sprintf(fmt.Sprintf("Error: %v\n", err.Error())))
				a.SendLoginAttempt(key, err.Error(), "failure")
				//Send error to keypad by sending an 'E'
				WriteByte(byte('E'))
			} else {
				//Unlock door by sending an 'O'
				WriteByte(byte('O'))
				log.Printf("Access Granted for user %s (%v)", user.IrcName, user.Key)
				m.Message(fmt.Sprintf("Access Granted for user %s", user.IrcName))
				a.SendLoginAttempt(key, fmt.Sprintf("Access Granted for user %s (%v)", user.IrcName, user.Key), "success")
			}

			//clear key and pincode
			pw = ""
			key = 0
		}

	}
}

func WriteByte(a byte) {
	b := []byte{0x00}
	b[0] = a
	sp.Write(b)
}
