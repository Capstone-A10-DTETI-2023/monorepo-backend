package controller

import (
	"time"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
)

func InsertDataSensor(ctx *fiber.Ctx) error {
	var sensorData model.NodeSensorData
	db := model.ConnectDBTS()

	if err := ctx.BodyParser(&sensorData); err != nil {
		return err
	}

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