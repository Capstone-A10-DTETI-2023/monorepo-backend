package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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

func ConnectMQTT() mqtt.Client {
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

		var NodeSensorData model.NodeSensorData
		if err := json.Unmarshal([]byte(payloadStr), &NodeSensorData); err != nil {
			log.Printf("MQTT: Error unmarshalling sensor data: %v", err)
		} else {
			if err := InsertDataSensor(NodeSensorData); err != nil {
				log.Printf("MQTT: Error inserting sensor data: %v", err)
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

func InsertDataSensor(sensorData model.NodeSensorData) error {
	db := model.ConnectDBTS()

	if sensorData.Timestamp == "" {
		sensorData.Timestamp = time.Now().String()
	}

	var timestamp time.Time
	bufferTime := sensorData.Timestamp.(string)
	timestamp, err := time.Parse("2006-01-02 15:04:05", bufferTime)
	timestampStr := timestamp.Format(time.RFC3339)
	if err != nil {
		return Error{"Error parsing timestamp"}
	}

	if sensorData.NodeID == "" || sensorData.SensorID == "" || sensorData.Value == "" {
		return Error{"Missing required field"}
	}

	sensorData = model.NodeSensorData{
		Timestamp: timestampStr,
		NodeID: sensorData.NodeID,
		SensorID: sensorData.SensorID,
		Value: sensorData.Value,
	}

	if err := sensorData.CheckDuplicateSensorData(db); err != nil {
		return err
	}

	if err := sensorData.CreateNodeSensorData(db); err != nil {
		return err
	}

	return nil
}