package utils

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	_mqttHost 		= os.Getenv("MQTT_HOST")
	_mqttPort 		= os.Getenv("MQTT_PORT")
	_mqttUser 		= os.Getenv("MQTT_USER")
	_mqttPass 		= os.Getenv("MQTT_PASS")
)


func CreateMQTTOpts(clientId string) *mqtt.ClientOptions {
	dsn := "tcp://" + _mqttHost + ":" + _mqttPort

	connectHandler := func(client mqtt.Client) {
		log.Println("Connected to MQTT broker")
	}
	disconnectHandler := func(client mqtt.Client, err error) {
		log.Printf("MQTT: Disconnected from MQTT broker: %v", err)
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(dsn)
	opts.SetClientID(clientId)
	opts.SetUsername(_mqttUser)
	opts.SetPassword(_mqttPass)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = disconnectHandler

	return opts
}