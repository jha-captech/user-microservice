package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/models"
)

type userFetcher interface {
	FetchUser(ctx context.Context, ID int) (models.User, error)
}

// HandleFetchUser is a Handler that returns a single user by ID.
//
// @Summary		Fetch a user by ID
// @Description	Fetch a user by ID
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		id			path		int	true				"User ID"
// @Success		200			{object}	handlers.responseUser
// @Failure		400			{object}	handlers.responseErr
// @Failure		500			{object}	handlers.responseErr
// @Router		/user/{ID}	[GET]
func HandleFetchUser(logger sLogger, service userFetcher) http.HandlerFunc {
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

		// get values from database
		user, err := service.FetchUser(ctx, ID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				// no user found
				encodeResponse(w, logger, http.StatusOK, responseUser{})
			default:
				logger.Error("error getting object by ID", "error", err)
				encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
					Error: "Internal server error",
				})
			}
			return
		}

		// return response
		userOut := mapOutput(user)
		encodeResponse(w, logger, http.StatusOK, responseUser{
			User: userOut,
		})
	}
}
