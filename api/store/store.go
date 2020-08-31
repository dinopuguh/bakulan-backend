package store

import (
	"net/http"

	"github.com/dinopuguh/bakulan-backend/api/address"
	"github.com/dinopuguh/bakulan-backend/api/auth"
	"github.com/dinopuguh/bakulan-backend/database"
	"github.com/dinopuguh/bakulan-backend/helpers"
	"github.com/dinopuguh/bakulan-backend/response"
	"github.com/gofiber/fiber"

	"gorm.io/gorm"
)

type Store struct {
	gorm.Model
	Name         string            `json:"name"`
	Email        string            `json:"email"`
	Password     string            `json:"password"`
	Phone        string            `json:"phone"`
	Open         string            `json:"open"`
	Close        string            `json:"close"`
	DeliveryTime string            `json:"delivery_time"`
	Address      []address.Address `json:"address" gorm:"polymorphic:Owner;polymorphicValue:stores"`
	TypeID       int               `json:"type_id"`
}

func GetAll(c *fiber.Ctx) {
	db := database.DBConn

	var stores []Store
	if res := db.Preload("Address").Find(&stores); res.Error != nil {
		c.Status(http.StatusServiceUnavailable).JSON(response.Error{Message: res.Error.Error()})
		return
	}

	c.JSON(stores)
}

func New(c *fiber.Ctx) {
	db := database.DBConn

	store := new(Store)
	if err := c.BodyParser(&store); err != nil {
		c.Status(http.StatusServiceUnavailable).JSON(response.Error{Message: err.Error()})
		return
	}

	var cst Store
	var err error
	res := db.Where("email = ?", store.Email).First(&cst)

	if res.RowsAffected > 0 {
		c.Status(http.StatusBadRequest).JSON(response.Error{Message: "Store with this email already exist."})
		return
	}

	store.Password, err = helpers.HashPassword(store.Password)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(response.Error{Message: err.Error()})
		return
	}

	if res = db.Create(store); res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(response.Error{Message: err.Error()})
		return
	}

	res = db.Where("email = ?", store.Email).First(&cst)

	token, err := auth.GenerateJWT(cst.Name, cst.Email)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(response.Error{Message: err.Error()})
		return
	}

	c.JSON(response.Auth{
		Owner:       cst,
		AccessToken: token,
	})
	return
}

func Login(c *fiber.Ctx) {
	db := database.DBConn

	login := new(auth.Login)
	if err := c.BodyParser(&login); err != nil {
		c.Status(http.StatusServiceUnavailable).JSON(response.Error{Message: err.Error()})
		return
	}

	var store Store
	res := db.Where("email = ?", login.Email).First(&store)

	if res.RowsAffected == 0 {
		c.Status(http.StatusNotFound).JSON(response.Error{Message: "Store not found."})
		return
	}

	if !helpers.CheckPasswordHash(login.Password, store.Password) {
		c.Status(http.StatusUnauthorized).JSON(response.Error{Message: "Password incorrect."})
		return
	}

	token, err := auth.GenerateJWT(store.Name, store.Email)
	if err != nil {
		c.Status(http.StatusUnauthorized).JSON(response.Error{Message: err.Error()})
		return
	}

	c.JSON(response.Auth{
		Owner:       store,
		AccessToken: token,
	})
}

func Delete(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn

	var store Store
	res := db.First(&store, id)

	if res.RowsAffected == 0 {
		c.Status(http.StatusNotFound).JSON(response.Error{Message: "Store not found."})
		return
	}

	if res = db.Delete(&store); res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(response.Error{Message: res.Error.Error()})
		return
	}

	c.JSON(response.Error{Message: "Store deleted."})
}
