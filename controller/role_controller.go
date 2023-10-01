package controller

import (
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/middleware"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RoleController struct {
	DB *gorm.DB
}

func (c *RoleController) GetAllRole(ctx *fiber.Ctx) error {
	if isAdmin := middleware.IsAdmin(ctx); !isAdmin {
		return fiber.ErrUnauthorized
	}
	
	var roles []model.Role

	if err := c.DB.Find(&roles).Error; err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    roles,
	})
}

func (c *RoleController) CreateRole(ctx *fiber.Ctx) error{
	if isAdmin := middleware.IsAdmin(ctx); !isAdmin {
		return fiber.ErrUnauthorized
	}

	var role model.Role
	if err := ctx.BodyParser(&role); err != nil {
		return err
	}

	if err := role.CreateRole(c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    role,
	})
}