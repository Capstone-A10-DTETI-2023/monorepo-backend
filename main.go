package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/controller"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/middleware"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/model"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var (
		APP_HOST = os.Getenv("APP_HOST")
		APP_PORT = os.Getenv("APP_PORT")
	)
	
    app := fiber.New()
	db := utils.ConnectDB()

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	
	app.Get("/ping", func (c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

	model.MigrateUser(db)
	model.MigrateRole(db)
	model.MigratePermission(db)
	user := app.Group("/users")
	user.Use(middleware.IsAuthenticated)
	userController := &controller.UserController{DB: db}
	user.Post("/", userController.CreateUser)
	user.Get("/", userController.GetAllUsers)
	user.Get("/:id", userController.GetUserByID)
	user.Put("/:id", userController.UpdateUserByID)

	auth := app.Group("/auth")
	authController := &controller.AuthController{DB: db}
	auth.Post("/login", authController.Login)

	role := app.Group("/roles")
	role.Use(middleware.IsAuthenticated)
	roleController := &controller.RoleController{DB: db}
	role.Get("/", roleController.GetAllRole)
	role.Post("/", roleController.CreateRole)

	listenAddr := fmt.Sprintf("%s:%s", APP_HOST, APP_PORT)
    log.Fatal(app.Listen(listenAddr))

}
