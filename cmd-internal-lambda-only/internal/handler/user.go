package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jha-captech/user-microservice/internal/models"
)

// ListUsersHandler returns a list of all users from the database.
func (h *Handler) ListUsersHandler() APIGatewayHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get values from db
		users, err := h.UserService.ListUsers()
		if err != nil {
			h.logger.Error("Encountered error while getting objects from the database", "err", err)
			return h.returnJSON(http.StatusInternalServerError, ResponseError{
				Error: "Internal server error",
			})
		}

		// return response
		return h.returnJSON(http.StatusOK, ResponseAllUsers{
			Users: users,
		})
	}
}

// FetchUsersHandler returns a single user based on an id passed as a path parameter on the request.
func (h *Handler) FetchUsersHandler() APIGatewayHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get and validate ID
		idString := request.PathParameters["ID"]
		ID, err := strconv.Atoi(idString)
		if err != nil {
			h.logger.Error("Failed to parse ID from path paramaters", "err", err)
			return h.returnJSON(http.StatusBadRequest, ResponseError{
				Error: "Not a valid ID",
			})
		}

		// get value from db
		user, err := h.UserService.FetchUser(ID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return h.returnJSON(http.StatusOK, ResponseOneUser{
					User: models.User{},
				})
			default:
				h.logger.Error("Encountered error while getting object from the database", "err", err)
				return h.returnJSON(http.StatusInternalServerError, ResponseError{
					Error: "Internal server error",
				})
			}
		}

		// return response
		return h.returnJSON(http.StatusOK, ResponseOneUser{
			User: user,
		})
	}
}

// UpdateUsersHandler updates a user by ID.
func (h *Handler) UpdateUsersHandler() APIGatewayHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get and validate ID
		idString := request.PathParameters["ID"]
		ID, err := strconv.Atoi(idString)
		if err != nil {
			h.logger.Error("Failed to parse ID from path paramaters", "err", err)
			return h.returnJSON(http.StatusBadRequest, ResponseError{
				Error: "Not a valid ID",
			})
		}

		// get and validate body as object
		var inputUser models.User
		err = json.Unmarshal([]byte(request.Body), &inputUser)
		if err != nil {
			h.logger.Error("Failed to unmarshal request body", "err", err)
			return h.returnJSON(http.StatusBadRequest, ResponseError{
				Error: "Missing values or malformed body",
			})
		}

		// update object in db
		user, err := h.UserService.UpdateUser(ID, inputUser)
		if err != nil {
			h.logger.Error("Encountered error while updating object in the database", "err", err)
			return h.returnJSON(http.StatusInternalServerError, ResponseError{
				Error: "Internal server error",
			})
		}

		// return response
		return h.returnJSON(http.StatusOK, ResponseOneUser{
			User: user,
		})
	}
}

// CreateUsersHandler creates a new user.
func (h *Handler) CreateUsersHandler() APIGatewayHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get and validate body as object
		var inputUser models.User
		err := json.Unmarshal([]byte(request.Body), &inputUser)
		if err != nil {
			h.logger.Error("Failed to unmarshal request body", "err", err)
			return h.returnJSON(http.StatusBadRequest, ResponseError{
				Error: "Missing values or malformed body",
			})
		}

		// create object in db
		ID, err := h.UserService.CreateUser(inputUser)
		if err != nil {
			h.logger.Error("Encountered error while creating object in the database", "err", err)
			return h.returnJSON(http.StatusInternalServerError, ResponseError{
				Error: "Internal server error",
			})
		}

		// return response
		return h.returnJSON(http.StatusOK, ResponseID{
			ObjectID: ID,
		})
	}
}

// DeleteUsersHandler deletes a user by ID.
func (h *Handler) DeleteUsersHandler() APIGatewayHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get and validate ID
		idString := request.PathParameters["ID"]
		ID, err := strconv.Atoi(idString)
		if err != nil {
			h.logger.Error("Failed to parse ID from path paramaters", "err", err)
			return h.returnJSON(http.StatusBadRequest, ResponseError{
				Error: "Not a valid ID",
			})
		}

		// check that object exists
		_, err = h.UserService.FetchUser(ID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				h.logger.Error("Object with given ID does not exist", "ID", ID, "err", err)
				return h.returnJSON(http.StatusBadRequest, ResponseError{
					Error: "Internal server error",
				})
			default:
				h.logger.Error("Encountered error while validating object in the database", "err", err)
				return h.returnJSON(http.StatusInternalServerError, ResponseError{
					Error: "Internal server error",
				})
			}
		}

		// delete returnedUser from db
		if err = h.UserService.DeleteUser(ID); err != nil {
			h.logger.Error("Encountered error while deleting object from the database", "err", err)
			return h.returnJSON(http.StatusInternalServerError, ResponseError{
				Error: "Internal server error",
			})
		}

		// return response
		return h.returnJSON(http.StatusOK, ResponseMessage{
			Message: "object successful deleted",
		})
	}
}
