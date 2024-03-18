package api

import (
	"encoding/json"
	"errors"
	api_models "filmoteka/api/models"
	"filmoteka/db"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-pg/pg/v10"
)

// getFilms godoc
// @Summary      Get films list
// @Description  Availible only for authenticated user, getting films list, they can be sorted by fields, default is rate. Also you can use filters in field.value template.
// @Tags         films
// @Accept       json
// @Produce      json
// @Router       /films [get]
// @Param sortBy query string false "Sort by field, default rate" example(name)
// @Param filter query string false "Filter by field (field.value), can be user all except actors" example(name.Name1)
// @Security BasicAuth
// @Success 200 {object} api_models.FilmsResponse
// @Failure 401 {object}  ErrorResponse
// @Failure 400 {object}  ErrorResponse
func getFilms(w http.ResponseWriter, r *http.Request) {
	_, err := checkBasicAuth(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		HandleError(w, err)
		return
	}

	sortBy := r.URL.Query().Get("sortBy")
	filter := r.URL.Query().Get("filter")
	var splits []string
	if filter != "" {
		splits = strings.Split(filter, ".")
	}
	if sortBy == "" || sortBy == "rate" {
		sortBy = "rate DESC"
	}

	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	films, err := db.GetFilms(pgdb, sortBy, splits)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	res := &api_models.FilmsResponse{
		Success: true,
		Error:   "",
		Films:   films,
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// createFilm godoc
// @Summary      Create film
// @Description  Availible only for admin user, creating film using data from request body and return new film
// @Tags         films
// @Accept       json
// @Produce      json
// @Router       /films [post]
// @Param Film body db.Film true "film info"
// @Param filmID query string true "Film Id"
// @Security BasicAuth
// @Success 200 {object} api_models.FilmResponse
// @Failure 401 {object}  ErrorResponse
// @Failure 400 {object}  ErrorResponse
func createFilm(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		w.WriteHeader(http.StatusUnauthorized)
		HandleError(w, err)
		return
	}

	req := &api_models.CreateFilmRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}
	err = Validate.Struct(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	pgdb, ok := r.Context().Value("DB").(*pg.DB)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	datetime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		slog.Debug(req.Date)
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}
	film, err := db.CreateFilm(pgdb, &db.Film{
		Name:        req.Name,
		Description: req.Description,
		Date:        datetime,
		Rate:        req.Rate,
	}, req.Actors)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	res := &api_models.FilmResponse{
		Success: true,
		Error:   "",
		Film:    film,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// updateFilm godoc
// @Summary      Update film
// @Description  Availible only for admin user, updating film using data from request body and return new film
// @Tags         films
// @Accept       json
// @Produce      json
// @Router       /films [put]
// @Param Film body db.Film true "film info"
// @Param filmID query string true "Film Id"
// @Security BasicAuth
// @Success 200 {object} api_models.FilmResponse
// @Failure 401 {object}  ErrorResponse
// @Failure 400 {object}  ErrorResponse
func updateFilm(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		w.WriteHeader(http.StatusUnauthorized)
		HandleError(w, err)
		return
	}

	req := &api_models.UpdateFilmRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}
	err = Validate.Struct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	filmID := chi.URLParam(r, "filmID")

	intFilmID, err := strconv.Atoi(filmID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}
	var datetime time.Time
	if req.Date != "" {
		datetime, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			slog.Debug(req.Date)
			w.WriteHeader(http.StatusBadRequest)
			HandleError(w, err)
			return
		}
	}

	film, err := db.UpdateFilm(pgdb, &db.Film{
		ID:          intFilmID,
		Name:        req.Name,
		Description: req.Description,
		Date:        datetime,
		Rate:        req.Rate,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	res := &api_models.FilmResponse{
		Success: true,
		Error:   "",
		Film:    film,
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// deleteFilm godoc
// @Summary      Delete film
// @Description  Availible only for admin user, deleting film by id from params
// @Tags         films
// @Accept       json
// @Produce      json
// @Param filmID query string true "Film Id"
// @Router       /films [delete]
// @Security BasicAuth
// @Success 200 {object} nil
// @Failure 401 {object}  ErrorResponse
// @Failure 400 {object}  ErrorResponse
func deleteFilm(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		w.WriteHeader(http.StatusUnauthorized)
		HandleError(w, err)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	filmID := chi.URLParam(r, "filmID")
	intFilmID, err := strconv.ParseInt(filmID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	err = db.DeleteFilm(pgdb, intFilmID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
