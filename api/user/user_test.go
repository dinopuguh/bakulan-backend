package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dinopuguh/bakulan-backend/api/user"
	"github.com/dinopuguh/bakulan-backend/database"
	"github.com/dinopuguh/bakulan-backend/response"
	"github.com/dinopuguh/bakulan-backend/routes"
	"github.com/stretchr/testify/assert"
)

var (
	createdUserId uint
)

func TestGetAll(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	router := routes.New()
	t.Run("Get all users", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)

		res, _ := router.Test(req, -1)
		resBody, _ := ioutil.ReadAll(res.Body)

		assert.Equalf(t, http.StatusOK, res.StatusCode, string(resBody))
	})
}

func TestNew(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	r := routes.New()

	type args struct {
		data        map[string]string
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid register", args{
			data: map[string]string{
				"name":     "Dino",
				"email":    "dino@email.com",
				"password": "123",
				"phone":    "0812345",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"Email already exist.", args{
			data: map[string]string{
				"name":     "Dino",
				"email":    "dino@email.com",
				"password": "123",
				"phone":    "0812345",
			},
			statusCode:  http.StatusBadRequest,
			contentType: "application/json",
		}},
		{"Body parser invalid", args{
			data: map[string]string{
				"name":     "Dino",
				"email":    "dino@email.com",
				"password": "123",
				"phone":    "0812345",
			},
			statusCode: http.StatusServiceUnavailable,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.args.data)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/users-register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)

			res, _ := r.Test(req, -1)
			resBody, _ := ioutil.ReadAll(res.Body)
			defer res.Body.Close()

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))

			if tt.args.statusCode == http.StatusOK {
				rb := response.Auth{}

				json.Unmarshal(resBody, &rb)

				u := new(user.User)
				userJson, _ := json.Marshal(rb.Owner.(map[string]interface{}))
				json.Unmarshal(userJson, &u)

				createdUserId = u.ID
			}
		})
	}
}

func TestLogin(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	r := routes.New()

	type args struct {
		data        map[string]string
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid login", args{
			data: map[string]string{
				"email":    "dino@email.com",
				"password": "123",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"User not found", args{
			data: map[string]string{
				"email":    "dinopuguh@email.com",
				"password": "123",
			},
			statusCode:  http.StatusNotFound,
			contentType: "application/json",
		}},
		{"Password incorrect", args{
			data: map[string]string{
				"email":    "dino@email.com",
				"password": "1234",
			},
			statusCode:  http.StatusUnauthorized,
			contentType: "application/json",
		}},
		{"Body parser invalid", args{
			data: map[string]string{
				"email":    "dino@email.com",
				"password": "123",
			},
			statusCode: http.StatusServiceUnavailable,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.args.data)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/users-login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)

			res, _ := r.Test(req, -1)
			resBody, _ := ioutil.ReadAll(res.Body)

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))
		})
	}
}

func TestDelete(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	r := routes.New()

	type args struct {
		userId     uint
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid delete user", args{
			userId:     createdUserId,
			statusCode: http.StatusOK,
		}},
		{"User not found", args{
			userId:     createdUserId + 1,
			statusCode: http.StatusNotFound,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := fmt.Sprintf("/api/v1/users/%d", tt.args.userId)
			req, _ := http.NewRequest(http.MethodDelete, endpoint, nil)

			res, _ := r.Test(req, -1)
			resBody, _ := ioutil.ReadAll(res.Body)

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))
		})
	}
}
