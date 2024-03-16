package api

import (
	"errors"
	"log/slog"
	"net/http"

	"filmoteka/db"

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
	user_role, err := db.GetUser(pgdb, user, pass)

	return user_role, err
}
