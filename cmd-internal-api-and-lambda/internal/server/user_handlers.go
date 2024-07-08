package server

import (
	"net/http"
	"strconv"

	"github.com/jha-captech/user-microservice/internal/models"
)

type responseOneUser struct {
	User models.User `json:"user"`
}

type responseAllUsers struct {
	Users []models.User `json:"users"`
}

type responseMessage struct {
	Message string `json:"message"`
}

type responseID struct {
	ObjectID int `json:"object_id"`
}

type responseError struct {
	Error string `json:"error"`
}

// handleListUsers is a Handler that returns a list of all users.
func (h *Handler) handleListUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get values from Database
		users, err := h.userService.ListUsers()
		if err != nil {
			h.logger.Error("error getting all locations", "error", err)
			encodeResponse(
				w,
				http.StatusInternalServerError,
				responseError{Error: "Error retrieving data"},
			)
			return
		}

		// return response
		encodeResponse(w, http.StatusOK, responseAllUsers{Users: users})
	}
}

// handleFetchUser is a Handler that returns a single user by ID.
func (h *Handler) handleFetchUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and validate ID
		idString := r.PathValue("id")
		ID, err := strconv.Atoi(idString)
		if err != nil {
			h.logger.Error("error getting ID", "error", err)
			encodeResponse(w, http.StatusBadRequest, responseError{Error: "Not a valid ID"})
			return
		}

		// get values from Database
		user, err := h.userService.FetchUser(ID)
		if err != nil {
			h.logger.Error("error getting locations buy ID", "ID", ID, "error", err)
			encodeResponse(
				w,
				http.StatusInternalServerError,
				responseError{Error: "Error retrieving data"},
			)
			return
		}

		// return response
		encodeResponse(w, http.StatusOK, responseOneUser{User: user})
	}
}

// handleUpdateUser is a Handler that updates a user based on a user object from the request body.
func (h *Handler) handleUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and validate ID
		idString := r.PathValue("id")
		ID, err := strconv.Atoi(idString)
		if err != nil {
			h.logger.Error("error getting ID", "error", err)
			encodeResponse(w, http.StatusBadRequest, responseError{Error: "Not a valid ID"})
			return
		}

		// get and validate body as object
		inputUser, err := decodeToStruct[models.User](r)
		if err != nil {
			h.logger.Error("BodyParser error", "error", err)
			encodeResponse(
				w,
				http.StatusBadRequest,
				responseError{Error: "missing values or malformed body"},
			)
			return
		}

		// update object in Database
		user, err := h.userService.UpdateUser(ID, inputUser)
		if err != nil {
			h.logger.Error("error updating object in Database", "error", err)
			encodeResponse(
				w,
				http.StatusInternalServerError,
				responseError{Error: "Error updating data"},
			)
			return
		}

		// return response
		encodeResponse(w, http.StatusOK, responseOneUser{User: user})
	}
}

// handleUpdateUser is a Handler that creates a user based on a user object from the request body.
func (h *Handler) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and validate body as object
		inputUser, err := decodeToStruct[models.User](r)
		if err != nil {
			h.logger.Error("BodyParser error", "error", err)
			encodeResponse(
				w,
				http.StatusBadRequest,
				responseError{Error: "missing values or malformed body"},
			)
			return
		}

		// create object in Database
		ID, err := h.userService.CreateUser(inputUser)
		if err != nil {
			h.logger.Error("error creating object to Database", "error", err)
			encodeResponse(
				w,
				http.StatusInternalServerError,
				responseError{Error: "Error creating object"},
			)
			return
		}

		// return response
		encodeResponse(w, http.StatusOK, responseID{ObjectID: ID})
	}
}

// handleUpdateUser is a Handler that deletes a user based on an ID.
func (h *Handler) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and validate ID
		idString := r.PathValue("id")
		ID, err := strconv.Atoi(idString)
		if err != nil {
			h.logger.Error("error getting ID", "error", err)
			encodeResponse(w, http.StatusBadRequest, responseError{Error: "Not a valid ID"})
			return
		}

		// check that object exists
		user, err := h.userService.FetchUser(ID)
		if err != nil {
			h.logger.Error("error getting object by ID", "error", err)
			encodeResponse(
				w,
				http.StatusInternalServerError,
				responseError{Error: "Error validating object"},
			)
			return
		}
		if user.ID == 0 {
			encodeResponse(w, http.StatusBadRequest, responseError{Error: "Object does not exist"})
			return
		}

		// delete user
		if err = h.userService.DeleteUser(ID); err != nil {
			h.logger.Error("error deleting object by ID", "ID", ID, "error", err)
			encodeResponse(
				w,
				http.StatusInternalServerError,
				responseError{Error: "Error deleting object."},
			)
			return
		}

		// return message
		encodeResponse(w, http.StatusOK, responseMessage{Message: "object successful deleted"})
	}
}
