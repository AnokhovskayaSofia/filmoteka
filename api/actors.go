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
	"github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/v5"
	"github.com/go-pg/pg/v10"
)

var validate *validator.Validate = validator.New()

// getActors godoc
// @Summary      List actors
// @Description  get actors list from db
// @Tags         actors
// @Accept       json
// @Produce      json
// @Router       /actors [get]
// @Security BasicAuth
// @Success 200 {object} api_models.ActorsResponse
// @Failure 401 {object}  api_models.ActorsResponse
func getActors(w http.ResponseWriter, r *http.Request) {
	_, err := checkBasicAuth(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		HandleActorError(w, err)
		return
	}

	pgdb, ok := r.Context().Value("DB").(*pg.DB)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorsError(w, err)
		return
	}

	actors, err := db.GetActors(pgdb)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorsError(w, err)
		return
	}

	res := &api_models.ActorsResponse{
		Success: true,
		Error:   "",
		Actors:  actors,
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorsError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// createActor godoc
// @Summary      Create actor
// @Description  create actor using data from request body and return new actor
// @Tags         actors
// @Accept       json
// @Produce      json
// @Router       /actors [post]
// @Param Actor body db.Actor true "actor info"
// @Security BasicAuth
// @Success 200 {object} api_models.ActorResponse
// @Failure 401 {object}  api_models.ActorResponse
func createActor(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		w.WriteHeader(http.StatusUnauthorized)
		HandleActorError(w, err)
		return
	}
	req := &api_models.CreateActorRequest{}
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(req)
	err = validate.Struct(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}

	pgdb, ok := r.Context().Value("DB").(*pg.DB)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, errors.New("could not get the DB from context"))
		return
	}

	birthday, err := time.Parse("2006-01-02", req.Birth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}
	actor, err := db.CreateActor(pgdb, &db.Actor{
		Name:  req.Name,
		Sex:   req.Sex,
		Birth: birthday,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}

	res := &api_models.ActorResponse{
		Success: true,
		Error:   "",
		Actor:   actor,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("error encoding after creating actor %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// updateActor godoc
// @Summary      Update actor
// @Description  update actor using id from request params and return actor
// @Tags         actors
// @Accept       json
// @Produce      json
// @Param actorID query string true "Actors Id"
// @Param Actor body db.Actor true "actor info"
// @Router       /actors [put]
// @Security BasicAuth
// @Success 200 {object} api_models.ActorResponse
// @Failure 401 {object}  api_models.ActorResponse
func updateActor(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		w.WriteHeader(http.StatusUnauthorized)
		HandleActorError(w, err)
		return
	}
	req := &api_models.UpdateActorRequest{}

	err = json.NewDecoder(r.Body).Decode(req)
	err = validate.Struct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}

	actorID := chi.URLParam(r, "actorID")

	intActorID, err := strconv.ParseInt(actorID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}
	datetime, err := time.Parse("2006-01-02", req.Birth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}
	actor, err := db.UpdateActor(pgdb, &db.Actor{
		ID:    intActorID,
		Name:  req.Name,
		Sex:   req.Sex,
		Birth: datetime,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
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

// deleteActor godoc
// @Summary      Delete actor
// @Description  delete actor using id from request params
// @Tags         actors
// @Accept       json
// @Produce      json
// @Param actorID query string true "Actors Id"
// @Router       /actors [delete]
// @Security BasicAuth
// @Success 200 {object} nil
// @Failure 401 {object} http.Response
func deleteActor(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		w.WriteHeader(http.StatusUnauthorized)
		HandleActorError(w, err)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}

	actorID := chi.URLParam(r, "actorID")
	intActorID, err := strconv.ParseInt(actorID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}

	err = db.DeleteActor(pgdb, intActorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleActorError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
