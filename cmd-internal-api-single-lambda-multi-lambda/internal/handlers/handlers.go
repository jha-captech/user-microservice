package handlers

import (
	"log/slog"

	"github.com/jha-captech/user-microservice/internal/user"
)

// ── Handler Struct And Constructor ───────────────────────────────────────────────────────────────

type Handler struct {
	logger      *slog.Logger
	userService user.Service
}

func New(logger *slog.Logger, us user.Service) Handler {
	return Handler{
		logger:      logger,
		userService: us,
	}
}
