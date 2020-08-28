package routes

import (
	"github.com/dinopuguh/bakulan-backend/api/address"
	"github.com/dinopuguh/bakulan-backend/api/auth"
	"github.com/dinopuguh/bakulan-backend/api/store"
	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt"
)

func New(app *fiber.App) {
	app.Get("/api/v1/stores", store.GetAll)
	app.Post("/api/v1/stores-register", store.New)
	app.Post("/api/v1/stores-login", store.Login)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: auth.MySigningKey,
		Claims:     &auth.JwtCustomClaims{},
	}))
	app.Post("/api/v1/address", address.New)
}
