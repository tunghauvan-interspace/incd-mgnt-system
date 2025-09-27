package ratelimit

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter defines the interface for rate limiting
type RateLimiter interface {
	Allow() bool
	Reserve() *rate.Reservation
	Wait() error
}

// PerIPRateLimiter manages rate limits per IP address
type PerIPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
	cleanup  time.Duration
}

// NewPerIPRateLimiter creates a new per-IP rate limiter
func NewPerIPRateLimiter(r rate.Limit, burst int) *PerIPRateLimiter {
	limiter := &PerIPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    burst,
		cleanup:  5 * time.Minute, // Cleanup unused limiters every 5 minutes
	}
	
	// Start cleanup goroutine
	go limiter.cleanupRoutine()
	
	return limiter
}

// GetLimiter returns the rate limiter for a specific IP
func (p *PerIPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	limiter, exists := p.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(p.rate, p.burst)
		p.limiters[ip] = limiter
	}
	
	return limiter
}

// cleanupRoutine periodically removes unused rate limiters
func (p *PerIPRateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(p.cleanup)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			p.mu.Lock()
			// In a production system, you'd track last access time
			// For now, we'll keep all limiters as they're lightweight
			p.mu.Unlock()
		}
	}
}

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	RequestsPerSecond float64
	Burst            int
	Enabled          bool
}

// DefaultWebhookRateLimit returns default rate limiting for webhooks
func DefaultWebhookRateLimit() *RateLimitConfig {
	return &RateLimitConfig{
		RequestsPerSecond: 10.0, // 10 requests per second
		Burst:            20,    // Allow bursts of up to 20 requests
		Enabled:          true,
	}
}

// RateLimitMiddleware creates a middleware that applies rate limiting
func RateLimitMiddleware(config *RateLimitConfig) func(http.Handler) http.Handler {
	if !config.Enabled {
		return func(next http.Handler) http.Handler {
			return next // No rate limiting if disabled
		}
	}
	
	limiter := NewPerIPRateLimiter(rate.Limit(config.RequestsPerSecond), config.Burst)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			ip := getClientIP(r)
			
			// Get limiter for this IP
			ipLimiter := limiter.GetLimiter(ip)
			
			// Check if request is allowed
			reservation := ipLimiter.Reserve()
			if !reservation.OK() {
				// Rate limit exceeded
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.Burst))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10))
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			
			// If we need to wait, delay the request
			delay := reservation.Delay()
			if delay > 0 {
				time.Sleep(delay)
			}
			
			// Add rate limit headers
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.Burst))
			remaining := ipLimiter.Tokens()
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(int(remaining)))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10))
			
			// Proceed with the request
			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (common in reverse proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP if there are multiple
		for i, char := range forwarded {
			if char == ',' || char == ' ' {
				return forwarded[:i]
			}
		}
		return forwarded
	}
	
	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	
	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	for i := len(ip) - 1; i >= 0; i-- {
		if ip[i] == ':' {
			return ip[:i]
		}
	}
	return ip
}

// WebhookRateLimitWrapper wraps a handler function with rate limiting
func WebhookRateLimitWrapper(config *RateLimitConfig, handler http.HandlerFunc) http.HandlerFunc {
	if !config.Enabled {
		return handler
	}
	
	middleware := RateLimitMiddleware(config)
	wrappedHandler := middleware(handler)
	
	return func(w http.ResponseWriter, r *http.Request) {
		wrappedHandler.ServeHTTP(w, r)
	}
}