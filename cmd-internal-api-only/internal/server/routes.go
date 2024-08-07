package server

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h Handler) {
	mux.HandleFunc("GET /api/health-check", h.handleHealthCheck())
	mux.HandleFunc("GET /api/user", h.handleListUsers())
	mux.HandleFunc("GET /api/user/{id}", h.handleFetchUser())
	mux.HandleFunc("PUT /api/user/{id}", h.handleUpdateUser())
	mux.HandleFunc("POST /api/user", h.handleCreateUser())
	mux.HandleFunc("DELETE /api/user/{id}", h.handleDeleteUser())
}
