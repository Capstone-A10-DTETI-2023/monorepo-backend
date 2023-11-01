package controller

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/middleware"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

type LoginRequest struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

type ResetPassRequest struct {
	Username 		string `json:"username"`
	Phone_Num 	string `json:"phone_num"`
}

type LoginResponse struct {
	AccessToken 	string `json:"access_token"`
	ID 				uint   `json:"id"`
	Role_ID			uint   `json:"role_id"`
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var loginRequest LoginRequest
	if err := ctx.BodyParser(&loginRequest); err != nil {
		return err
	}

	var user model.User
	if err := c.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		return err
	}

	if passwordMatch := model.CompareHashAndPassword(user.Password, loginRequest.Password); !passwordMatch {
		return fiber.ErrUnauthorized
	}

	token, err := middleware.GenerateJWT(user.ID, user.RoleID)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	ctx.Cookie(&fiber.Cookie{
		Name: "token",
		Value: token,
		Expires: time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data": LoginResponse{
			AccessToken:  token,
			ID: 		   user.ID,
			Role_ID: 	   user.RoleID,
		},
	})
}

func (c *AuthController) ResetPassword(ctx *fiber.Ctx) error {
	tokenWA := os.Getenv("TOKEN_WA")

	var resetPassRequest ResetPassRequest
	if err := ctx.BodyParser(&resetPassRequest); err != nil {
		return err
	}
	phoneNum := resetPassRequest.Phone_Num
	username := resetPassRequest.Username

	var user model.User
	if err := c.DB.Where("name = ?", username).Where("phone_num = ?", phoneNum).First(&user).Error; err != nil {
		return err
	}

	newPassword, _:= utils.GenerateRandomPassword(12)
	user.Password = newPassword
	if err := user.HashPassword(); err != nil {
		return err
	}

	if err := user.UpdateUser(c.DB); err != nil {
		return err
	}

	request := model.WhatsAppNotificationRequest{
		Message: "Halo, " + user.Name + "! Password Anda telah direset. Password baru Anda adalah: " + newPassword + "  Silahkan login dan ubah password Anda.",
		Schedule: "0",
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
