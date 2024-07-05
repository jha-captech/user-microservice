package route

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"user-microservice/internal/database/entity"
)

type responseMessage struct {
	Message string `json:"message"`
}

type responseID struct {
	ObjectID int `json:"object_id"`
}

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

// NewHandler creates and returns a new Handler struct.
func NewHandler(userService userService, logger *slog.Logger) Handler {
	return Handler{
		userService: userService,
		logger:      logger,
	}
}

// SetUpRoutes sets up routes using a *chi.Mux.
func SetUpRoutes(r *chi.Mux, h Handler) {
	r.Route("/api", func(r chi.Router) {
		// For demonstration purposes, there are two versions of health check that show the two ways of
		// using http handlers with chi and the standard library.
		// "/health-check" uses an anonymous handler func inside a route function
		// "/health-check/v2" uses a names function as a closure for a http.HandlerFunc
		r.Route("/health-check", healthCheck(h))

		r.Route("/health-check/v2", func(r chi.Router) {
			r.Get("/", handleHealthCheck(h))
		})

		r.Route("/user", userRoutes(h))
	})

	notFound(r, h)
}
