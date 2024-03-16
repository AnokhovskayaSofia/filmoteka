package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	api_models "filmoteka/api/models"
	db_models "filmoteka/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"
)

func StartAPI(pgdb *pg.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.RequestID, middleware.Recoverer, middleware.WithValue("DB", pgdb))

	r.Route("/films", func(r chi.Router) {
		r.Get("/", getFilms)
		r.Post("/", createFilm)
		r.Put("/{filmID}", updateFilm)
		r.Delete("/{filmID}", deleteFilm)
	})
	r.Route("/actors", func(r chi.Router) {
		r.Get("/", getActors)
		r.Post("/", createActor)
		r.Put("/{actorID}", updateActor)
		r.Delete("/{actorID}", deleteActor)
	})

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value("DB").(*pg.DB)
		if !ok {
			msg := "could not get the DB from context"
			slog.Error(msg)
			w.Write([]byte(msg))
			return
		}
		w.Write([]byte("OK"))
		w.WriteHeader(http.StatusOK)
	})

	slog.Info("Success start API routes")
	return r
}

func checkBasicAuth(r *http.Request) (string, error) {
	user, pass, ok := r.BasicAuth()
	var err error
	if !ok {
		err := errors.New("failed to get username and password")
		return "", err
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	user_role, err := db_models.GetUser(pgdb, user, pass)

	return user_role, err
}

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

	films, err := db_models.GetFilms(pgdb, sortBy, splits)
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
	if err != nil || auth_role != db_models.Admin {
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
	film, err := db_models.CreateFilm(pgdb, &db_models.Film{
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
	if err != nil || auth_role != db_models.Admin {
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

	film, err := db_models.UpdateFilm(pgdb, &db_models.Film{
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
	if err != nil || auth_role != db_models.Admin {
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

	err = db_models.DeleteFilm(pgdb, intFilmID)
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

func getActors(w http.ResponseWriter, r *http.Request) {
	_, err := checkBasicAuth(r)
	if err != nil {
		res := &api_models.ActorsResponse{
			Success: false,
			Error:   err.Error(),
			Actors:  nil,
		}
		slog.Error("error getting actors %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
		return
	}

	pgdb, ok := r.Context().Value("DB").(*pg.DB)

	if !ok {
		res := &api_models.ActorsResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Actors:  nil,
		}
		err := json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	actors, err := db_models.GetActors(pgdb)
	if err != nil {
		res := &api_models.ActorsResponse{
			Success: false,
			Error:   err.Error(),
			Actors:  nil,
		}
		slog.Error("error getting actors %v\n", err)
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			slog.Error("error sending response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := &api_models.ActorsResponse{
		Success: true,
		Error:   "",
		Actors:  actors,
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("error encoding actors: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func createActor(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db_models.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		slog.Error("error getting actor %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
		return
	}
	req := &api_models.CreateActorRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	user, password, _ := r.BasicAuth()
	slog.Info("Request", user, password)

	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
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
		res := &api_models.ActorResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Actor:   nil,
		}
		err := json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	birthday, err := time.Parse("2006-01-02", req.Birth)
	if err != nil {
		slog.Error("failed to parse date %v\n", err)
	}
	actor, err := db_models.CreateActor(pgdb, &db_models.Actor{
		Name:  req.Name,
		Sex:   req.Sex,
		Birth: birthday,
	})
	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := &api_models.ActorResponse{
		Success: true,
		Error:   "",
		Actor:   actor,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("error encoding after creating actor %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func updateActor(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db_models.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		slog.Error("error getting actor %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
		return
	}
	req := &api_models.UpdateActorRequest{}

	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
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
		res := &api_models.ActorResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Actor:   nil,
		}
		err := json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	actorID := chi.URLParam(r, "actorID")

	intActorID, err := strconv.ParseInt(actorID, 10, 64)
	datetime, err := time.Parse("2006-01-02", req.Birth)

	actor, err := db_models.UpdateActor(pgdb, &db_models.Actor{
		ID:    intActorID,
		Name:  req.Name,
		Sex:   req.Sex,
		Birth: datetime,
	})
	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := &api_models.ActorResponse{
		Success: true,
		Error:   "",
		Actor:   actor,
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("error encoding actor: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteActor(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db_models.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		slog.Error("error getting actor %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(res)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   "could not get the DB from context",
			Actor:   nil,
		}
		err := json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	actorID := chi.URLParam(r, "actorID")
	intActorID, err := strconv.ParseInt(actorID, 10, 64)
	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
		}
		err = json.NewEncoder(w).Encode(res)

		if err != nil {
			slog.Error("error sending response %v\n", err)
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db_models.DeleteActor(pgdb, intActorID)
	if err != nil {
		res := &api_models.ActorResponse{
			Success: false,
			Error:   err.Error(),
			Actor:   nil,
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
