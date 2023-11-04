package controller

import (
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SysSettingController struct {
	DB *gorm.DB
}

type NodePressureRefResponse struct {
	ID        uint   `json:"id"`
	NodeID    uint   `json:"node_id"`
	Pressure  float64 `json:"pressure"`
}

func NewSysSettingController(db *gorm.DB) *SysSettingController {
	return &SysSettingController{
		DB: db,
	}
}

func (c *SysSettingController) SetNodePressureRef(ctx *fiber.Ctx) error {
	var nodePresRef model.NodePressureRef

	nodeID := ctx.Params("node_id")
	if err := c.DB.First(&nodePresRef, "node_id = ?", nodeID).Error; err != nil {
		return err
	}

	if err := ctx.BodyParser(&nodePresRef); err != nil {
		return err
	}

	if err := nodePresRef.UpdateNodePressureRef(c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    nodePresRef,
	})
}

func (c *SysSettingController) GetAllNodePressureRef(ctx *fiber.Ctx) error{
	var nodePresRef []model.NodePressureRef

	if err := c.DB.Find(&nodePresRef).Error; err != nil {
		return err
	}

	var nodePresRefRes []NodePressureRefResponse
	for _, nodePR := range nodePresRef {
		nodePresRefRes = append(nodePresRefRes, NodePressureRefResponse{
			ID: nodePR.ID,
			NodeID: nodePR.NodeID,
			Pressure: nodePR.Pressure,
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    nodePresRefRes,
	})
}