package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/consumer"
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
	if len(os.Args) < 2 {
		log.Fatal("Missing argument. Please specify either 'server', 'bootstrap', or 'migrate'")
	}

	runArgs := os.Args[1]
	switch runArgs {
	case "server":
		migrate()
		server()
	case "bootstrap":
		migrate()
		bootstrap()
	case "migrate":
		migrate()
	default:
		log.Fatal("Invalid argument")
	}
}

func server() {
	var (
		_appHost = os.Getenv("APP_HOST")
		_appPort = os.Getenv("APP_PORT")
	)
    app := fiber.New()

	// Initialize DB Connection with GORM to PostgreSQL
	db := utils.ConnectDB()

	// Initialize DB Connection with PGX to QuestDB
	dbTs := utils.ConnectTSDB()
	if err := dbTs.Ping(context.Background()); err != nil {
		panic(err)
	}
	defer dbTs.Close(context.Background())
	log.Printf("Connected to TSDB QuestDB")

	// Initialize MQTT Connection
	go func(){
		mqtt := consumer.ConnectMQTT()
		defer mqtt.Disconnect(250)
		consumer.SubMQTT(mqtt, os.Getenv("MQTT_TOPIC_CONSUMER"), 0)
	}()

	// Initialize Fiber Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))
	
	// Initialize Fiber Routes
	app.Get("/ping", func (c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })
	
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
	role.Put("/:id", roleController.UpdateRole)

	permissions := app.Group("/permissions")
	permissions.Use(middleware.IsAuthenticated)
	permissionController := &controller.PermissionController{DB: db}
	permissions.Get("/", permissionController.GetAllPermission)

	notif := app.Group("/notif")
	notif.Use(middleware.IsAuthenticated)
	notifController := &controller.NotifController{DB: db}
	notif.Get("/user", notifController.GetAllNotifPref)
	notif.Put("/user/:userID", notifController.UpdateNotifPreference)
	notif.Get("/user/:userID", notifController.GetNotifPreference)
	notif.Post("/testing/wa", notifController.SendWhatsAppNotification)

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

	sensorData := app.Group("/tsdata")
	sensorData.Use(middleware.IsAuthenticated)
	sensorData.Post("/sensor", controller.InsertDataSensor)

	// Start Fiber App
	listenAddr := fmt.Sprintf("%s:%s", _appHost, _appPort)
    log.Fatal(app.Listen(listenAddr))
}

func migrate() error {
	// Initialize DB Connection with GORM to PostgreSQL
	db := utils.ConnectDB()

	// Initialize DB Migration
	model.MigrateRole(db)
	model.MigrateUser(db)
	model.MigratePermission(db)
	model.MigrateNotification(db)
	model.MigrateNode(db)
	model.MigrateSensor(db)
	model.MigrateNodeData()

	log.Println("Migration completed")
	return nil
}

func bootstrap() error {
	// Initialize DB Connection with GORM to PostgreSQL
	db := utils.ConnectDB()

	// Initialize DB Bootstrap
	model.BootstrapRole(db)
	model.BootstrapAccount(db)
	model.BootstrapAccountNotif(db)
	model.BootstrapPermission(db)
	log.Println("Bootstrap completed")
	return nil
}