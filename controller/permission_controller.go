package controller

import (
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/middleware"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PermissionController struct {
	DB *gorm.DB
}

func (c *PermissionController) GetAllPermission(ctx *fiber.Ctx) error {
	if isAdmin := middleware.IsAdmin(ctx); !isAdmin {
		return fiber.ErrUnauthorized
	}

	var perms []model.Permission

	if err := c.DB.Find(&perms).Error; err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":   perms,
	})
}