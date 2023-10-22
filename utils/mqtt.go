package utils

import (
	"fmt"
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

func ConnectMQTT() mqtt.Client {
	dsn := "tcp://" + _mqttHost + ":" + _mqttPort

	connectHandler := func(client mqtt.Client) {
		log.Println("Connected to MQTT broker")
	}
	disconnectHandler := func(client mqtt.Client, err error) {
		log.Printf("Disconnected from MQTT broker: %v", err)
	}

	messageReceiveHandler := func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(dsn)
	opts.SetClientID("capstone-backend-a10")
	opts.SetUsername(_mqttUser)
	opts.SetPassword(_mqttPass)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = disconnectHandler
	opts.SetDefaultPublishHandler(messageReceiveHandler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("MQTT: %s", token.Error()))
	}

	return client
}

func SubMQTT(client mqtt.Client, topic string, qos int) {
	token := client.Subscribe(topic, byte(qos), nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}
