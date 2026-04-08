package middleware

import (
	"context"
	"net/http"
	"time"
)

// TimeoutMiddleware adds request timeout
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Create a request with the new context
			req := r.WithContext(ctx)
			handler.ServeHTTP(w, req)
		})
	}
}
