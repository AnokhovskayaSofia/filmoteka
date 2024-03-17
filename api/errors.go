package api

import (
	"encoding/json"
	
	api_models "filmoteka/api/models"

	"log/slog"
	"net/http"
)

func HandleActorsError(w http.ResponseWriter, err error)  {
	res := &api_models.ActorsResponse{
		Success: false,
		Error:   err.Error(),
		Actors:   nil,
	}
	slog.Error("error getting actor %s\n", err)
	
	json.NewEncoder(w).Encode(res)
}

func HandleActorError(w http.ResponseWriter, err error)  {
	res := &api_models.ActorResponse{
		Success: false,
		Error:   err.Error(),
		Actor:   nil,
	}
	slog.Error("error getting actor %s\n", err)
	
	json.NewEncoder(w).Encode(res)
}