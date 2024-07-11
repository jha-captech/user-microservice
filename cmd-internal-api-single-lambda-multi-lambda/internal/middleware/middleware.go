package middleware

import (
	"net/http"
)

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
