package handler

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jha-captech/user-microservice/internal/user"
)

type APIGatewayHandler func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Handler struct {
	logger      *slog.Logger
	userService user.Service
}

func NewHandler(logger *slog.Logger, us user.Service) Handler {
	return Handler{
		logger:      logger,
		userService: us,
	}
}
