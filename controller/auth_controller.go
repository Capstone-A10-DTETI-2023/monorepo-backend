package controller

import (
	"time"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/middleware"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
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

	token, err := middleware.GenerateJWT(user.ID, user.Role_ID)

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
			Role_ID: 	   user.Role_ID,
		},
	})
}