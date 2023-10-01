package controller

import (
	"strconv"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/middleware"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

type UserResponse struct {
	ID 	  		uint   `json:"id"`
	Role_ID		uint   `json:"role_id"`
	Name 	  	string `json:"name"`
	Email 	  	string `json:"email"`
	Phone_Num 	string `json:"phone_num"`
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	if isAdmin := middleware.IsAdmin(ctx); !isAdmin {
		return fiber.ErrUnauthorized
	}
	
	var user model.User
	var notifPref model.Notification
	if err := ctx.BodyParser(&user); err != nil {
		return err
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	if err := user.CreateUser(c.DB); err != nil {
		return err
	}

	notifPref = model.Notification{
		User_ID: user.ID,
		Email: false,
		Whatsapp: false,
		Firebase: false,
	}
	
	if err := notifPref.CreateNotification(c.DB); err != nil {
		return err
	}

	var userResponses []UserResponse
	userResponses = append(userResponses, UserResponse{
		ID:        	user.ID,
		Role_ID:  	user.Role_ID,
		Name:      	user.Name,
		Email:     	user.Email,
		Phone_Num: 	user.Phone_Num,
	})

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    userResponses,
	})
}

func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	if isAdmin := middleware.IsAdmin(ctx); !isAdmin {
		return fiber.ErrUnauthorized
	}

	var users []model.User

	if err := c.DB.Find(&users).Error; err != nil {
		return err
	}

	var userResponses []UserResponse
    for _, user := range users {
        userResponses = append(userResponses, UserResponse{
            ID:        	user.ID,
			Role_ID: 	user.Role_ID,
            Name:      	user.Name,
            Email:     	user.Email,
            Phone_Num: 	user.Phone_Num,
        })
    }

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    userResponses,
	})
}

func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	if isAdmin := middleware.IsAdmin(ctx); !isAdmin {
		return fiber.ErrUnauthorized
	}

	userID := ctx.Params("id")

	var user model.User
	if err := c.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	var userResponses []UserResponse
	userResponses = append(userResponses, UserResponse{
		ID:        	user.ID,
		Role_ID:  	user.Role_ID,
		Name:      	user.Name,
		Email:     	user.Email,
		Phone_Num: 	user.Phone_Num,
	})

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    userResponses,
	})
}

func (c *UserController) UpdateUserByID(ctx *fiber.Ctx) error {
	var isAllowed bool

	token := ctx.Cookies("token")
	claims, err := middleware.ParseJWT(token)
	userID, _ := strconv.Atoi(ctx.Params("id"))
	if err != nil || claims.ExpiresAt <= 0 || (claims.ID != uint(userID) && claims.Role_ID != 1) {
		isAllowed = false
	} else {
		isAllowed = true
	}
	

	if !isAllowed {
		return fiber.ErrUnauthorized
	}

	var user model.User
	if err := c.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	if err := ctx.BodyParser(&user); err != nil {
		return err
	}
	
	if claims.Role_ID != 1 && user.Role_ID == 1 {
		return fiber.ErrUnauthorized
	}

	if err := user.UpdateUser(c.DB); err != nil {
		return err
	}

	var userResponses []UserResponse
	userResponses = append(userResponses, UserResponse{
		ID:        	user.ID,
		Role_ID:  	user.Role_ID,
		Name:      	user.Name,
		Email:     	user.Email,
		Phone_Num: 	user.Phone_Num,
	})

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    userResponses,
	})
}
