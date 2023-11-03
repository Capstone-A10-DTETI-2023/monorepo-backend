package controller

import (
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SensorController struct {
	DB *gorm.DB
}

type SensorResponse struct {
	ID 			uint   `json:"id"`
	Node_ID		uint   `json:"node_id"`
	Name 			string `json:"name"`
	Unit 			string `json:"unit"`
}

func (c *SensorController) AddNewSensor(ctx *fiber.Ctx) error {
	var sensor model.Sensor
	if err := ctx.BodyParser(&sensor); err != nil {
		return err
	}

	if err := sensor.CreateSensor(c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    sensor,
	})
}

func (c *SensorController) GetAllSensors(ctx *fiber.Ctx) error {
	var sensors []model.Sensor
	if err := c.DB.Find(&sensors).Error; err != nil {
		return err
	}

	var sensorRes []SensorResponse
	for _, sensor := range sensors {
		sensorRes = append(sensorRes, SensorResponse{
			ID: sensor.ID,
			Node_ID: sensor.NodeID,
			Name: sensor.Name,
			Unit: sensor.Unit,
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    sensorRes,
	})
}

func (c *SensorController) DeleteSensor(ctx *fiber.Ctx) error {
	sensorID := ctx.Params("id")

	var sensor model.Sensor
	if err := c.DB.Where("id = ?", sensorID).First(&sensor).Error; err != nil {
		return err
	}

	if err := sensor.DeleteSensor(c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}

func (c *SensorController) UpdateSensorByID (ctx *fiber.Ctx) error {
	sensorID := ctx.Params("id")

	var sensor model.Sensor
	if err := c.DB.Where("id = ?", sensorID).First(&sensor).Error; err != nil {
		return err
	}

	if err := ctx.BodyParser(&sensor); err != nil {
		return err
	}

	if err := sensor.UpdateSensor(c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}

func (c *SensorController) GetSensorByNodeID(ctx *fiber.Ctx) error {
	nodeID := ctx.Params("node_id")

	var sensor []model.Sensor
	if err := c.DB.Where("node_id = ?", nodeID).Find(&sensor).Error; err != nil {
		return err
	}

	var sensorRes []SensorResponse
	for _, sensor := range sensor {
		sensorRes = append(sensorRes, SensorResponse{
			ID: sensor.ID,
			Node_ID: sensor.NodeID,
			Name: sensor.Name,
			Unit: sensor.Unit,
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    sensorRes,
	})
}