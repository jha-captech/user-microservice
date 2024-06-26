package route

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"user-microservice/internal/database/entity"
)

type responseOneUser struct {
	User entity.User `json:"user"`
}

type responseAllUsers struct {
	Users []entity.User `json:"users"`
}

func userRoutes(h Handler) func(r chi.Router) {
	return func(r chi.Router) {
		// list all users
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			// get values from db
			users, err := h.userService.List()
			if err != nil {
				h.logger.Error("error getting all locations", "error", err)
				encode(
					w,
					http.StatusInternalServerError,
					responseError{Error: "Error retrieving data"},
				)
				return
			}

			// return response
			encode(w, http.StatusOK, responseAllUsers{Users: users})
		})

		// fetch a user by ID
		r.Get("/{ID}", func(w http.ResponseWriter, r *http.Request) {
			// get and validate ID
			idString := chi.URLParam(r, "ID")
			ID, err := strconv.Atoi(idString)
			if err != nil {
				h.logger.Error("error getting ID", "error", err)
				encode(w, http.StatusBadRequest, responseError{Error: "Not a valid ID"})
				return
			}

			// get values from db
			user, err := h.userService.Fetch(ID)
			if err != nil {
				h.logger.Error("error getting all locations", "error", err)
				encode(
					w,
					http.StatusInternalServerError,
					responseError{Error: "Error retrieving data"},
				)
				return
			}

			// return response
			encode(w, http.StatusOK, responseOneUser{User: user})
		})
	}
}
