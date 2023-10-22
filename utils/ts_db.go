package utils

import (
	"context"
	"os"

	pgx "github.com/jackc/pgx/v5"
)

var (
	_dbTsHost 		= os.Getenv("DBTS_HOST")
	_dbTsPort 		= os.Getenv("DBTS_PORT")
	_dbTsUser 		= os.Getenv("DBTS_USER")
	_dbTsPass 		= os.Getenv("DBTS_PASS")
	_dbTsName 		= os.Getenv("DBTS_NAME")
	_dbTsSSLMode 	= os.Getenv("DBTS_SSLMODE")
	_dbTsTimeZone 	= os.Getenv("DBTS_TIMEZONE")
)

func ConnectTSDB() *pgx.Conn {
	dsn := "postgresql://" + _dbTsUser + ":" + _dbTsPass + "@" + _dbTsHost + ":" + _dbTsPort + "/" + _dbTsName
	connConfig, _ := pgx.ParseConfig(dsn)
	db, err := pgx.ConnectConfig(context.Background(), connConfig)

	if err != nil {
		panic(err)
	}

	return db
}
