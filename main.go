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

	model.MigrateRole(db)
	model.MigrateUser(db)
	model.MigratePermission(db)
	model.MigrateNotification(db)
	model.MigrateNode(db)
	model.MigrateSensor(db)

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

	notif := app.Group("/notifpref")
	notif.Use(middleware.IsAuthenticated)
	notifController := &controller.NotifController{DB: db}
	notif.Put("/:userID", notifController.UpdateNotifPreference)

	node := app.Group("/nodes")
	node.Use(middleware.IsAuthenticated)
	nodeController := &controller.NodeController{DB: db}
	node.Post("/", nodeController.RegisterNewNode)
	node.Get("/", nodeController.GetAllNodes)

	sensor := app.Group("/sensors")
	sensor.Use(middleware.IsAuthenticated)
	sensorController := &controller.SensorController{DB: db}
	sensor.Post("/", sensorController.AddNewSensor)
	sensor.Get("/", sensorController.GetAllSensors)

	listenAddr := fmt.Sprintf("%s:%s", APP_HOST, APP_PORT)
    log.Fatal(app.Listen(listenAddr))

}
