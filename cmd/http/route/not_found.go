package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func notFound(r chi.Router) {
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		encode(w, http.StatusNotFound, responseMessage{Message: "Page not found"})
	})
}
