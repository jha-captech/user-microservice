package handlers

import (
	"context"
	"net/http"

	"github.com/jha-captech/user-microservice/internal/log"
	"github.com/jha-captech/user-microservice/internal/models"
)

type userCreator interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
}

// HandleCreateUser is a Handler that creates a user based on a user object from the request body.
//
// @Summary		Create a user
// @Description	Create a user
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		user		body		handlers.inputUser	true	"User Object"
// @Success		201			{object}	handlers.responseID
// @Failure		400			{object}	handlers.responseErr
// @Failure		500			{object}	handlers.responseErr
// @Failure		409			{object}	handlers.responseErr
// @Router		/user		[POST]
func HandleCreateUser(logger log.Logger, service userCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup
		ctx := r.Context()

		// get and validate body as object
		userIn, problems, err := decodeValidateBody[inputUser, models.User](r)
		if err != nil {
			switch {
			case len(problems) > 0:
				logger.Error("Problems validating input", "error", err, "problems", problems)
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					ValidationErrors: problems,
				})
			default:
				logger.Error("BodyParser error", "error", err)
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					Error: "missing values or malformed body",
				})
			}
			return
		}

		// create object in database
		ID, err := service.CreateUser(ctx, userIn)
		if err != nil {
			logger.Error("error creating object to database", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error creating object",
			})
			return
		}

		// return response
		encodeResponse(w, logger, http.StatusCreated, responseID{
			ObjectID: ID,
		})
	}
}
