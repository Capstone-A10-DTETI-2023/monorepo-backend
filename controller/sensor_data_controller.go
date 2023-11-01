package controller

import (
	"time"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
)

func InsertDataSensor(ctx *fiber.Ctx) error {
	var sensorData model.NodeSensorData
	db := model.ConnectDBTS()

	if sensorData.Timestamp == "" {
		sensorData.Timestamp = time.Now().String()
	}

	var timestamp time.Time
	bufferTime := sensorData.Timestamp.(string)
	timestamp, err := time.Parse("2006-01-02 15:04:05", bufferTime)
	timestampStr := timestamp.Format(time.RFC3339)
	if err != nil {
		return err
	}

	if sensorData.NodeID == "" || sensorData.SensorID == "" || sensorData.Value == "" {
		return fiber.ErrBadRequest
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

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data": sensorData,
	})
}

func GetSensorData(ctx *fiber.Ctx) error {
	db := model.ConnectDBTS()

	nodeID := ctx.Query("node_id")
    sensorID := ctx.Query("sensor_id")
	fromTS := ctx.Query("from")
	toTS := ctx.Query("to")
	order := ctx.Query("order_by")
	limit := ctx.Query("limit")

	var timestamp time.Time
	bufferFromTime := fromTS
	timestamp, _ = time.Parse("2006-01-02 15:04:05", bufferFromTime)
	fromTimeStr := timestamp.Format(time.RFC3339)

	bufferToTime := toTS
	timestamp, _ = time.Parse("2006-01-02 15:04:05", bufferToTime)
	toTimeStr := timestamp.Format(time.RFC3339)

	if nodeID == "" || sensorID == "" {
		return fiber.ErrBadRequest
	}

	sensorDataDB, err := model.GetSensorData(db, nodeID, sensorID, fromTimeStr, toTimeStr, order, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data": fiber.Map{
			"id_node": nodeID,
			"id_sensor": sensorID,
			"sensor_data": sensorDataDB,
		},
	})
}