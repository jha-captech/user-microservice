package handlers

import (
	"log/slog"
	"net/http"
)

func HandleHealthCheck(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Health check called")
		encodeResponse(w, logger, http.StatusOK, responseMsg{
			Message: "hello world",
		})
	}
}
