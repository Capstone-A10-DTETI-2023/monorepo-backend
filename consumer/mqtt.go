package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/controller"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pusher/pusher-http-go/v5"
	"gorm.io/gorm"
)

var (
	_mqttHost 		= os.Getenv("MQTT_HOST")
	_mqttPort 		= os.Getenv("MQTT_PORT")
	_mqttUser 		= os.Getenv("MQTT_USER")
	_mqttPass 		= os.Getenv("MQTT_PASS")
)

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func ConnectMQTT(dbData *gorm.DB, websocket *pusher.Client) mqtt.Client {
	dsn := "tcp://" + _mqttHost + ":" + _mqttPort

	connectHandler := func(client mqtt.Client) {
		log.Println("Connected to MQTT broker")
	}
	disconnectHandler := func(client mqtt.Client, err error) {
		log.Printf("MQTT: Disconnected from MQTT broker: %v", err)
	}

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

	opts := mqtt.NewClientOptions()
	opts.AddBroker(dsn)
	opts.SetClientID("capstone-backend-a10")
	opts.SetUsername(_mqttUser)
	opts.SetPassword(_mqttPass)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = disconnectHandler
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
