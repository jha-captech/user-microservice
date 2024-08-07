package server

import "net/http"

func (h *Handler) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("Health check called")
		encodeResponse(w, http.StatusOK, map[string]string{"message": "hello world"})
	}
}
