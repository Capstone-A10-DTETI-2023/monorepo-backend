package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/controller"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pusher/pusher-http-go/v5"
	"gorm.io/gorm"
)

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func ConnectConsumerMQTT(dbData *gorm.DB, websocket *pusher.Client) mqtt.Client {
	messageReceiveHandler := func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("MQTT: Received message from topic: %s\n", msg.Topic())
		payloadStr := string(msg.Payload())
		log.Println(payloadStr)

		var NodeSensorData model.NodeSensorData
		if err := json.Unmarshal([]byte(payloadStr), &NodeSensorData); err != nil {
			log.Printf("MQTT: Error unmarshalling sensor data: %v", err)
		} else {
			if err := controller.InsertDataSensorToDB(NodeSensorData, dbData); err != nil {
				log.Printf("MQTT: Error inserting sensor data: %v", err)
			} else {
				log.Printf("MQTT: Inserted sensor data to DB")
				websocket.Trigger("sensordata", "new-sensor-data", NodeSensorData)
			}
		}
	}

	opts := utils.CreateMQTTOpts("mqtt-consumer-backend-a10")
	opts.SetOrderMatters(false)
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
	log.Printf("MQTT: Subscribed to topic: %s\n", topic)
}
