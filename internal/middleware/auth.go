package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
)

type contextKey string

const (
	// ClaimsContextKey is the key used to store JWT claims in the request context
	ClaimsContextKey contextKey = "claims"
	// UserIDContextKey is the key used to store user ID in the request context
	UserIDContextKey contextKey = "user_id"
)

// AuthMiddleware provides JWT authentication middleware
func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check Bearer prefix
			tokenParts := strings.SplitN(authHeader, " ", 2)
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := tokenParts[1]

			// Validate token
			claims, err := authService.ValidateToken(token)
			if err != nil {
				var statusCode int
				switch err {
				case services.ErrTokenExpired:
					statusCode = http.StatusUnauthorized
				case services.ErrInvalidToken:
					statusCode = http.StatusUnauthorized
				default:
					statusCode = http.StatusInternalServerError
				}
				http.Error(w, err.Error(), statusCode)
				return
			}

			// Add claims and user ID to request context
			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)
			
			// Add IP address and user agent for audit logging
			ctx = context.WithValue(ctx, "ip_address", GetClientIP(r))
			ctx = context.WithValue(ctx, "user_agent", r.UserAgent())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuthMiddleware provides optional JWT authentication middleware
// If token is present, it validates and adds claims to context
// If token is missing, the request proceeds without authentication
func OptionalAuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			
			ctx := r.Context()
			// Add IP address and user agent for audit logging
			ctx = context.WithValue(ctx, "ip_address", GetClientIP(r))
			ctx = context.WithValue(ctx, "user_agent", r.UserAgent())

			if authHeader != "" {
				tokenParts := strings.SplitN(authHeader, " ", 2)
				if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
					token := tokenParts[1]
					if claims, err := authService.ValidateToken(token); err == nil {
						ctx = context.WithValue(ctx, ClaimsContextKey, claims)
						ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)
					}
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole creates middleware that requires specific roles
func RequireRole(authService *services.AuthService, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsContextKey).(*services.Claims)
			if !ok || claims == nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			if !authService.HasAnyRole(claims, roles) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission creates middleware that requires specific permissions
func RequirePermission(authService *services.AuthService, permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsContextKey).(*services.Claims)
			if !ok || claims == nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			if !authService.HasAnyPermission(claims, permissions) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetClientIP extracts the real client IP from various headers
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to remote address
	ip := r.RemoteAddr
	// Remove port if present
	if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
		ip = ip[:colonIndex]
	}

	return ip
}

// GetClaimsFromContext extracts claims from request context
func GetClaimsFromContext(ctx context.Context) (*services.Claims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*services.Claims)
	return claims, ok
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	return userID, ok
}