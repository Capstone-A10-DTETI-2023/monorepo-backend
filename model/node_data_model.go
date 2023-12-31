package model

import (
	"context"
	"fmt"
	"log"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/utils"
	"github.com/gofiber/fiber/v2"
	pgx "github.com/jackc/pgx/v5"
)

type NodeSensorData struct {
	Timestamp interface{}     	`json:"timestamp"`
	NodeID    string          	`json:"node_id"`
	SensorID  string          	`json:"sensor_id"`
	Value     string       		`json:"value"`
}

type NodeSensorDataGet struct {
	Timestamp 	interface{}     	`json:"timestamp"`
	Value   	string     			`json:"value"`
}

type NodeActuatorData struct {
	NodeID    	string          	`json:"node_id"`
	ActuatorID  string          	`json:"actuator_id"`
	Action   	string       		`json:"action"`
	Value     	string       		`json:"value"`
	Timestamp 	interface{}     	`json:"timestamp"`
}

type NodeActuatorDataMQTT struct {
	NodeID   	string          	`json:"node_id"`
	ActuatorID  string          	`json:"actuator_id"`
	Action   	string       		`json:"action"`
	Value     	float64       		`json:"value"`
}

func ConnectDBTS() *pgx.Conn {
	return utils.ConnectTSDB()
}

func DropNodeData() {
	db := ConnectDBTS()
	defer db.Close(context.Background())
	dropSensorDataTable := fmt.Sprintf("DROP TABLE IF EXISTS %s", "sensor_data")
	_, err := db.Exec(context.Background(), dropSensorDataTable)
	if err != nil {
		panic(err)
	}

	dropActuatorDataTable := fmt.Sprintf("DROP TABLE IF EXISTS %s", "actuator_data")
	_, err = db.Exec(context.Background(), dropActuatorDataTable)
	if err != nil {
		panic(err)
	}
}

func MigrateNodeData() {
	db := ConnectDBTS()
	defer db.Close(context.Background())
	initSensorDataTable := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s) TIMESTAMP(timestamp) PARTITION BY DAY;", 
		"sensor_data", 
		"node_id SYMBOL INDEX, sensor_id SYMBOL INDEX, value SYMBOL INDEX, timestamp TIMESTAMP")
	_, err := db.Exec(context.TODO(), initSensorDataTable)
	if err != nil {
		panic(err)
	}

	initActuatorDataTable := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s) TIMESTAMP(timestamp) PARTITION BY DAY;", 
		"actuator_data", 
		"node_id SYMBOL INDEX, actuator_id SYMBOL INDEX, action SYMBOL INDEX, value SYMBOL INDEX, timestamp TIMESTAMP")
	_, err = db.Exec(context.TODO(), initActuatorDataTable)
	if err != nil {
		panic(err)
	}
}

func (n *NodeSensorData) CreateNodeSensorData(db *pgx.Conn) error {
	query := fmt.Sprintf("INSERT INTO %s (node_id, sensor_id, value, timestamp) VALUES ($1, $2, $3, $4)", "sensor_data")
	
	_, err := db.Prepare(context.Background(), "insert_sensor_data", query)
	if err != nil {
		log.Printf("Error preparing insert sensor data: %v", err)
		return err
	}

	_, err = db.Exec(context.Background(), "insert_sensor_data", n.NodeID, n.SensorID, n.Value, n.Timestamp)
	if err != nil {
		log.Printf("Error inserting sensor data: %v", err)
		return err
	}

	defer db.Close(context.Background())
	log.Printf("Inserting sensor data: %v success", n)
	return nil
}

func (n *NodeActuatorData) CreateNodeActuatorData(db *pgx.Conn) error {
	query := fmt.Sprintf("INSERT INTO %s (node_id, actuator_id, action, value, timestamp) VALUES ($1, $2, $3, $4, $5)", "actuator_data")
	
	_, err := db.Prepare(context.Background(), "insert_actuator_data", query)
	if err != nil {
		log.Printf("Error preparing insert actuator data: %v", err)
		return err
	}

	_, err = db.Exec(context.Background(), "insert_actuator_data", n.NodeID, n.ActuatorID, n.Action, n.Value, n.Timestamp)
	if err != nil {
		log.Printf("Error inserting actuator data: %v", err)
		return err
	}

	defer db.Close(context.Background())
	log.Printf("Inserting actuator data: %v success", n)
	return nil

}

func (n *NodeSensorData) CheckDuplicateSensorData(db *pgx.Conn) error {
	query := fmt.Sprintf("SELECT value, timestamp FROM %s WHERE node_id = '%s' AND sensor_id = '%s' AND timestamp = '%s' LIMIT 1", "sensor_data", n.NodeID, n.SensorID, n.Timestamp)

	var dbData NodeSensorData
	data := db.QueryRow(context.TODO(), query)
	if err := data.Scan(&dbData.Value, &dbData.Timestamp); err != nil {
		log.Println(err)
		if err != pgx.ErrNoRows {
			log.Printf("Error querying sensor data: %v", err)
			return err
		}
		if err == pgx.ErrNoRows {
			log.Printf("No duplicate sensor data found")
			return nil
		}
	}

	log.Println("Duplicate sensor data found")
	return fiber.ErrConflict
}

func (n *NodeActuatorData) CheckDuplicateActuatorData(db *pgx.Conn) error {
	query := fmt.Sprintf("SELECT action, value, timestamp FROM %s WHERE node_id = '%s' AND actuator_id = '%s' AND timestamp = '%s' LIMIT 1", "actuator_data", n.NodeID, n.ActuatorID, n.Timestamp)

	var dbData NodeActuatorData
	data := db.QueryRow(context.TODO(), query)
	if err := data.Scan(&dbData.Action, &dbData.Value, &dbData.Timestamp); err != nil {
		log.Println(err)
		if err != pgx.ErrNoRows {
			return err
		}
		if err == pgx.ErrNoRows {
			log.Printf("No duplicate actuator data found")
			return nil
		}
	}

	log.Println("Duplicate actuator data found")
	return fiber.ErrConflict
}

func GetSensorData(db *pgx.Conn, nodeID, sensorID, fromTs, toTs, orderBy, limit string) ([]NodeSensorDataGet, error) {
	query := fmt.Sprintf("SELECT value, timestamp FROM %s WHERE node_id = '%s' AND sensor_id = '%s' AND timestamp BETWEEN '%s' AND '%s' ORDER BY timestamp %s LIMIT %s", "sensor_data", nodeID, sensorID, fromTs, toTs, orderBy, limit)

	var dbData []NodeSensorDataGet
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var data NodeSensorDataGet
		rows.Scan(&data.Value, &data.Timestamp)
		dbData = append(dbData, data)
	}

	return dbData, nil
}

func (n *NodeSensorData) GetLastSensorData(db *pgx.Conn) (NodeSensorDataGet, error) {
	query := fmt.Sprintf("SELECT value, timestamp FROM %s WHERE node_id = '%s' AND sensor_id = '%s' ORDER BY timestamp DESC LIMIT 1", "sensor_data", n.NodeID, n.SensorID)

	var dbData NodeSensorDataGet
	data := db.QueryRow(context.Background(), query)
	if err := data.Scan(&dbData.Value, &dbData.Timestamp); err != nil {
		log.Printf("Error querying sensor data: %v", err)
		return dbData, err
	}

	return dbData, nil
}

func (n *NodeActuatorData) GetLastActuatorData(db *pgx.Conn) (NodeActuatorData, error) {
	query := fmt.Sprintf("SELECT action, value, timestamp FROM %s WHERE node_id = '%s' AND actuator_id = '%s' ORDER BY timestamp DESC LIMIT 1", "actuator_data", n.NodeID, n.ActuatorID)

	var dbData NodeActuatorData
	data := db.QueryRow(context.Background(), query)
	if err := data.Scan(&dbData.Action, &dbData.Value, &dbData.Timestamp); err != nil {
		log.Printf("Error querying actuator data: %v", err)
		return dbData, err
	}

	return dbData, nil
}
