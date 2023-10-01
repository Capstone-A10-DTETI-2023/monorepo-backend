package utils

import (
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
)

var (
	dbHost 		= os.Getenv("DB_HOST")
	dbPort 		= os.Getenv("DB_PORT")
	dbUser 		= os.Getenv("DB_USER")
	dbPass 		= os.Getenv("DB_PASS")
	dbName 		= os.Getenv("DB_NAME")
	dbSSLMode 	= os.Getenv("DB_SSLMODE")
	dbTimeZone 	= os.Getenv("DB_TIMEZONE")
	dbMinConn, _ 	= strconv.Atoi(os.Getenv("DB_MIN_CONN"))
	dbMaxConn, _ 	= strconv.Atoi(os.Getenv("DB_MAX_CONN"))
	dbMaxLT, _ 		= strconv.Atoi(os.Getenv("DB_MAX_LIFETIME"))
)

func ConnectDB() *gorm.DB {
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSSLMode + " TimeZone=" + dbTimeZone
//   dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
  	}

	dbPool, _ := db.DB()
	dbPool.SetMaxIdleConns(dbMinConn)
	dbPool.SetMaxOpenConns(dbMaxConn)
	dbPool.SetConnMaxLifetime(time.Duration(dbMaxLT)*time.Hour)
	
	if err := dbPool.Ping(); err == nil {
		log.Println("Database connection established")
	}

  	return db
}