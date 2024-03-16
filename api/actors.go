package api

import (
	"encoding/json"
	"errors"
	api_models "filmoteka/api/models"
	"filmoteka/db"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-pg/pg/v10"
)

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

	actors, err := db.GetActors(pgdb)
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
	if err != nil || auth_role != db.Admin {
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
	actor, err := db.CreateActor(pgdb, &db.Actor{
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
	if err != nil || auth_role != db.Admin {
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

	actor, err := db.UpdateActor(pgdb, &db.Actor{
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
	if err != nil || auth_role != db.Admin {
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

	err = db.DeleteActor(pgdb, intActorID)
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
