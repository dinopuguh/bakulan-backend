package main

import (
	"log"

	"github.com/dinopuguh/bakulan-backend/api/address"
	"github.com/dinopuguh/bakulan-backend/api/store"
	"github.com/dinopuguh/bakulan-backend/api/user"
	"github.com/dinopuguh/bakulan-backend/database"
	"github.com/dinopuguh/bakulan-backend/routes"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
)

func migrateDatabase() {
	database.DBConn.AutoMigrate(&store.Store{})
	database.DBConn.AutoMigrate(&user.User{})
	database.DBConn.AutoMigrate(&address.Address{})

	log.Println("Database migrated")
}

func main() {
	app := fiber.New()
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}

	migrateDatabase()

	app.Use(cors.New())
	routes.New(app)

	app.Listen(3000)
}
