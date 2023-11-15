package producer

import (
	"log"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func ConnectProducerMQTT() mqtt.Client {
	opts := utils.CreateMQTTOpts("mqtt-producer-backend-a10")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func PublishMQTT(client mqtt.Client, topic string, qos byte, retained bool, payload string) {
	token := client.Publish(topic, qos, retained, payload)
	token.Wait()
	log.Printf("MQTT: Published to topic: %s\n", topic)
}
