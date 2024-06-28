package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"user-microservice/internal/database/entity"
)

type (
	Request  = events.APIGatewayProxyRequest
	Response = events.APIGatewayProxyResponse

	APIGatewayHandler func(context.Context, Request) (Response, error)
)

type responseError struct {
	Error string `json:"error"`
}

type userService interface {
	List() ([]entity.User, error)
	Fetch(int) (entity.User, error)
	Update(int, entity.User) (entity.User, error)
	Create(entity.User) (int, error)
	Delete(int) error
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

func Run(h Handler) APIGatewayHandler {
	return func(ctx context.Context, request Request) (Response, error) {
		switch request.HTTPMethod {
		case http.MethodGet:
			if request.PathParameters != nil {
				return h.fetchUser(request)
			}
			return h.fetchUser(request)

		case http.MethodPut:
			return h.updateUser(request)

		case http.MethodPost:
			return h.createUser(request)

		case http.MethodDelete:
			return h.deleteUser(request)

		default:
			return Response{StatusCode: http.StatusMethodNotAllowed}, nil
		}
	}
}
