package controller

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/producer"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ActuatorDataController struct {
	DB *gorm.DB
}

func NewActuatorDataController(db *gorm.DB) *ActuatorDataController {
	return &ActuatorDataController{
		DB: db,
	}
}

func (c *ActuatorDataController) InsertDataActuator(ctx *fiber.Ctx) error {
	var actuatorData model.NodeActuatorData
	actuatorData.Timestamp = time.Now().Format(time.RFC3339)
	if err := ctx.BodyParser(&actuatorData); err != nil {
		return err
	}
	
	if err := InsertDataActuatorToDB(actuatorData, c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    actuatorData,
	})
}

func InsertDataActuatorToDB(actuatorData model.NodeActuatorData, dbData *gorm.DB) error {
	dbTs := model.ConnectDBTS()
	var actuator model.Actuator

	if actuatorData.Action == "" {
		return Error{"Missing required field"}
	}

	if err := dbData.Where("id = ?", actuatorData.ActuatorID).First(&actuator).Error; err != nil {
		return err
	}

	reqNodeId, _ := strconv.Atoi(actuatorData.NodeID)
	if uint(reqNodeId) != actuator.NodeID {
		return Error{"Actuator does not belong to Node ID provided"}
	}

	if err := actuatorData.CheckDuplicateActuatorData(dbTs); err != nil {
		return Error{"Duplicate data"}
	}

	if err := actuatorData.CreateNodeActuatorData(dbTs); err != nil {
		return err
	}

	var mqttModel model.NodeActuatorDataMQTT
	mqttModel.NodeID = actuatorData.NodeID
	mqttModel.ActuatorID = actuatorData.ActuatorID
	mqttModel.Action = actuatorData.Action
	mqttModel.Value, _ = strconv.ParseFloat(actuatorData.Value, 64)

	mqttJSON, _ := json.Marshal(mqttModel)
	mqttCli := producer.ConnectProducerMQTT()
	producer.PublishMQTT(mqttCli, "actuatorData", 0, true, string(mqttJSON))

	return nil
}	

func (c *ActuatorDataController) GetLastActuatorData (ctx *fiber.Ctx) error {
	db := model.ConnectDBTS()

	nodeID := ctx.Query("node_id")
	actuatorID := ctx.Query("actuator_id")

	var actuatorData model.NodeActuatorData
	actuatorData.NodeID = nodeID
	actuatorData.ActuatorID = actuatorID
	actuatorData, err := actuatorData.GetLastActuatorData(db)
	if err != nil {
		return err
	}
	actuatorData.NodeID = nodeID
	actuatorData.ActuatorID = actuatorID

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    actuatorData,
	})
}
