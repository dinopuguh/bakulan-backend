package routes

import (
	"github.com/dinopuguh/bakulan-backend/api/address"
	"github.com/dinopuguh/bakulan-backend/api/auth"
	"github.com/dinopuguh/bakulan-backend/api/store"
	"github.com/dinopuguh/bakulan-backend/api/user"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt"
)

func New() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())

	app.Post("/api/v1/stores-register", store.New)
	app.Post("/api/v1/stores-login", store.Login)

	app.Post("/api/v1/users-register", user.New)
	app.Post("/api/v1/users-login", user.Login)

	app.Get("/api/v1/stores", store.GetAll)
	app.Delete("/api/v1/stores/:id", store.Delete)

	app.Get("/api/v1/users", user.GetAll)
	app.Delete("/api/v1/users/:id", user.Delete)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: auth.MySigningKey,
		Claims:     &auth.JwtCustomClaims{},
	}))

	app.Post("/api/v1/address", address.New)

	return app
}
