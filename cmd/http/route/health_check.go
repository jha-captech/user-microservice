package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func healthCheck(h Handler) func(r chi.Router) {
	return func(r chi.Router) {
		// @Summary		Health check response
		// @Description	Health check response
		// @Tags		health-check
		// @Accept		json
		// @Produce		json
		// @Success		200				{object}	routes.responseMessage
		// @Router		/health-check 	[GET]
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			h.logger.Info("Health check called")
			encode(w, http.StatusOK, responseMessage{Message: "Hello World"})
		})
	}
}

// @Summary		Health check response
// @Description	Health check response
// @Tags		health-check
// @Accept		json
// @Produce		json
// @Success		200					{object}	routes.responseMessage
// @Router		/health-check/v2	[GET]
func handleHealthCheck(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("Health check v2 called")
		encode(w, http.StatusOK, responseMessage{Message: "Hello World"})
	}
}
