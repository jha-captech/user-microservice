package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jha-captech/user-microservice/internal/models"
)

type responseUser struct {
	User models.User `json:"user"`
}

type responseUsers struct {
	Users []models.User `json:"users"`
}

type responseMsg struct {
	Message string `json:"message"`
}

type responseID struct {
	ObjectID int `json:"object_id"`
}

type responseErr struct {
	Error string `json:"error"`
}

// encodeResponse encodes a struct of type T as a JSON response.
func encodeResponse(w http.ResponseWriter, logger *slog.Logger, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Error while marshaling data", "err", err, "data", data)
		http.Error(w, `{"Error": "Internal server error"}`, http.StatusInternalServerError)
	}
}
