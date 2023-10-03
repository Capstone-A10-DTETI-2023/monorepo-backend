package controller

import (
	"strconv"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/middleware"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NotifController struct {
	DB *gorm.DB
}

func (c *NotifController) UpdateNotifPreference(ctx *fiber.Ctx) error {
	var isAllowed bool

	token := ctx.Cookies("token")
	claims, err := middleware.ParseJWT(token)
	userID, _ := strconv.Atoi(ctx.Params("userID"))
	if err != nil || claims.ExpiresAt <= 0 || (claims.ID != uint(userID) && claims.Role_ID != 1) {
		isAllowed = false
	} else {
		isAllowed = true
	}

	if !isAllowed {
		return fiber.ErrUnauthorized
	}
	
	var notifPref model.Notification
	if err := ctx.BodyParser(&notifPref); err != nil {
		return err
	}

	notifPref.UserID = uint(userID)
	notifPref.ID = uint(userID)
	if err := notifPref.UpdateNotification(c.DB); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    notifPref,
	})
}