package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func notFound(r chi.Router, h Handler) {
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("Route not found called")
		encode(w, http.StatusNotFound, responseMessage{Message: "Page not found"})
	})
}
