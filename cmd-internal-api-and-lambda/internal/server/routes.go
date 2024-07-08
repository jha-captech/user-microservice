package server

import (
	"net/http"
)

type Options func(*routesOptions)

type routesOptions struct {
	useHealthCheck bool
}

func WithEnableHealthCheck(enableHealthCheck bool) Options {
	return func(options *routesOptions) {
		options.useHealthCheck = enableHealthCheck
	}
}

func RegisterRoutes(mux *http.ServeMux, h Handler, options ...Options) {
	opts := routesOptions{
		useHealthCheck: true,
	}

	for _, fn := range options {
		fn(&opts)
	}

	if opts.useHealthCheck {
		mux.HandleFunc("GET /api/health-check", h.handleHealthCheck())
	}

	mux.HandleFunc("GET /api/user", h.handleListUsers())
	mux.HandleFunc("GET /api/user/{id}", h.handleFetchUser())
	mux.HandleFunc("PUT /api/user/{id}", h.handleUpdateUser())
	mux.HandleFunc("POST /api/user", h.handleCreateUser())
	mux.HandleFunc("DELETE /api/user/{id}", h.handleDeleteUser())
}
