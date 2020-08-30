package user

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

type User struct {
	gorm.Model
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Phone    string            `json:"phone"`
	Address  []address.Address `json:"address" gorm:"polymorphic:Owner;polymorphicValue:users"`
}

func GetAll(c *fiber.Ctx) {
	db := database.DBConn

	var users []User
	if res := db.Preload("Address").Find(&users); res.Error != nil {
		c.Status(http.StatusServiceUnavailable).JSON(response.Error{Message: res.Error.Error()})
		return
	}

	c.JSON(users)
}

func New(c *fiber.Ctx) {
	db := database.DBConn

	user := new(User)
	if err := c.BodyParser(&user); err != nil {
		c.Status(http.StatusServiceUnavailable).JSON(response.Error{Message: err.Error()})
		return
	}

	var cu User
	var err error
	res := db.Where("email = ?", user.Email).First(&cu)

	if res.RowsAffected > 0 {
		c.Status(http.StatusBadRequest).JSON(response.Error{Message: "User with this email already exist."})
		return
	}

	user.Password, err = helpers.HashPassword(user.Password)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(response.Error{Message: err.Error()})
		return
	}

	if res = db.Create(user); res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(response.Error{Message: res.Error.Error()})
		return
	}

	res = db.Where("email = ?", user.Email).First(&cu)

	token, err := auth.GenerateJWT(cu.Name, cu.Email)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(response.Error{Message: err.Error()})
		return
	}

	c.JSON(response.Auth{
		Owner:       cu,
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

	var user User
	res := db.Where("email = ?", login.Email).First(&user)

	if res.RowsAffected == 0 {
		c.Status(http.StatusNotFound).JSON(response.Error{Message: "User not found."})
		return
	}

	if !helpers.CheckPasswordHash(login.Password, user.Password) {
		c.Status(http.StatusUnauthorized).JSON(response.Error{Message: "Password incorrect."})
		return
	}

	token, err := auth.GenerateJWT(user.Name, user.Email)
	if err != nil {
		c.Status(http.StatusUnauthorized).JSON(response.Error{Message: err.Error()})
		return
	}

	c.JSON(response.Auth{
		Owner:       user,
		AccessToken: token,
	})
}

func Delete(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn

	var user User
	res := db.First(&user, id)

	if res.RowsAffected == 0 {
		c.Status(http.StatusNotFound).JSON(response.Error{Message: "User not found."})
		return
	}

	if res = db.Delete(&user); res.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(response.Error{Message: res.Error.Error()})
		return
	}

	c.JSON(response.Error{Message: "User deleted."})
}
