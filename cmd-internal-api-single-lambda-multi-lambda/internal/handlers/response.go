package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jha-captech/user-microservice/internal/log"
	"github.com/jha-captech/user-microservice/internal/models"
)

type outputUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	UserID    int    `json:"user_id"`
}

func mapOutput(user models.User) outputUser {
	return outputUser{
		ID:        int(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		UserID:    int(user.UserID),
	}
}

func mapMultipleOutput(user []models.User) []outputUser {
	usersOut := make([]outputUser, len(user))
	for i := 0; i < len(user); i++ {
		userOut := mapOutput(user[i])
		usersOut[i] = userOut
	}

	return usersOut
}

type responseUser struct {
	User outputUser `json:"user"`
}

type responseUsers struct {
	Users []outputUser `json:"users"`
}

type responseMsg struct {
	Message string `json:"message"`
}

type responseID struct {
	ObjectID int `json:"object_id"`
}

type responseErr struct {
	Error            string            `json:"error,omitempty"`
	ValidationErrors map[string]string `json:"validation_errors,omitempty"`
}

// encodeResponse encodes data as a JSON response.
func encodeResponse(w http.ResponseWriter, logger log.Logger, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Error while marshaling data", "err", err, "data", data)
		http.Error(w, `{"Error": "Internal server error"}`, http.StatusInternalServerError)
	}
}
