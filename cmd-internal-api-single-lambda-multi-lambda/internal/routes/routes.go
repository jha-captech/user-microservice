package routes

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/handlers"
	"github.com/jha-captech/user-microservice/internal/service"
)

func RegisterRoutes(r *chi.Mux, logger *slog.Logger, svs *service.Service) {
	r.Get("/api/health-check", handlers.HandleHealthCheck(logger))

	r.Get("/api/user", handlers.HandleListUsers(logger, svs))
	r.Get("/api/user/{ID}", handlers.HandleFetchUser(logger, svs))
	r.Put("/api/user/{ID}", handlers.HandleUpdateUser(logger, svs))
	r.Post("/api/user", handlers.HandleCreateUser(logger, svs))
	r.Delete("/api/user/{ID}", handlers.HandleDeleteUser(logger, svs))
}
