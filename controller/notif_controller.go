package controller

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

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

func (c *NotifController) GetNotifPreference(ctx *fiber.Ctx) error {
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
	if err := c.DB.Where("user_id = ?", userID).First(&notifPref).Error; err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    notifPref,
	})
}

func (c *NotifController) GetAllNotifPref(ctx *fiber.Ctx) error {
	isAdmin := middleware.IsAdmin(ctx)
	if !isAdmin {
		return fiber.ErrUnauthorized
	}

	var notifPref []model.Notification
	if err := c.DB.Find(&notifPref).Error; err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    notifPref,
	})
}

func (c *NotifController) SendWhatsAppNotification(ctx *fiber.Ctx) error {
	tokenWA := os.Getenv("TOKEN_WA")

	isAdmin := middleware.IsAdmin(ctx)
	if !isAdmin {
		return fiber.ErrUnauthorized
	}

	var request model.WhatsAppNotificationRequest
	if err := ctx.BodyParser(&request); err != nil {
		return err
	}

	var user model.User
	if err := c.DB.Where("id = ?", request.UserID).First(&user).Error; err != nil {
		return err
	}

	data := url.Values{}
	data.Set("target", user.Phone_Num)
	data.Set("message", request.Message)
	data.Set("schedule", request.Schedule)

	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, "https://api.fonnte.com/send", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", tokenWA)

	result, err := client.Do(r)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}
