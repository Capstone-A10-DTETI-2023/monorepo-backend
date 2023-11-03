package controller

import (
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NodeController struct {
	DB *gorm.DB
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

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    nodes,
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
