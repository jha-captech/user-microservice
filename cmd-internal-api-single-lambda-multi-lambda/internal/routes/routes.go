package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/handlers"
)

func RegisterRoutes(r *chi.Mux, h handlers.Handler) {
	r.Get("/api/health-check", h.HandleHealthCheck())

	r.Get("/api/user", h.HandleListUsers())
	r.Get("/api/user/{ID}", h.HandleFetchUser())
	r.Put("/api/user/{ID}", h.HandleUpdateUser())
	r.Post("/api/user", h.HandleCreateUser())
	r.Delete("/api/user/{ID}", h.HandleDeleteUser())
}
