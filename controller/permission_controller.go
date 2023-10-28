package controller

import (
	"strconv"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/middleware"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PermissionController struct {
	DB *gorm.DB
}

type PermissionResponse struct {
	RoleID 					uint   `json:"role_id"`
	Read_Realtime_Data 		bool   `json:"read_realtime_data"`
	Read_Historical_Data 	bool   `json:"read_historical_data"`
	Change_Actuator 		bool   `json:"change_actuator"`
	User_Management 		bool   `json:"user_management"`
	Node_Management 		bool   `json:"node_management"`
}

func (c *PermissionController) GetAllPermission(ctx *fiber.Ctx) error {
	if isAdmin := middleware.IsAdmin(ctx); !isAdmin {
		return fiber.ErrUnauthorized
	}

	var perms []model.Permission

	if err := c.DB.Find(&perms).Error; err != nil {
		return err
	}

	var permResponses []PermissionResponse
	for _, perm := range perms {
		permResponses = append(permResponses, PermissionResponse{
			RoleID: perm.RoleID,
			Read_Realtime_Data: perm.Read_Realtime_Data,
			Read_Historical_Data: perm.Read_Historical_Data,
			Change_Actuator: perm.Change_Actuator,
			User_Management: perm.User_Management,
			Node_Management: perm.Node_Management,
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":   permResponses,
	})
}

func (c *PermissionController) GetPermissionsByID(ctx *fiber.Ctx) error {
	roleID, _ := strconv.Atoi(ctx.Params("id"))

	var perms model.Permission
	if err := c.DB.Where("role_id = ?", roleID).First(&perms).Error; err != nil {
		return err
	}

	permsRes := PermissionResponse{
		RoleID: perms.RoleID,
		Read_Realtime_Data: perms.Read_Realtime_Data,
		Read_Historical_Data: perms.Read_Historical_Data,
		Change_Actuator: perms.Change_Actuator,
		User_Management: perms.User_Management,
		Node_Management: perms.Node_Management,
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    permsRes,
	})
}

func (c *PermissionController) UpdatePermission(ctx *fiber.Ctx) error {
	if isAdmin := middleware.IsAdmin(ctx); !isAdmin {
		return fiber.ErrUnauthorized
	}

	roleID, _ := strconv.Atoi(ctx.Params("id"))

	var perms model.Permission
	if err := c.DB.Where("role_id = ?", roleID).First(&perms).Error; err != nil {
		return err
	}

	if err := ctx.BodyParser(&perms); err != nil {
		return err
	}

	if err := perms.UpdatePermission(c.DB); err != nil {
		return err
	}

	permsRes := PermissionResponse{
		RoleID: perms.RoleID,
		Read_Realtime_Data: perms.Read_Realtime_Data,
		Read_Historical_Data: perms.Read_Historical_Data,
		Change_Actuator: perms.Change_Actuator,
		User_Management: perms.User_Management,
		Node_Management: perms.Node_Management,
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    permsRes,
	})
}
