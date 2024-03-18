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

// getActors godoc
// @Summary      List actors
// @Description  Availible only for authenticated user, getting actors list from db
// @Tags         actors
// @Accept       json
// @Produce      json
// @Router       /actors [get]
// @Security BasicAuth
// @Success 200 {object} api_models.ActorsResponse
// @Failure 401 {object}  ErrorResponse
// @Failure 400 {object} ErrorResponse
func getActors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := checkBasicAuth(r)
	if err != nil {
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

	actors, err := db.GetActors(pgdb)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	res := &api_models.ActorsResponse{
		Success: true,
		Error:   "",
		Actors:  actors,
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

}

// createActor godoc
// @Summary      Create actor
// @Description  Availible only for admin user, creating actor using data from request body and return new actor
// @Tags         actors
// @Accept       json
// @Produce      json
// @Router       /actors [post]
// @Param Actor body db.Actor true "actor info"
// @Security BasicAuth
// @Success 200 {object} api_models.ActorResponse
// @Failure 401 {object}  ErrorResponse
// @Failure 400 {object} ErrorResponse
func createActor(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		w.WriteHeader(http.StatusUnauthorized)
		HandleError(w, err)
		return
	}
	req := &api_models.CreateActorRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	err = Validate.Struct(req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	pgdb, ok := r.Context().Value("DB").(*pg.DB)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, errors.New("could not get the DB from context"))
		return
	}

	birthday, err := time.Parse("2006-01-02", req.Birth)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}
	actor, err := db.CreateActor(pgdb, &db.Actor{
		Name:  req.Name,
		Sex:   req.Sex,
		Birth: birthday,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
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
		HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// updateActor godoc
// @Summary      Update actor
// @Description  Availible only for admin user, updating actor using id from request params and return actor
// @Tags         actors
// @Accept       json
// @Produce      json
// @Param actorID query string true "Actors Id"
// @Param Actor body db.Actor true "actor info"
// @Router       /actors [put]
// @Security BasicAuth
// @Success 200 {object} api_models.ActorResponse
// @Failure 401 {object}  ErrorResponse
// @Failure 400 {object} ErrorResponse
func updateActor(w http.ResponseWriter, r *http.Request) {
	auth_role, err := checkBasicAuth(r)
	if err != nil || auth_role != db.Admin {
		if err == nil {
			err = errors.New("wrong access level")
		}
		w.WriteHeader(http.StatusUnauthorized)
		HandleError(w, err)
		return
	}
	req := &api_models.UpdateActorRequest{}

	err = json.NewDecoder(r.Body).Decode(req)
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

	actorID := chi.URLParam(r, "actorID")

	intActorID, err := strconv.ParseInt(actorID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}
	var datetime time.Time
	if req.Birth != "" {
		datetime, err = time.Parse("2006-01-02", req.Birth)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			HandleError(w, err)
			return
		}
	}
	actor, err := db.UpdateActor(pgdb, &db.Actor{
		ID:    intActorID,
		Name:  req.Name,
		Sex:   req.Sex,
		Birth: datetime,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
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
// @Description  Availible only for admin user, deleting actor using id from request params
// @Tags         actors
// @Accept       json
// @Produce      json
// @Param actorID query string true "Actors Id"
// @Router       /actors [delete]
// @Security BasicAuth
// @Success 200 {object} nil
// @Failure 401 {object}  ErrorResponse
// @Failure 400 {object} ErrorResponse
func deleteActor(w http.ResponseWriter, r *http.Request) {
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

	actorID := chi.URLParam(r, "actorID")
	intActorID, err := strconv.ParseInt(actorID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	err = db.DeleteActor(pgdb, intActorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
