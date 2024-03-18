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

func TestGetFilms(t *testing.T) {

	method := "GET"
	url := "/films"

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
				films := api_models.FilmsResponse{}
				err := json.Unmarshal(writer.Body.Bytes(), &films)
				if err != nil {
					assert.Fail(t, "Cant parse response to api_models.FilmsResponse")
				}
			}
		})
	}
}

func TestCreateFilms(t *testing.T) {

	method := "POST"
	url := "/films"
	// "2001-12-12" := time.Parse("2001-12-12", "2001-12-12")

	testCases := []struct {
		name        string
		username    string
		password    string
		film_name   string
		film_desc   string
		film_date   string
		film_rate   int
		film_actors []int
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
			film_name:   "Film3",
			film_desc:   "Film desc3",
			film_date:   "2000-10-10",
			film_rate:   4,
			film_actors: []int{1, 2},
			code:        200,
		},
		{
			name:        "Admin Auth Invalide Name",
			username:    "admin",
			password:    "admin",
			film_name:   "",
			film_desc:   "Film desc",
			film_date:   "2001-12-12",
			film_rate:   4,
			film_actors: []int{1, 2},
			code:        400,
		},
		{
			name:        "Admin Auth Invalide Rate",
			username:    "admin",
			password:    "admin",
			film_name:   "Film",
			film_desc:   "Film desc",
			film_date:   "2001-12-12",
			film_rate:   11,
			film_actors: []int{1, 2},
			code:        400,
		},
		{
			name:        "Admin Auth Invalide Missing Data",
			username:    "admin",
			password:    "admin",
			film_name:   "Film",
			film_desc:   "Film desc",
			film_rate:   1,
			film_actors: []int{1, 2},
			code:        400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			if tc.film_name != "" {
				body, _ = json.Marshal(api_models.CreateFilmRequest{

					Name:        tc.film_name,
					Description: tc.film_desc,
					Date:        tc.film_date,
					Rate:        tc.film_rate,
					Actors:      tc.film_actors,
				})
				slog.Debug(string(body))
			}

			slog.Debug(string(body))
			request, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
			if tc.username != "" && tc.password != "" {
				request.SetBasicAuth(tc.username, tc.password)
			}
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, tc.code, writer.Code)

			if writer.Code == 200 {
				resp := api_models.FilmResponse{}
				json.Unmarshal(writer.Body.Bytes(), &resp)
				exists := resp.Success
				film := resp.Film

				assert.Equal(t, true, exists)
				assert.Equal(t, film.Name, tc.film_name)
			}
		})
	}
}

func TestUpdateFilms(t *testing.T) {

	method := "PUT"
	url := "/films"

	body, _ := json.Marshal(map[string]string{
		"name":  "TempName",
		"sex":   "female",
		"birth": "2001-01-01",
	})

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	request.SetBasicAuth("admin", "admin")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	testCases := []struct {
		name        string
		username    string
		password    string
		film_id     string
		film_name   string
		film_desc   string
		film_date   string
		film_rate   int
		film_actors []int
		code        int
		error       string
	}{
		{
			name:    "No Auth",
			film_id: "1",
			code:    401,
		},
		{
			name:     "Client Auth",
			username: "client",
			password: "client",
			film_id:  "1",
			code:     401,
		},
		{
			name:        "Admin Auth Valid Data",
			username:    "admin",
			password:    "admin",
			film_id:     "1",
			film_name:   "Film",
			film_desc:   "Film desc",
			film_date:   "2001-12-12",
			film_rate:   4,
			film_actors: []int{1, 2},
			code:        200,
		},
		{
			name:        "Admin Auth Invalide Rate",
			username:    "admin",
			password:    "admin",
			film_id:     "1",
			film_name:   "Film",
			film_desc:   "Film desc",
			film_date:   "2001-12-12",
			film_rate:   11,
			film_actors: []int{1, 2},
			code:        400,
		},
		{
			name:        "Admin Auth Missing Data",
			username:    "admin",
			password:    "admin",
			film_id:     "1",
			film_name:   "Film",
			film_desc:   "Film desc",
			film_rate:   0,
			film_actors: []int{1, 2},
			code:        200,
		},
		{
			name:        "Admin Auth Missing Data No Id",
			username:    "admin",
			password:    "admin",
			film_name:   "Film",
			film_desc:   "Film desc",
			film_rate:   0,
			film_actors: []int{1, 2},
			code:        405,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			url := "/films"
			if tc.film_name != "" {

				body, _ = json.Marshal(api_models.UpdateFilmRequest{
					Name:        tc.film_name,
					Description: tc.film_desc,
					Date:        tc.film_date,
					Rate:        tc.film_rate,
					Actors:      tc.film_actors,
				})

			}
			if tc.film_id != "" {
				url = url + "/" + tc.film_id
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
				resp := api_models.FilmResponse{}
				err := json.Unmarshal(writer.Body.Bytes(), &resp)
				if err != nil {
					panic(err)
				}
				exists := resp.Success
				film := resp.Film
				assert.Equal(t, true, exists)
				assert.Equal(t, film.Name, tc.film_name)
			}
		})
	}
}

func TestDeleteFilms(t *testing.T) {

	method := "DELETE"

	testCases := []struct {
		name     string
		username string
		password string
		film_id  string
		wrong_id string
		code     int
		error    string
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
			name:    "No Auth",
			film_id: "1",
			code:    401,
		},
		{
			name:     "Client Auth",
			film_id:  "1",
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
			film_id:  "1",
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
			// body_, _ := json.Marshal(map[string]string{
			// 	"name":  "TempName",
			// 	"sex":   "male",
			// 	"birth": "2001-01-01",
			// })
			// // Добавление пользователя в базу
			// request, _ := http.NewRequest("POST", "/films", bytes.NewBuffer(body_))
			// request.SetBasicAuth("admin", "admin")
			// writer := httptest.NewRecorder()
			// router.ServeHTTP(writer, request)

			// Получение всех пользователей, чтобы посчитать сколько было до удаления
			request, _ := http.NewRequest("GET", "/films", bytes.NewBufferString(""))
			request.SetBasicAuth("admin", "admin")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			films := api_models.FilmsResponse{}
			err := json.Unmarshal(writer.Body.Bytes(), &films)
			if err != nil {
				panic(err)
			}
			cap_films_init := len(films.Films)
			if tc.film_id != "" {
				tc.film_id = strconv.FormatInt(int64(films.Films[0].ID), 10)
			}
			slog.Debug(strconv.Itoa(cap_films_init))

			var body []byte
			url := "/films"

			if tc.film_id != "" {
				url = url + "/" + tc.film_id
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
				request, _ = http.NewRequest("GET", "/films", bytes.NewBufferString(""))
				request.SetBasicAuth("admin", "admin")
				writer = httptest.NewRecorder()
				router.ServeHTTP(writer, request)

				films := api_models.FilmsResponse{}
				err := json.Unmarshal(writer.Body.Bytes(), &films)
				if err != nil {
					panic(err)
				}
				cap_films := len(films.Films)
				assert.Equal(t, cap_films, cap_films_init-1)
			}
		})
	}
}
