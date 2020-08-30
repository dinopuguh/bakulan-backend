package database

import (
	"fmt"
	"strconv"

	"github.com/dinopuguh/bakulan-backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn   *gorm.DB
	User     = config.Config("DB_USER")
	Password = config.Config("DB_PASSWORD")
	Host     = config.Config("DB_HOST")
	DB       = config.Config("DB_NAME")
	Port     = config.Config("DB_PORT")
)

func Connect() (err error) {
	port, err := strconv.Atoi(Port)
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%d sslmode=disable", User, Password, Host, DB, port)
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}
