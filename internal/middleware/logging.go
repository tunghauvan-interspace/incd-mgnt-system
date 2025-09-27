package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate a unique request ID
			requestID := uuid.New().String()
			
			// Set request ID in context
			ctx := services.SetRequestID(r.Context(), requestID)
			r = r.WithContext(ctx)
			
			// Add request ID to response headers
			w.Header().Set("X-Request-ID", requestID)
			
			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(logger *services.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log incoming request
			logger.InfoWithRequest(r.Context(), "HTTP request received", map[string]interface{}{
				"method": r.Method,
				"path":   r.URL.Path,
				"query":  r.URL.RawQuery,
				"remote": r.RemoteAddr,
				"user_agent": r.UserAgent(),
			})
			
			next.ServeHTTP(w, r)
		})
	}
}