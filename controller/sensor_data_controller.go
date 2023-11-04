package controller

import (
	"strconv"
	"time"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

type SensorDataController struct {
	DB *gorm.DB
}

func (c *SensorDataController) InsertDataSensor(ctx *fiber.Ctx) error {
	var sensorData model.NodeSensorData

	if err := ctx.BodyParser(&sensorData); err != nil {
		return err
	}

	if sensorData.Value == "" {
		return Error{"Missing required field"}
	}

	if err := InsertDataSensorToDB(sensorData, c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data": sensorData,
	})
}

func InsertDataSensorToDB(sensorData model.NodeSensorData, dbData *gorm.DB) error {
	dbTs := model.ConnectDBTS()
	var sensor model.Sensor

	if err := dbData.Where("id = ?", sensorData.SensorID).First(&sensor).Error; err != nil {
		return err
	}

	reqNodeId, _ := strconv.Atoi(sensorData.NodeID)
	if  uint(reqNodeId) != sensor.NodeID {
		return Error{"Sensor does not belong to Node ID provided"}
	}

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

	if err := sensorData.CheckDuplicateSensorData(dbTs); err != nil {
		return err
	}

	if err := sensorData.CreateNodeSensorData(dbTs); err != nil {
		return err
	}

	if sensor.Alarm {
		switch sensor.AlarmType {
		case 1:
			sensorVal, _ := strconv.ParseFloat(sensorData.Value, 64)
			if sensorVal < sensor.AlarmLow {
				var phoneNum []string
				rows, err := dbData.Table("users").Select("users.phone_num").Joins("left join notifications on notifications.user_id = users.id").Where("notifications.whatsapp = ?", true).Rows()
				if err != nil {
					return err
				}
				defer rows.Close()
				for rows.Next() {
					var phone string
					rows.Scan(&phone)
					phoneNum = append(phoneNum, phone)
				}

				var nodeName string
				if err := dbData.Table("nodes").Select("nodes.name").Where("nodes.id = ?", sensor.NodeID).Scan(&nodeName).Error; err != nil {
					return err
				}

				message := "Sensor " + sensor.Name + " on " + nodeName + " is below the threshold. Current value: " + strconv.FormatFloat(sensorVal, 'f', 2, 64) + " " + sensor.Unit
				for _, phone := range phoneNum {
					utils.SendWAMessage(phone, message, "0")
				}
			}

		case 2:
			sensorVal, _ := strconv.ParseFloat(sensorData.Value, 64)
			if sensorVal > sensor.AlarmHigh {
				var phoneNum []string
				rows, err := dbData.Table("users").Select("users.phone_num").Joins("left join notifications on notifications.user_id = users.id").Where("notifications.whatsapp = ?", true).Rows()
				if err != nil {
					return err
				}
				defer rows.Close()
				for rows.Next() {
					var phone string
					rows.Scan(&phone)
					phoneNum = append(phoneNum, phone)
				}

				var nodeName string
				if err := dbData.Table("nodes").Select("nodes.name").Where("nodes.id = ?", sensor.NodeID).Scan(&nodeName).Error; err != nil {
					return err
				}

				message := "Sensor " + sensor.Name + " on " + nodeName + " is above the threshold. Current value: " + strconv.FormatFloat(sensorVal, 'f', 2, 64) + " " + sensor.Unit
				for _, phone := range phoneNum {
					utils.SendWAMessage(phone, message, "0")
				}
			}

		case 3:
			sensorVal, _ := strconv.ParseFloat(sensorData.Value, 64)
			if sensorVal < sensor.AlarmLow || sensorVal > sensor.AlarmHigh {
				var phoneNum []string
				rows, err := dbData.Table("users").Select("users.phone_num").Joins("left join notifications on notifications.user_id = users.id").Where("notifications.whatsapp = ?", true).Rows()
				if err != nil {
					return err
				}
				defer rows.Close()
				for rows.Next() {
					var phone string
					rows.Scan(&phone)
					phoneNum = append(phoneNum, phone)
				}

				var nodeName string
				if err := dbData.Table("nodes").Select("nodes.name").Where("nodes.id = ?", sensor.NodeID).Scan(&nodeName).Error; err != nil {
					return err
				}

				message := "Sensor " + sensor.Name + " on " + nodeName + " is outside the predefined threshold. Current value: " + strconv.FormatFloat(sensorVal, 'f', 2, 64) + " " + sensor.Unit
				for _, phone := range phoneNum {
					utils.SendWAMessage(phone, message, "0")
				}
			}
		}
	}

	return nil
}

func (c *SensorDataController) GetSensorData(ctx *fiber.Ctx) error {
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

func (c *SensorDataController) GetLastSensorData(ctx *fiber.Ctx) error {
	db := model.ConnectDBTS()

	nodeID := ctx.Query("node_id")
	sensorID := ctx.Query("sensor_id")

	if nodeID == "" || sensorID == "" {
		return fiber.ErrBadRequest
	}

	sensorData := model.NodeSensorData{
		NodeID: nodeID,
		SensorID: sensorID,
	}

	sensorDataDB, err := sensorData.GetLastSensorData(db)
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