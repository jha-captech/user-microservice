package handlers

import (
	"log/slog"

	"github.com/jha-captech/user-microservice/internal/service"
)

// ── Handler Struct And Constructor ───────────────────────────────────────────────────────────────

type Handler struct {
	logger      *slog.Logger
	userService service.Service
}

func New(logger *slog.Logger, us service.Service) Handler {
	return Handler{
		logger:      logger,
		userService: us,
	}
}
