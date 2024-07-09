package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jha-captech/user-microservice/internal/models"
)

// ListUsersHandler returns a list of all users from the database.
func (h *Handler) ListUsersHandler() APIGatewayHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get values from db
		users, err := h.userService.ListUsers()
		if err != nil {
			return h.returnErr("Not a valid ID", err, http.StatusInternalServerError)
		}

		// return response
		return h.returnJSON(http.StatusOK, ResponseAllUsers{Users: users})
	}
}

// FetchUsersHandler returns a single user based on an id passed as a path parameter on the request.
func (h *Handler) FetchUsersHandler() APIGatewayHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get and validate ID
		idString := request.PathParameters["ID"]
		ID, err := strconv.Atoi(idString)
		if err != nil {
			return h.returnErr("Not a valid ID", err, http.StatusBadRequest)
		}

		// get values from db
		user, err := h.userService.FetchUser(ID)
		if err != nil {
			return h.returnErr("Error retrieving data", err, http.StatusInternalServerError)
		}

		// return response
		return h.returnJSON(http.StatusOK, ResponseOneUser{User: user})
	}
}

// UpdateUsersHandler updates a user by ID.
func (h *Handler) UpdateUsersHandler() APIGatewayHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get and validate ID
		idString := request.PathParameters["ID"]
		ID, err := strconv.Atoi(idString)
		if err != nil {
			return h.returnErr("Not a valid ID", err, http.StatusBadRequest)
		}

		// get and validate body as object
		var inputUser models.User
		err = json.Unmarshal([]byte(request.Body), &inputUser)
		if err != nil {
			return h.returnErr("Missing values or malformed body", err, http.StatusBadRequest)
		}

		// update object in database
		user, err := h.userService.UpdateUser(ID, inputUser)
		if err != nil {
			return h.returnErr("Error updating object", err, http.StatusInternalServerError)
		}

		// return response
		return h.returnJSON(http.StatusOK, ResponseOneUser{User: user})
	}
}

// CreateUsersHandler creates a new user.
func (h *Handler) CreateUsersHandler() APIGatewayHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get and validate body as object
		var inputUser models.User
		err := json.Unmarshal([]byte(request.Body), &inputUser)
		if err != nil {
			return h.returnErr("missing values or malformed body", err, http.StatusBadRequest)
		}

		// create object in database
		ID, err := h.userService.CreateUser(inputUser)
		if err != nil {
			return h.returnErr("Error updating object", err, http.StatusInternalServerError)
		}

		// return response
		return h.returnJSON(http.StatusOK, ResponseID{ObjectID: ID})
	}
}

// DeleteUsersHandler deletes a user by ID.
func (h *Handler) DeleteUsersHandler() APIGatewayHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// get and validate ID
		idString := request.PathParameters["ID"]
		ID, err := strconv.Atoi(idString)
		if err != nil {
			return h.returnErr("Not a valid ID", err, http.StatusBadRequest)
		}

		// check that object exists
		user, err := h.userService.FetchUser(ID)
		if err != nil {
			return h.returnErr("Error validating object", err, http.StatusInternalServerError)
		}
		if user.ID == 0 {
			return h.returnErr("Object does not exist", err, http.StatusInternalServerError)
		}

		// delete user
		if err = h.userService.DeleteUser(ID); err != nil {
			return h.returnErr("Error deleting object.", err, http.StatusInternalServerError)
		}

		// return response
		return h.returnJSON(http.StatusOK, ResponseMessage{Message: "object successful deleted"})
	}
}
