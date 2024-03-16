package api_models

import db_models "filmoteka/db"

type ActorsResponse struct {
	Success bool               `json:"success"`
	Error   string             `json:"error"`
	Actors  []*db_models.Actor `json:"actors"`
}

type ActorResponse struct {
	Success bool             `json:"success"`
	Error   string           `json:"error"`
	Actor   *db_models.Actor `json:"actor"`
}

type CreateActorRequest struct {
	Name  string `json:"name"`
	Sex   string `json:"sex" validate:"oneof=male female"`
	Birth string `json:"birth"`
}

type UpdateActorRequest struct {
	Name  string `json:"name,omitempty"`
	Sex   string `json:"sex,omitempty"`
	Birth string `json:"birth,omitempty"`
}
