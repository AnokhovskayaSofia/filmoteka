package api_models

import (
	db_models "filmoteka/db"
)

type FilmsResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Films   []*db_models.Film `json:"film"`
}

type FilmResponse struct {
	Success bool            `json:"success"`
	Error   string          `json:"error"`
	Film    *db_models.Film `json:"film"`
}

type CreateFilmRequest struct {
	Name        string `json:"name" validate:"min=1,max=150"`
	Description string `json:"description" validate:"max=1000"`
	Date        string `json:"date"`
	Rate        int    `json:"rate" validate:"min=0,min=10"`
	Actors      []int  `json:"actors"`
}

type UpdateFilmRequest struct {
	Name        string `json:"name,omitempty" validate:"min=1,max=150"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	Date        string `json:"date,omitempty"`
	Rate        int    `json:"rate,omitempty" validate:"min=0,min=10"`
	Actors      []int  `json:"actors,omitempty"`
}
