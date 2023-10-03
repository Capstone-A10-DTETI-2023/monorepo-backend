package controller

import (
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SensorController struct {
	DB *gorm.DB
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

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    sensors,
	})
}
