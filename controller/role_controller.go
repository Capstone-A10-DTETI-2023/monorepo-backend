package controller

import (
	"strconv"

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

	permsRole := model.Permission{
		RoleID: role.ID,
		Read_Realtime_Data: true,
		Read_Historical_Data: false,
		Change_Actuator: false,
		User_Management: false,
		Node_Management: false,
	}

	if err := permsRole.CreatePermission(c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    role,
	})
}

func (c *RoleController) UpdateRole(ctx *fiber.Ctx) error{
	if isAdmin := middleware.IsAdmin(ctx); !isAdmin {
		return fiber.ErrUnauthorized
	}

	roleID, _ := strconv.Atoi(ctx.Params("id"))

	var role model.Role
	if err := c.DB.Where("id = ?", roleID).First(&role).Error; err != nil {
		return err
	}

	if err := ctx.BodyParser(&role); err != nil {
		return err
	}

	if err := role.UpdateRole(c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    role,
	})
}
