package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/handlers"
	"github.com/jha-captech/user-microservice/internal/service"
)

type sLogger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
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

func RegisterRoutes(router *chi.Mux, logger sLogger, svs *service.User, opts ...Option) {
	options := routerOptions{
		registerHealthRoute: false,
	}
	for _, opt := range opts {
		opt(&options)
	}

	if options.registerHealthRoute {
		router.Get("/api/health-check", handlers.HandleHealth(logger))
	}

	router.Get("/api/user", handlers.HandleListUsers(logger, svs))
	router.Get("/api/user/{ID}", handlers.HandleFetchUser(logger, svs))
	router.Put("/api/user/{ID}", handlers.HandleUpdateUser(logger, svs))
	router.Post("/api/user", handlers.HandleCreateUser(logger, svs))
	router.Delete("/api/user/{ID}", handlers.HandleDeleteUser(logger, svs))
}
