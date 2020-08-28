package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)

func Connect() (err error) {
	dsn := "user=postgres password=postgres host=localhost dbname=bakulan port=5432 sslmode=disable"
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}
