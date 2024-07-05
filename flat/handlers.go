package main

import (
	"log/slog"
	"net/http"
	"strconv"
)

// ── Handler Struct And Constructor ───────────────────────────────────────────────────────────────

type handler struct {
	logger  *slog.Logger
	service userService
}

func newHandler(logger *slog.Logger, us userService) handler {
	return handler{
		logger:  logger,
		service: us,
	}
}

// ── Healthcheck Handler ──────────────────────────────────────────────────────────────────────────

func (h *handler) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("Health check called")
		encodeResponse(w, http.StatusOK, map[string]string{"message": "hello world"})
	}
}

// ── User Handlers ────────────────────────────────────────────────────────────────────────────────

type responseOneUser struct {
	User User `json:"user"`
}

type responseAllUsers struct {
	Users []User `json:"users"`
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

// handleListUsers is a handler that returns a list of all users.
func (h *handler) handleListUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get values from db
		users, err := h.service.listUsers()
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

// handleFetchUser is a handler that returns a single user by ID.
func (h *handler) handleFetchUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and validate ID
		idString := r.PathValue("id")
		ID, err := strconv.Atoi(idString)
		if err != nil {
			h.logger.Error("error getting ID", "error", err)
			encodeResponse(w, http.StatusBadRequest, responseError{Error: "Not a valid ID"})
			return
		}

		// get values from db
		user, err := h.service.fetchUser(ID)
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

// handleUpdateUser is a handler that updates a user based on a user object from the request body.
func (h *handler) handleUpdateUser() http.HandlerFunc {
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
		inputUser, err := decodeToStruct[User](r)
		if err != nil {
			h.logger.Error("BodyParser error", "error", err)
			encodeResponse(
				w,
				http.StatusBadRequest,
				responseError{Error: "missing values or malformed body"},
			)
			return
		}

		// update object in database
		user, err := h.service.updateUser(ID, inputUser)
		if err != nil {
			h.logger.Error("error updating object in db", "error", err)
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

// handleUpdateUser is a handler that creates a user based on a user object from the request body.
func (h *handler) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and validate body as object
		inputUser, err := decodeToStruct[User](r)
		if err != nil {
			h.logger.Error("BodyParser error", "error", err)
			encodeResponse(
				w,
				http.StatusBadRequest,
				responseError{Error: "missing values or malformed body"},
			)
			return
		}

		// create object in database
		ID, err := h.service.createUser(inputUser)
		if err != nil {
			h.logger.Error("error creating object to db", "error", err)
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

// handleUpdateUser is a handler that deletes a user based on an ID.
func (h *handler) handleDeleteUser() http.HandlerFunc {
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
		user, err := h.service.fetchUser(ID)
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
		if err = h.service.deleteUser(ID); err != nil {
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
