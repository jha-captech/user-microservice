package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// ── Middleware Setup ─────────────────────────────────────────────────────────────────────────────

type Middleware func(http.Handler) http.Handler

func CreateStack(fn ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(fn) - 1; i >= 0; i-- {
			x := fn[i]
			next = x(next)
		}
		return next
	}
}

// ── Request Logger ───────────────────────────────────────────────────────────────────────────────

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func LoggerMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			logger.Info(
				"Incoming request",
				"status code",
				wrapped.statusCode,
				"elapses time",
				time.Since(start),
				"request address",
				r.RemoteAddr,
				"method",
				r.Method,
				"route",
				r.URL.Path,
			)
		})
	}
}

// ── Recovery ─────────────────────────────────────────────────────────────────────────────────────

func RecoveryMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					logger.Error("panic recovered", "panic", err)

					jsonBody, _ := json.Marshal(
						map[string]string{
							"error": "There was an internal server error",
						},
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(jsonBody)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// ── CORS ─────────────────────────────────────────────────────────────────────────────────────────

type CORSOptions struct {
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
}

func CORSMiddleware(opts CORSOptions) Middleware {
	// set defaults if not set
	if opts.allowedOrigins == nil {
		opts.allowedOrigins = []string{"*"}
	}
	if opts.allowedMethods == nil {
		opts.allowedMethods = []string{"GET", "POST", "PUT", "DELETE"}
	}
	if opts.allowedHeaders == nil {
		opts.allowedHeaders = []string{"Content-Type", "Authorization"}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(opts.allowedOrigins, ","))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(opts.allowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(opts.allowedHeaders, ","))

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
