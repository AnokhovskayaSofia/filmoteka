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

func getFilms(w http.ResponseWriter, r *http.Request) {
	_, err := checkBasicAuth(r)
	if err != nil {
		res := &api_models.FilmsResponse{
			Success: false,
			Error:   err.Error(),
			Films:   nil,
		}
		slog.Error("error getting films %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
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
		res := &api_models.FilmsResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Films:   nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			slog.Error("error sending response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	films, err := db.GetFilms(pgdb, sortBy, splits)
	if err != nil {
		res := &api_models.FilmsResponse{
			Success: false,
			Error:   err.Error(),
			Films:   nil,
		}
		slog.Error("error getting films %v\n", err)
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			slog.Error("error sending response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := &api_models.FilmsResponse{
		Success: true,
		Error:   "",
		Films:   films,
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("error encoding films: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func createFilm(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		slog.Error("error creating film %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
		return
	}

	req := &api_models.CreateFilmRequest{}
	err = json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pgdb, ok := r.Context().Value("DB").(*pg.DB)

	if !ok {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Film:    nil,
		}
		err := json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	datetime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		slog.Error("failed to parse date %v\n", err)
	}
	film, err := db.CreateFilm(pgdb, &db.Film{
		Name:        req.Name,
		Description: req.Description,
		Date:        datetime,
		Rate:        req.Rate,
	}, req.Actors)
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := &api_models.FilmResponse{
		Success: true,
		Error:   "",
		Film:    film,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("error encoding after creating film %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func updateFilm(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		slog.Error("error updating film %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
		return
	}

	req := &api_models.UpdateFilmRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Film:    nil,
		}
		err := json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filmID := chi.URLParam(r, "filmID")

	intFilmID, err := strconv.Atoi(filmID)
	datetime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
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
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(film)
	if err != nil {
		slog.Error("error encoding film: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteFilm(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		slog.Error("error deleting film %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Film:    nil,
		}
		err := json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filmID := chi.URLParam(r, "filmID")
	intFilmID, err := strconv.ParseInt(filmID, 10, 64)
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.DeleteFilm(pgdb, intFilmID)
	if err != nil {
		res := &api_models.FilmResponse{
			Success: false,
			Error:   err.Error(),
			Film:    nil,
		}
		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
