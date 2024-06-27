package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	"user-microservice/internal/database/entity"
)

type (
	Request  = events.APIGatewayProxyRequest
	Response = events.APIGatewayProxyResponse

	APIGatewayHandler func(context.Context, Request) (Response, error)
)

type userService interface {
	List() ([]entity.User, error)
	Fetch(ID int) (entity.User, error)
}

type Handler struct {
	userService userService
	logger      *slog.Logger
}

func New(userService userService, logger *slog.Logger) Handler {
	return Handler{
		userService: userService,
		logger:      logger,
	}
}

func Run(handler Handler) APIGatewayHandler {
	return func(ctx context.Context, request Request) (Response, error) {
		switch {
		case request.HTTPMethod == http.MethodGet && request.PathParameters != nil:
			return handler.fetchUser(request)
		case request.HTTPMethod == http.MethodPost:
			return handler.fetchUser(request)
		default:
			return Response{StatusCode: http.StatusMethodNotAllowed}, nil
		}
	}
}

type responseError struct {
	Error string `json:"error"`
}
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
		return Response{StatusCode: http.StatusInternalServerError, Body: respBody}, err
	}

	// get values from db
	user, err := h.userService.Fetch(ID)
	if err != nil {
		h.logger.Error("error getting all locations", "error", err)
		respBody, _ := structToJSON(responseError{Error: "Error retrieving data"})
		return Response{StatusCode: http.StatusInternalServerError, Body: respBody}, err
	}

	// return response
	respBody, _ := structToJSON(responseOneUser{User: user})
	return Response{StatusCode: http.StatusOK, Body: respBody}, err
}
