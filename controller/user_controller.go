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
	Role_Name 	string `json:"role_name"`
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
		UserID: user.ID,
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
		Role_ID:  	user.RoleID,
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

	var users []UserResponse
	rows, err := c.DB.Table("users").Select("users.id, users.name, users.email, users.phone_num, users.role_id, roles.name").Joins("left join roles on roles.id = users.role_id").Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var user UserResponse
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone_Num, &user.Role_ID, &user.Role_Name)
		users = append(users, user)
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    users,
	})
}

func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	if isAdmin := middleware.IsAdmin(ctx); !isAdmin {
		return fiber.ErrUnauthorized
	}

	userID := ctx.Params("id")

	var user UserResponse
	row, err := c.DB.Table("users").Select("users.id, users.name, users.email, users.phone_num, users.role_id, roles.name").Where("users.id = ?", userID).Joins("left join roles on roles.id = users.role_id").Where("users.id = ?", userID).Rows()
	if err != nil {
		return err
	}
	for row.Next() {
		row.Scan(&user.ID, &user.Name, &user.Email, &user.Phone_Num, &user.Role_ID, &user.Role_Name)
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    user,
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
	
	if claims.Role_ID != 1 && user.RoleID == 1 {
		return fiber.ErrUnauthorized
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	if err := user.UpdateUser(c.DB); err != nil {
		return err
	}

	var userResponses []UserResponse
	userResponses = append(userResponses, UserResponse{
		ID:        	user.ID,
		Role_ID:  	user.RoleID,
		Name:      	user.Name,
		Email:     	user.Email,
		Phone_Num: 	user.Phone_Num,
	})

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    userResponses,
	})
}
