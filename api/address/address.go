package address

import (
	"github.com/dinopuguh/bakulan-backend/database"
	"github.com/gofiber/fiber"
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	Name        string  `json:"name"`
	Latitude    float32 `json:"latitude" gorm:"type:decimal(10,2)"`
	Longitude   float32 `json:"longitude" gorm:"type:decimal(10,2)"`
	PostalCode  int     `json:"postal_code"`
	SubDistrict string  `json:"sub_district"`
	City        string  `json:"city"`
	Province    string  `json:"province"`
	OwnerID     int     `json:"owner_id"`
	OwnerType   string  `json:"owner_type"`
}

func GetById(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn

	var address Address
	db.Find(&address, id)

	c.JSON(address)
}

func New(c *fiber.Ctx) {
	db := database.DBConn

	address := new(Address)
	// if err := c.BodyParser(&address); err != nil {
	// 	c.Status(503).Send(err)
	// 	return
	// }
	address.Name = "Munung"
	address.Latitude = 12.03
	address.Longitude = 11.00
	address.PostalCode = 61234
	address.SubDistrict = "Jatikalen"
	address.City = "Nganjuk"
	address.Province = "Jawa Timur"
	address.OwnerID = 1
	address.OwnerType = "stores"

	db.Create(&address)

	c.JSON(address)
}
