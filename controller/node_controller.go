package controller

import (
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NodeController struct {
	DB *gorm.DB
}

type NodeResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Coordinate []string  `json:"coordinate"`
}

func (c *NodeController) RegisterNewNode(ctx *fiber.Ctx) error {

	var node model.Node
	if err := ctx.BodyParser(&node); err != nil {
		return err
	}
	
	if err := node.CreateNode(c.DB); err != nil {
		return err
	}
	
	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    node,
	})
}

func (c *NodeController) GetAllNodes(ctx *fiber.Ctx) error {
	var nodes []model.Node
	if err := c.DB.Find(&nodes).Error; err != nil {
		return err
	}

	var nodeResponses []NodeResponse
	for _, node := range nodes {
		nodeResponses = append(nodeResponses, NodeResponse{
			ID:        node.ID,
			Name:      node.Name,
			Coordinate: []string{node.Latitude, node.Longitude},
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    nodeResponses,
	})
}

func (c *NodeController) DeleteNodeByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var node model.Node
	if err := c.DB.First(&node, id).Error; err != nil {
		return err
	}

	if err := c.DB.Delete(&node).Error; err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}
