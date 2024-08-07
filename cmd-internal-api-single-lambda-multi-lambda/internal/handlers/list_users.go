package handlers

import (
	"context"
	"net/http"

	"github.com/jha-captech/user-microservice/internal/models"
)

type userLister interface {
	ListUsers(ctx context.Context) ([]models.User, error)
}

// HandleListUsers is a Handler that returns a list of all users.
//
// @Summary		List all users
// @Description	List all users
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		200		{object}	handlers.responseUsers
// @Failure		500		{object}	handlers.responseErr
// @Router		/user	[GET]
func HandleListUsers(logger sLogger, service userLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get values from database
		users, err := service.ListUsers(ctx)
		if err != nil {
			logger.Error("error getting all locations", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		// return response
		usersOut := mapMultipleOutput(users)
		encodeResponse(w, logger, http.StatusOK, responseUsers{
			Users: usersOut,
		})
	}
}
