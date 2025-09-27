package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
)

// MetricsMiddleware provides HTTP request instrumentation
func MetricsMiddleware(metricsService *services.MetricsService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Increment in-flight requests
			metricsService.IncrementHTTPInFlight()
			defer metricsService.DecrementHTTPInFlight()
			
			// Create a response writer wrapper to capture status code
			wrapper := &responseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // default status code
			}
			
			// Call the next handler
			next.ServeHTTP(wrapper, r)
			
			// Record metrics
			duration := time.Since(start)
			statusCode := strconv.Itoa(wrapper.statusCode)
			
			metricsService.RecordHTTPRequest(r.Method, r.URL.Path, statusCode, duration)
		})
	}
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader captures the status code
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	if !w.written {
		w.statusCode = statusCode
		w.written = true
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write captures when a response is written (for default 200 status)
func (w *responseWriterWrapper) Write(data []byte) (int, error) {
	if !w.written {
		w.written = true
	}
	return w.ResponseWriter.Write(data)
}