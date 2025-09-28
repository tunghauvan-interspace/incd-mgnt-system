package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrTokenExpired       = errors.New("token has expired")
	ErrInvalidToken       = errors.New("invalid token")
)

// Claims represents JWT claims with user information
type Claims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// AuthService handles authentication operations
type AuthService struct {
	jwtSecret      []byte
	jwtExpiration  time.Duration
	refreshExpiration time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(jwtSecret string, jwtExpiration, refreshExpiration time.Duration) *AuthService {
	return &AuthService{
		jwtSecret:         []byte(jwtSecret),
		jwtExpiration:     jwtExpiration,
		refreshExpiration: refreshExpiration,
	}
}

// HashPassword hashes a password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// ValidatePassword validates a password against its hash
func (s *AuthService) ValidatePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateTokens generates both access and refresh tokens for a user
func (s *AuthService) GenerateTokens(user *models.User) (*models.AuthResponse, error) {
	// Extract role names and permissions
	roles := make([]string, len(user.Roles))
	permissionMap := make(map[string]bool)
	
	for i, role := range user.Roles {
		roles[i] = role.Name
		for _, permission := range role.Permissions {
			permissionMap[permission.Name] = true
		}
	}
	
	// Convert permissions map to slice
	permissions := make([]string, 0, len(permissionMap))
	for perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	// Create access token
	now := time.Now()
	expiresAt := now.Add(s.jwtExpiration)
	
	claims := &Claims{
		UserID:      user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "incident-management-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Create refresh token
	refreshClaims := &jwt.RegisteredClaims{
		ID:        uuid.New().String(),
		Subject:   user.ID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiration)),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "incident-management-system",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &models.AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshTokenString,
		User:         *user,
		ExpiresAt:    expiresAt,
	}, nil
}

// ValidateToken validates and parses a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshToken generates a new access token from a valid refresh token
func (s *AuthService) RefreshToken(refreshTokenString string, user *models.User) (*models.AuthResponse, error) {
	token, err := jwt.ParseWithClaims(refreshTokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.Subject != user.ID {
			return nil, ErrInvalidToken
		}
		return s.GenerateTokens(user)
	}

	return nil, ErrInvalidToken
}

// GeneratePasswordResetToken generates a secure token for password reset
func (s *AuthService) GeneratePasswordResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HasPermission checks if a user has a specific permission
func (s *AuthService) HasPermission(claims *Claims, permission string) bool {
	for _, perm := range claims.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// HasRole checks if a user has a specific role
func (s *AuthService) HasRole(claims *Claims, role string) bool {
	for _, r := range claims.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole checks if a user has any of the specified roles
func (s *AuthService) HasAnyRole(claims *Claims, roles []string) bool {
	for _, requiredRole := range roles {
		if s.HasRole(claims, requiredRole) {
			return true
		}
	}
	return false
}

// HasAnyPermission checks if a user has any of the specified permissions
func (s *AuthService) HasAnyPermission(claims *Claims, permissions []string) bool {
	for _, requiredPerm := range permissions {
		if s.HasPermission(claims, requiredPerm) {
			return true
		}
	}
	return false
}