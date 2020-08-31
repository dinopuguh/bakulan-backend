package store_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dinopuguh/bakulan-backend/api/store"
	"github.com/dinopuguh/bakulan-backend/database"
	"github.com/dinopuguh/bakulan-backend/response"
	"github.com/dinopuguh/bakulan-backend/routes"
	"github.com/stretchr/testify/assert"
)

var (
	createdStoreId uint
)

func TestGetAll(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	router := routes.New()
	t.Run("Get all users", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/stores", nil)

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
		data        map[string]interface{}
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid register", args{
			data: map[string]interface{}{
				"name":          "Lumintu",
				"email":         "lumintu@email.com",
				"password":      "123",
				"phone":         "0812345",
				"open":          "08:00:00",
				"close":         "13:00:00",
				"delivery_time": "09:00:00",
				"type_id":       1,
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"Email already exist.", args{
			data: map[string]interface{}{
				"name":          "Lumintu",
				"email":         "lumintu@email.com",
				"password":      "123",
				"phone":         "0812345",
				"open":          "08:00:00",
				"close":         "13:00:00",
				"delivery_time": "09:00:00",
				"type_id":       1,
			},
			statusCode:  http.StatusBadRequest,
			contentType: "application/json",
		}},
		{"Body parser invalid", args{
			data: map[string]interface{}{
				"name":          "Lumintu",
				"email":         "lumintu@email.com",
				"password":      "123",
				"phone":         "0812345",
				"open":          "08:00:00",
				"close":         "13:00:00",
				"delivery_time": "09:00:00",
				"type_id":       1,
			},
			statusCode: http.StatusServiceUnavailable,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.args.data)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/stores-register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)

			res, _ := r.Test(req, -1)
			resBody, _ := ioutil.ReadAll(res.Body)
			defer res.Body.Close()

			as := assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))

			if tt.args.statusCode == http.StatusOK && as {
				rb := response.Auth{}

				json.Unmarshal(resBody, &rb)

				s := new(store.Store)
				storeJson, _ := json.Marshal(rb.Owner.(map[string]interface{}))
				json.Unmarshal(storeJson, &s)

				createdStoreId = s.ID
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
				"email":    "lumintu@email.com",
				"password": "123",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"Store not found", args{
			data: map[string]string{
				"email":    "slamet@email.com",
				"password": "123",
			},
			statusCode:  http.StatusNotFound,
			contentType: "application/json",
		}},
		{"Password incorrect", args{
			data: map[string]string{
				"email":    "lumintu@email.com",
				"password": "1234",
			},
			statusCode:  http.StatusUnauthorized,
			contentType: "application/json",
		}},
		{"Body parser invalid", args{
			data: map[string]string{
				"email":    "lumintu@email.com",
				"password": "123",
			},
			statusCode: http.StatusServiceUnavailable,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.args.data)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/stores-login", bytes.NewBuffer(reqBody))
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
		storeId    uint
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid delete store", args{
			storeId:    createdStoreId,
			statusCode: http.StatusOK,
		}},
		{"Store not found", args{
			storeId:    createdStoreId + 1,
			statusCode: http.StatusNotFound,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := fmt.Sprintf("/api/v1/stores/%d", tt.args.storeId)
			req, _ := http.NewRequest(http.MethodDelete, endpoint, nil)

			res, _ := r.Test(req, -1)
			resBody, _ := ioutil.ReadAll(res.Body)

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))
		})
	}
}
