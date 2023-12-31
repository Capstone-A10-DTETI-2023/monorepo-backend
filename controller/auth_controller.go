package controller

import (
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

	message := "Halo, " + user.Name + "! Password Anda telah direset. Password baru Anda adalah: " + newPassword + "  Silahkan login dan ubah password Anda."

	if err := utils.SendWAMessage(user.Phone_Num, message, "0"); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name: "token",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {

	var user model.User
	if err := ctx.BodyParser(&user); err != nil {
		return err
	}

	newPassword, _ := utils.GenerateRandomPassword(12)
	user.Password = newPassword
	user.RoleID = 3

	if err := user.HashPassword(); err != nil {
		return err
	}

	if err := user.CreateUser(c.DB); err != nil {
		return err
	}

	message := "Halo, " + user.Name + "! Akun Anda telah dibuat. Password login Anda adalah: " + newPassword + "  Silahkan login dan ubah password Anda."

	if err := utils.SendWAMessage(user.Phone_Num, message, "0"); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}

