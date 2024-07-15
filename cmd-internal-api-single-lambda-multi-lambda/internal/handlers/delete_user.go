package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/models"
)

type deleteUserServicer interface {
	FetchUser(ID int) (models.User, error)
	DeleteUser(ID int) error
}

// HandleDeleteUser is a Handler that deletes a user based on an ID.
//
// @Summary		Delete a user by ID
// @Description	Delete a user by ID
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		id			path		int	true				"User ID"
// @Success		202			{object}	handlers.responseMsg
// @Failure		400			{object}	handlers.responseErr
// @Failure		500			{object}	handlers.responseErr
// @Router		/user/{ID}	[DELETE]
func HandleDeleteUser(logger *slog.Logger, service deleteUserServicer) http.HandlerFunc {
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
		encodeResponse(w, logger, http.StatusAccepted, responseMsg{
			Message: "object successful deleted",
		})
	}
}
