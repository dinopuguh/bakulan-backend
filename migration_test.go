package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dinopuguh/bakulan-backend/database"
	"github.com/dinopuguh/bakulan-backend/routes"
	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}

	migrateDatabase()

	r := routes.New()

	type args struct {
		endpoint   string
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{"Get all users", args{
			endpoint:   "/api/v1/users",
			statusCode: http.StatusOK,
		}},
		{"Get all stores", args{
			endpoint:   "/api/v1/stores",
			statusCode: http.StatusOK,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.args.endpoint, nil)

			res, _ := r.Test(req, -1)
			resBody, _ := ioutil.ReadAll(res.Body)

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))
		})
	}
}
