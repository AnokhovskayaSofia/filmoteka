package tests

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	api_models "filmoteka/api/models"

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
			if writer.Code == 200 {
				actors := api_models.ActorsResponse{}
				err := json.Unmarshal(writer.Body.Bytes(), &actors)
				if err != nil {
					assert.Fail(t, "Cant parse response to api_models.ActorsResponse")
				}
			}
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
			name:        "Admin Auth Valid Data",
			username:    "admin",
			password:    "admin",
			actor_id:    "1",
			actor_name:  "NewActor1",
			actor_sex:   "male",
			actor_birth: "2001-01-01",
			code:        200,
		},
		{
			name:       "Admin Auth Valid Data Not Full",
			username:   "admin",
			password:   "admin",
			actor_id:   "1",
			actor_name: "NewActor 1",
			code:       200,
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

func TestDeleteActors(t *testing.T) {

	method := "DELETE"

	testCases := []struct {
		name        string
		username    string
		password    string
		actor_id    string
		wrong_id    string
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
			name:     "Admin Auth Valid Data No ID",
			username: "admin",
			password: "admin",
			code:     405,
		},
		{
			name:     "Admin Auth Valid Data",
			username: "admin",
			password: "admin",
			actor_id: "1",
			code:     200,
		},
		{
			name:     "Admin Auth Invalide Sex",
			username: "admin",
			password: "admin",
			wrong_id: "45",
			code:     400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body_, _ := json.Marshal(map[string]string{
				"name":  "TempName",
				"sex":   "male",
				"birth": "2001-01-01",
			})
			// Добавление пользователя в базу
			request, _ := http.NewRequest("POST", "/actors", bytes.NewBuffer(body_))
			request.SetBasicAuth("admin", "admin")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			// Получение всех пользователей, чтобы посчитать сколько было до удаления
			request, _ = http.NewRequest("GET", "/actors", bytes.NewBufferString(""))
			request.SetBasicAuth("admin", "admin")
			writer = httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			actors := api_models.ActorsResponse{}
			err := json.Unmarshal(writer.Body.Bytes(), &actors)
			if err != nil {
				panic(err)
			}
			cap_actors_init := len(actors.Actors)
			if tc.actor_id != "" {
				tc.actor_id = strconv.FormatInt(actors.Actors[0].ID, 10)
			}
			slog.Debug(strconv.Itoa(cap_actors_init))

			var body []byte
			url := "/actors"

			if tc.actor_id != "" {
				url = url + "/" + tc.actor_id
			}
			if tc.wrong_id != "" {
				url = url + "/" + tc.wrong_id
			}
			request, _ = http.NewRequest(method, url, bytes.NewBuffer(body))
			if tc.username != "" && tc.password != "" {
				request.SetBasicAuth(tc.username, tc.password)
			}

			writer = httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, tc.code, writer.Code)

			if writer.Code == 200 {
				request, _ = http.NewRequest("GET", "/actors", bytes.NewBufferString(""))
				request.SetBasicAuth("admin", "admin")
				writer = httptest.NewRecorder()
				router.ServeHTTP(writer, request)

				actors := api_models.ActorsResponse{}
				err := json.Unmarshal(writer.Body.Bytes(), &actors)
				if err != nil {
					panic(err)
				}
				cap_actor := len(actors.Actors)
				assert.Equal(t, cap_actor, cap_actors_init-1)
			}
		})
	}
}
