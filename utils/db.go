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
	_dbHost 		= os.Getenv("DB_HOST")
	_dbPort 		= os.Getenv("DB_PORT")
	_dbUser 		= os.Getenv("DB_USER")
	_dbPass 		= os.Getenv("DB_PASS")
	_dbName 		= os.Getenv("DB_NAME")
	_dbSSLMode 		= os.Getenv("DB_SSLMODE")
	_dbTimeZone 	= os.Getenv("DB_TIMEZONE")
	dbMinConn, _ 	= strconv.Atoi(os.Getenv("DB_MIN_CONN"))
	dbMaxConn, _ 	= strconv.Atoi(os.Getenv("DB_MAX_CONN"))
	dbMaxLT, _ 		= strconv.Atoi(os.Getenv("DB_MAX_LIFETIME"))
)

func ConnectDB() *gorm.DB {
	dsn := "host=" + _dbHost + " user=" + _dbUser + " password=" + _dbPass + " dbname=" + _dbName + " port=" + _dbPort + " sslmode=" + _dbSSLMode + " TimeZone=" + _dbTimeZone
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
		log.Printf("Connected to DB Postgres")
	}

  	return db
}