package database

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn   *gorm.DB
	User     = os.Getenv("BAKULAN_DB_USER")
	Password = os.Getenv("BAKULAN_DB_PASSWORD")
	Host     = os.Getenv("BAKULAN_DB_HOST")
	DB       = os.Getenv("BAKULAN_DB_NAME")
	Port     = os.Getenv("BAKULAN_DB_PORT")
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
