package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"stock-tracker/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	username string
	password string
	host     string
	port     string
	name     string
	Client   *gorm.DB
}

type CloseConn func()

func closeConn(conn *sql.DB) CloseConn {
	return func() {
		_ = conn.Close()
	}
}

var Dns = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
	config.Database.Host,
	config.Database.Username,
	config.Database.Password,
	config.Database.Name,
	config.Database.Port,
)

func NewDatabase() (*gorm.DB, error) {
	logLevel := logger.Silent
	if config.Debug {
		logLevel = logger.Info
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	client, err := gorm.Open(postgres.Open(Dns), &gorm.Config{
		Logger: newLogger,
	})

	return client, err
}
