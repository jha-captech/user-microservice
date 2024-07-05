package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	"user-microservice/internal/database/entity"
)

// ── Constants And Errors ─────────────────────────────────────────────────────────────────────────

// ── Return Structs ───────────────────────────────────────────────────────────────────────────────

type responseError struct {
	Error string `json:"error"`
}

type responseMessage struct {
	Message string `json:"message"`
}

type responseID struct {
	ObjectID int `json:"object_id"`
}

type responseOneUser struct {
	User entity.User `json:"user"`
}

type responseAllUsers struct {
	Users []entity.User `json:"users"`
}

// ── Method Handlers ──────────────────────────────────────────────────────────────────────────────

// listUsers returns a list of all users from the database.
func (h Handler) listUsers(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get values from db
	users, err := h.userService.List()
	if err != nil {
		h.logger.Error("error getting all locations", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Error retrieving data"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, err
	}

	// return response
	respBody, _ := structToJSON(responseAllUsers{Users: users})
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: respBody}, err
}

// fetchUser returns a single user based on an id passed as a path parameter on the request.
func (h Handler) fetchUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get and validate ID
	idString := request.PathParameters["ID"]
	ID, err := strconv.Atoi(idString)
	if err != nil {
		h.logger.Error("error getting ID", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Not a valid ID"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, nil
	}

	// get values from db
	user, err := h.userService.Fetch(ID)
	if err != nil {
		h.logger.Error("error getting all locations", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Error retrieving data"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, nil
	}

	// return response
	respBody, _ := structToJSON(responseOneUser{User: user})
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: respBody}, nil
}

// updateUser updates a user by ID.
func (h Handler) updateUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get and validate ID
	idString := request.PathParameters["ID"]
	ID, err := strconv.Atoi(idString)
	if err != nil {
		h.logger.Error("error getting ID", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Not a valid ID"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, nil
	}

	// get and validate body as object
	var inputUser entity.User
	err = json.Unmarshal([]byte(request.Body), &inputUser)
	if err != nil {
		h.logger.Error("Error Unmarshalling request body", "error", err)
		respBody, _ := structToJSON(responseError{Error: "missing values or malformed body"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, nil
	}

	// update object in database
	user, err := h.userService.Update(ID, inputUser)
	if err != nil {
		h.logger.Error("error updating object in db", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Error updating data"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, nil
	}

	// return response
	respBody, _ := structToJSON(responseOneUser{User: user})
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: respBody}, nil
}

// createUser creates a new user.
func (h Handler) createUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get and validate body as object
	var inputUser entity.User
	err := json.Unmarshal([]byte(request.Body), &inputUser)
	if err != nil {
		h.logger.Error("Error Unmarshalling request body", "error", err)
		respBody, _ := structToJSON(responseError{Error: "missing values or malformed body"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, nil
	}

	// create object in database
	id, err := h.userService.Create(inputUser)
	if err != nil {
		h.logger.Error("error creating object in db", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Error creating object"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, nil
	}

	// return response
	responseBody, _ := json.Marshal(responseID{ObjectID: id})
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}

// deleteUser deletes a user by ID.
func (h Handler) deleteUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get and validate ID
	idString := request.PathParameters["ID"]
	ID, err := strconv.Atoi(idString)
	if err != nil {
		h.logger.Error("error getting ID", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Not a valid ID"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, nil
	}

	// check that object exists
	user, err := h.userService.Fetch(ID)
	if err != nil {
		h.logger.Error("error getting object by ID", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Error validating object"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, nil
	}
	if user.ID == 0 {
		respBody, _ := structToJSON(responseError{Error: "Object does not exist"})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, nil
	}

	// delete user
	if err = h.userService.Delete(ID); err != nil {
		h.logger.Error("error deleting object by ID", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Error deleting object."})
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       respBody,
		}, nil
	}

	// return response
	responseBody, _ := json.Marshal(responseMessage{Message: "object successful deleted"})
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}
