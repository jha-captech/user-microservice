package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/models"
	"github.com/jha-captech/user-microservice/internal/service"
)

// HandleListUsers is a Handler that returns a list of all users.
func HandleListUsers(logger *slog.Logger, service *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get values from database
		users, err := service.ListUsers()
		if err != nil {
			logger.Error("error getting all locations", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		// return response
		encodeResponse(w, logger, http.StatusOK, responseUsers{
			Users: users,
		})
	}
}

// HandleFetchUser is a Handler that returns a single user by ID.
func HandleFetchUser(logger *slog.Logger, service *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and validate ID
		idString := chi.URLParam(r, "ID")
		ID, err := strconv.Atoi(idString)
		if err != nil {
			logger.Error("error getting ID", "error", err)
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "Not a valid ID",
			})
			return
		}

		// get values from database
		user, err := service.FetchUser(ID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				encodeResponse(w, logger, http.StatusOK, responseUser{})
			default:
				logger.Error("error getting object by ID", "error", err)
				encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
					Error: "Error validating object",
				})
			}
			return
		}

		// return response
		encodeResponse(w, logger, http.StatusOK, responseUser{
			User: user,
		})
	}
}

// HandleUpdateUser is a Handler that updates a user based on a user object from the request body.
func HandleUpdateUser(logger *slog.Logger, service *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and validate ID
		idString := chi.URLParam(r, "ID")
		ID, err := strconv.Atoi(idString)
		if err != nil {
			logger.Error("error getting ID", "error", err)
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "Not a valid ID",
			})
			return
		}

		// get and validate body as object
		inputUser, err := decodeRequestBody[models.User](r)
		if err != nil {
			logger.Error("BodyParser error", "error", err)
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "missing values or malformed body",
			})
			return
		}

		// update object in database
		user, err := service.UpdateUser(ID, inputUser)
		if err != nil {
			logger.Error("error updating object in database", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error updating data",
			})
			return
		}

		// return response
		encodeResponse(w, logger, http.StatusOK, responseUser{
			User: user,
		})
	}
}

// HandleCreateUser is a Handler that creates a user based on a user object from the request body.
func HandleCreateUser(logger *slog.Logger, service *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and validate body as object
		inputUser, err := decodeRequestBody[models.User](r)
		if err != nil {
			logger.Error("BodyParser error", "error", err)
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "missing values or malformed body",
			})
			return
		}

		// create object in database
		ID, err := service.CreateUser(inputUser)
		if err != nil {
			logger.Error("error creating object to database", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error creating object",
			})
			return
		}

		// return response
		encodeResponse(w, logger, http.StatusOK, responseID{
			ObjectID: ID,
		})
	}
}

// HandleDeleteUser is a Handler that deletes a user based on an ID.
func HandleDeleteUser(logger *slog.Logger, service *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and validate ID
		idString := chi.URLParam(r, "ID")
		ID, err := strconv.Atoi(idString)
		if err != nil {
			logger.Error("error getting ID", "error", err)
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "Not a valid ID",
			})
			return
		}

		// check that object exists
		_, err = service.FetchUser(ID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				logger.Error("Object does not exist", "error", err)
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					Error: "Object does not exist",
				})
			default:
				logger.Error("error getting object by ID", "error", err)
				encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
					Error: "Error validating object",
				})
			}
			return
		}

		// delete user
		if err = service.DeleteUser(ID); err != nil {
			logger.Error("error deleting object by ID", "ID", ID, "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error deleting object.",
			})
			return
		}

		// return message
		encodeResponse(w, logger, http.StatusOK, responseMsg{
			Message: "object successful deleted",
		})
	}
}
