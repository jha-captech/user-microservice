package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/log"
	"github.com/jha-captech/user-microservice/internal/models"
)

type userDeleter interface {
	FetchUser(ctx context.Context, ID int) (models.User, error)
	DeleteUser(ctx context.Context, ID int) error
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
func HandleDeleteUser(logger log.Logger, service userDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

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
		_, err = service.FetchUser(ctx, ID)
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
		if err = service.DeleteUser(ctx, ID); err != nil {
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
