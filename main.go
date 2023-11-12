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
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/service"
	"github.com/Capstone-A10-DTETI-2023/monorepo-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/robfig/cron/v3"

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
	defer dbTs.Close(context.Background())
	if err := dbTs.Ping(context.Background()); err != nil {
		panic(err)
	}
	log.Printf("Connected to TSDB QuestDB")

	// Initialize MQTT Connection
	go func(){
		db := utils.ConnectDB()
		websocket := utils.ConnectWS()
		mqtt := consumer.ConnectMQTT(db, websocket)
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
	user.Delete("/:id", userController.DeleteUser)

	auth := app.Group("/auth")
	authController := &controller.AuthController{DB: db}
	auth.Post("/login", authController.Login)
	auth.Post("/reset-password", authController.ResetPassword)
	auth.Post("/logout", authController.Logout)
	auth.Post("/register", authController.Register)

	role := app.Group("/roles")
	role.Use(middleware.IsAuthenticated)
	roleController := &controller.RoleController{DB: db}
	role.Get("/", roleController.GetAllRole)
	role.Post("/", roleController.CreateRole)
	role.Put("/:id", roleController.UpdateRole)
	role.Delete("/:id", roleController.DeleteRoleByID)

	permissions := app.Group("/permissions")
	permissions.Use(middleware.IsAuthenticated)
	permissionController := &controller.PermissionController{DB: db}
	permissions.Get("/", permissionController.GetAllPermission)
	permissions.Get("/:id", permissionController.GetPermissionsByID)
	permissions.Put("/:id", permissionController.UpdatePermission)

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
	node.Delete("/:id", nodeController.DeleteNodeByID)
	node.Put("/:id", nodeController.UpdateNodeByID)

	sensor := app.Group("/sensors")
	sensor.Use(middleware.IsAuthenticated)
	sensorController := &controller.SensorController{DB: db}
	sensor.Post("/", sensorController.AddNewSensor)
	sensor.Get("/", sensorController.GetAllSensors)
	sensor.Delete("/:id", sensorController.DeleteSensor)
	sensor.Put("/:id", sensorController.UpdateSensorByID)
	sensor.Get("/node/:node_id", sensorController.GetSensorByNodeID)

	sensorData := app.Group("/tsdata")
	sensorData.Use(middleware.IsAuthenticated)
	sensorDataController := &controller.SensorDataController{DB: db}
	sensorData.Post("/sensor", sensorDataController.InsertDataSensor)
	sensorData.Get("/sensor", sensorDataController.GetSensorData)
	sensorData.Get("/sensor/last", sensorDataController.GetLastSensorData)

	systemSettings := app.Group("/sys-setting")
	systemSettings.Use(middleware.IsAuthenticated)
	systemSettingsController := &controller.SysSettingController{DB: db}
	systemSettings.Get("/nodepref", systemSettingsController.GetAllNodePressureRef)
	systemSettings.Put("/nodepref/:node_id", systemSettingsController.SetNodePressureRef)

	leakage := app.Group("/leakage")
	leakage.Use(middleware.IsAuthenticated)
	leakageController := &controller.LeakageController{DB: db}
	leakage.Get("/sensmat", leakageController.GetSensMat)
	leakage.Get("/resmat", leakageController.GetResidualMatrix)
	leakage.Get("/sensor/last", leakageController.GetLatestPresSensorData)
	leakage.Get("/status", leakageController.GetLeakageStatus)

	log.Println("Starting scheduler")
   	scheduler := cron.New()
	scheduler.AddFunc("*/1 * * * *", scheduleLeakDetection)
	go scheduler.Start()

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
	model.MigrateNodePressureRef(db)
	model.MigrateSystemSetting(db)
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
	model.BootstrapSystemSetting(db)
	model.DropNodeData()
	model.MigrateNodeData()
	log.Println("Bootstrap completed")
	return nil
}

func scheduleLeakDetection() {
	dbPG := utils.ConnectDB()
	dbTs := model.ConnectDBTS()

	websocket := utils.ConnectWS()

	sensMat, _ := service.CalculateSensMatrix(dbPG)
	resMat, _ := service.CalculateResidualMatrix(dbPG, dbTs)
	nodeLeaking, _ := service.GetLeakageNode(sensMat, resMat, dbPG)
	if nodeLeaking != -1 {
		log.Println("Node " + fmt.Sprint(nodeLeaking) + " terdeteksi bocor")
		var phoneNum []string
		rows, _ := dbPG.Table("users").Select("users.phone_num").Joins("left join notifications on notifications.user_id = users.id").Where("notifications.whatsapp = ?", true).Rows()
		defer rows.Close()
		for rows.Next() {
			var phone string
			rows.Scan(&phone)
			phoneNum = append(phoneNum, phone)
		}

		message := "Halo! Node " + fmt.Sprint(nodeLeaking) + " terdeteksi bocor. Silahkan cek node tersebut."
		for _, phone := range phoneNum {
			utils.SendWAMessage(phone, message, "0")
		}

		data := map[string]string{"node_leak": fmt.Sprint(nodeLeaking)}
		websocket.Trigger("my-channel", "leakage", data)
	}
}