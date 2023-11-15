package controller

import (
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ActuatorController struct {
	DB *gorm.DB
}

type ActuatorResponse struct {
	ID 				uint   `json:"id"`
	NodeID 			uint   `json:"node_id"`
	Name 			string `json:"name"`
	ActuatorType 	int    `json:"actuator_type"`
	Unit 			string `json:"unit"`
	Interval 		int    `json:"interval"`
}

func NewActuatorController(db *gorm.DB) *ActuatorController {
	return &ActuatorController{
		DB: db,
	}
}

func (c *ActuatorController) CreateActuator(ctx *fiber.Ctx) error {
	var actuator model.Actuator
	if err := ctx.BodyParser(&actuator); err != nil {
		return err
	}

	if err := actuator.CreateActuator(c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    actuator,
	})
}

func (c *ActuatorController) GetAllActuators(ctx *fiber.Ctx) error {
	var actuators []model.Actuator
	if err := c.DB.Find(&actuators).Error; err != nil {
		return err
	}

	var actuatorResponses []ActuatorResponse
	for _, actuator := range actuators {
		actuatorResponses = append(actuatorResponses, ActuatorResponse{
			ID:        actuator.ID,
			NodeID:   actuator.NodeID,
			Name:      actuator.Name,
			ActuatorType: actuator.ActuatorType,
			Unit: actuator.Unit,
			Interval: actuator.Interval,
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    actuatorResponses,
	})
}

func (c *ActuatorController) GetActuatorByID(ctx *fiber.Ctx) error {
	var actuator model.Actuator
	if err := c.DB.First(&actuator, ctx.Params("id")).Error; err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    actuator,
	})
}

func (c *ActuatorController) UpdateActuatorByID(ctx *fiber.Ctx) error {
	var actuator model.Actuator
	if err := c.DB.First(&actuator, ctx.Params("id")).Error; err != nil {
		return err
	}

	if err := ctx.BodyParser(&actuator); err != nil {
		return err
	}

	if err := actuator.UpdateActuator(c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    actuator,
	})
}

func (c *ActuatorController) DeleteActuatorByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var actuator model.Actuator
	if err := c.DB.First(&actuator, id).Error; err != nil {
		return err
	}

	if err := actuator.DeleteActuator(c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}