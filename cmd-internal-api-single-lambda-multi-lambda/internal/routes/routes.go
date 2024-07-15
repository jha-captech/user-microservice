package routes

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/handlers"
	"github.com/jha-captech/user-microservice/internal/service"
)

type Option func(*routerOptions)

type routerOptions struct {
	registerHealthRoute bool
}

// WithRegisterHealthRoute controls whether a healthcheck route will be registered. If `false` is
// passed in or this function is not called, the default is `false`.
func WithRegisterHealthRoute(registerHealthRoute bool) Option {
	return func(options *routerOptions) {
		options.registerHealthRoute = registerHealthRoute
	}
}

func RegisterRoutes(r *chi.Mux, logger *slog.Logger, svs *service.User, opts ...Option) {
	options := routerOptions{
		registerHealthRoute: false,
	}
	for _, opt := range opts {
		opt(&options)
	}

	if options.registerHealthRoute {
		r.Get("/api/health-check", handlers.HandleHealth(logger))
	}

	r.Get("/api/user", handlers.HandleListUsers(logger, svs))
	r.Get("/api/user/{ID}", handlers.HandleFetchUser(logger, svs))
	r.Put("/api/user/{ID}", handlers.HandleUpdateUser(logger, svs))
	r.Post("/api/user", handlers.HandleCreateUser(logger, svs))
	r.Delete("/api/user/{ID}", handlers.HandleDeleteUser(logger, svs))
}
