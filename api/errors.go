package api

import (
	"encoding/json"

	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Error msg"`
}

func HandleError(w http.ResponseWriter, err error) {
	res := &ErrorResponse{
		Success: false,
		Error:   err.Error(),
	}
	slog.Error("error %s\n", err)

	json.NewEncoder(w).Encode(res)
}
