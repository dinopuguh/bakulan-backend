package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var MySigningKey = []byte(os.Getenv("BAKULAN_JWT_TOKEN"))

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(name, email string) (string, error) {
	claims := &JwtCustomClaims{
		name,
		email,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(MySigningKey)
	if err != nil {
		return "", fmt.Errorf("Failed to generate JWT")
	}

	return t, nil
}

func GetJWTClaims(c echo.Context) *JwtCustomClaims {
	u := c.Get("user").(*jwt.Token)
	return u.Claims.(*JwtCustomClaims)
}
