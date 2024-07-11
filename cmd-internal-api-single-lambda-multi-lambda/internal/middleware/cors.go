package middleware

import (
	"net/http"
	"strings"
)

type CORSOptions struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

func CORSMiddleware(opts CORSOptions) Middleware {
	// set defaults if not set
	if opts.AllowedOrigins == nil {
		opts.AllowedOrigins = []string{"*"}
	}
	if opts.AllowedMethods == nil {
		opts.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE"}
	}
	if opts.AllowedHeaders == nil {
		opts.AllowedHeaders = []string{"Content-Type", "Authorization"}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(opts.AllowedOrigins, ","))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(opts.AllowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(opts.AllowedHeaders, ","))

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
