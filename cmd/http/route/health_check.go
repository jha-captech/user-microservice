package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func healthCheck(h Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			encode(w, http.StatusOK, responseMessage{Message: "Hello World"})
		})
	}
}

func handleHealthCheck(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		encode(w, http.StatusOK, responseMessage{Message: "Hello World"})
	}
}
