package tests

import (
	"bytes"
	"encoding/json"
	api_models "filmoteka/api/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetActors(t *testing.T) {

	method := "GET"
	url := "/actors"

	testCases := []struct {
		name     string
		username string
		password string
		code     int
		error    string
	}{
		{
			name: "No Auth",
			code: 401,
		},
		{
			name:     "Client Auth",
			username: "client",
			password: "client",
			code:     200,
		},
		{
			name:     "Admin Auth",
			username: "admin",
			password: "admin",
			code:     200,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(method, url, bytes.NewBufferString(""))
			if tc.username != "" && tc.password != "" {
				request.SetBasicAuth(tc.username, tc.password)
			}
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, tc.code, writer.Code)
		})
	}
}

func TestCreateActors(t *testing.T) {

	method := "POST"
	url := "/actors"

	testCases := []struct {
		name        string
		username    string
		password    string
		actor_name  string
		actor_sex   string
		actor_birth string
		code        int
		error       string
	}{
		{
			name: "No Auth",
			code: 401,
		},
		{
			name:     "Client Auth",
			username: "client",
			password: "client",
			code:     401,
		},
		{
			name:        "Admin Auth Valid Data",
			username:    "admin",
			password:    "admin",
			actor_name:  "Actor1",
			actor_sex:   "male",
			actor_birth: "2001-01-01",
			code:        200,
		},
		{
			name:        "Admin Auth Invalide Sex",
			username:    "admin",
			password:    "admin",
			actor_name:  "Actor1",
			actor_sex:   "other",
			actor_birth: "2001-01-01",
			code:        400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			if tc.actor_name != "" {
				body, _ = json.Marshal(map[string]string{
					"name":  tc.actor_name,
					"sex":   tc.actor_sex,
					"birth": tc.actor_birth,
				})
			}
			request, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
			if tc.username != "" && tc.password != "" {
				request.SetBasicAuth(tc.username, tc.password)
			}
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, tc.code, writer.Code)

			if writer.Code == 200 {
				resp := api_models.ActorResponse{}
				json.Unmarshal(writer.Body.Bytes(), &resp)
				exists := resp.Success
				actor := resp.Actor

				assert.Equal(t, true, exists)
				assert.Equal(t, actor.Name, tc.actor_name)
			}
		})
	}
}

func TestUpdateActors(t *testing.T) {

	method := "PUT"

	body, _ := json.Marshal(map[string]string{
		"name":  "TempName",
		"sex":   "female",
		"birth": "2001-01-01",
	})

	request, _ := http.NewRequest("POST", "/actors", bytes.NewBuffer(body))
	request.SetBasicAuth("admin", "admin")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	testCases := []struct {
		name        string
		username    string
		password    string
		actor_id    string
		actor_name  string
		actor_sex   string
		actor_birth string
		code        int
		error       string
	}{
		{
			name: "No Auth No ID",
			code: 405,
		},
		{
			name:     "Client Auth No ID",
			username: "client",
			password: "client",
			code:     405,
		},
		{
			name:     "No Auth",
			actor_id: "1",
			code:     401,
		},
		{
			name:     "Client Auth",
			actor_id: "1",
			username: "client",
			password: "client",
			code:     401,
		},
		{
			name:        "Admin Auth Valid Data No ID",
			username:    "admin",
			password:    "admin",
			actor_name:  "NewActor1",
			actor_sex:   "male",
			actor_birth: "2001-01-01",
			code:        405,
		},
		{
			name:       "Admin Auth Valid Data",
			username:   "admin",
			password:   "admin",
			actor_id:   "1",
			actor_name: "NewActor1",
			// actor_sex:  "male",
			// actor_birth: "2001-01-01",
			code: 200,
		},
		{
			name:        "Admin Auth Invalide Sex",
			username:    "admin",
			password:    "admin",
			actor_id:    "1",
			actor_name:  "Actor1",
			actor_sex:   "other",
			actor_birth: "2001-01-01",
			code:        400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			url := "/actors"
			if tc.actor_name != "" {
				if tc.actor_sex != "" {
					body, _ = json.Marshal(map[string]string{
						"name": tc.actor_name,
						"sex":  tc.actor_sex,
					})
				} else {
					body, _ = json.Marshal(map[string]string{
						"name": tc.actor_name,
					})
				}
			}
			if tc.actor_id != "" {
				url = url + "/" + tc.actor_id
			}
			slog.Debug(url)
			request, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
			if tc.username != "" && tc.password != "" {
				request.SetBasicAuth(tc.username, tc.password)
			}

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, tc.code, writer.Code)

			if writer.Code == 200 {
				resp := api_models.ActorResponse{}
				json.Unmarshal(writer.Body.Bytes(), &resp)
				exists := resp.Success
				actor := resp.Actor

				assert.Equal(t, true, exists)
				assert.Equal(t, actor.Name, tc.actor_name)
			}
		})
	}
}

// func TestDeleteActors(t *testing.T) {

// 	method := "DELETE"
// 	url := "/actors"

// 	testCases := []struct {
// 		name     string
// 		username string
// 		password string
// 		code     int
// 		error    string
// 	}{
// 		{
// 			name: "No Auth",
// 			code: 401,
// 		},
// 		{
// 			name:     "Client Auth",
// 			username: "client",
// 			password: "client",
// 			code:     401,
// 		},
// 		{
// 			name:     "Admin Auth",
// 			username: "admin",
// 			password: "admin",
// 			code:     200,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			request, _ := http.NewRequest(method, url, bytes.NewBufferString(""))
// 			if tc.username != "" && tc.password != "" {
// 				request.SetBasicAuth(tc.username, tc.password)
// 			}
// 			writer := httptest.NewRecorder()
// 			router.ServeHTTP(writer, request)
// 			assert.Equal(t, tc.code, writer.Code)
// 		})
// 	}
// }
