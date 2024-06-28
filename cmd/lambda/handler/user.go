package handler

import (
	"net/http"
	"strconv"

	"user-microservice/internal/database/entity"
)

type responseOneUser struct {
	User entity.User `json:"user"`
}

type responseAllUsers struct {
	Users []entity.User `json:"users"`
}

// listUsers returns a list of all users from the database.
func (h Handler) listUsers(request Request) (Response, error) {
	// get values from db
	users, err := h.userService.List()
	if err != nil {
		h.logger.Error("error getting all locations", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Error retrieving data"})
		return Response{StatusCode: http.StatusInternalServerError, Body: respBody}, err
	}

	// return response
	respBody, _ := structToJSON(responseAllUsers{Users: users})
	return Response{StatusCode: http.StatusOK, Body: respBody}, err
}

// fetchUser returns a single user based on an id passed as a path parameter on the request.
func (h Handler) fetchUser(request Request) (Response, error) {
	// get and validate ID
	idString := request.PathParameters["ID"]
	ID, err := strconv.Atoi(idString)
	if err != nil {
		h.logger.Error("error getting ID", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Not a valid ID"})
		return Response{StatusCode: http.StatusInternalServerError, Body: respBody}, nil
	}

	// get values from db
	user, err := h.userService.Fetch(ID)
	if err != nil {
		h.logger.Error("error getting all locations", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Error retrieving data"})
		return Response{StatusCode: http.StatusInternalServerError, Body: respBody}, nil
	}

	// return response
	respBody, _ := structToJSON(responseOneUser{User: user})
	return Response{StatusCode: http.StatusOK, Body: respBody}, nil
}

func (h Handler) updateUser(request Request) (Response, error) {
	return Response{}, nil
}

func (h Handler) createUser(request Request) (Response, error) {
	return Response{}, nil
}

func (h Handler) deleteUser(request Request) (Response, error) {
	return Response{}, nil
}
