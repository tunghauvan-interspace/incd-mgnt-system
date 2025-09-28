package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/middleware"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userService *services.UserService
	authService *services.AuthService
	logger      *services.Logger
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(userService *services.UserService, authService *services.AuthService, logger *services.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		authService: authService,
		logger:      logger,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode registration request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validateRegisterRequest(&req); err != nil {
		h.logger.Error("Registration validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create user
	user, err := h.userService.Register(r.Context(), &req)
	if err != nil {
		h.logger.Error("User registration failed", map[string]interface{}{
			"username": req.Username,
			"email":    req.Email,
			"error":    err.Error(),
		})

		switch err {
		case services.ErrUsernameExists:
			http.Error(w, "Username already exists", http.StatusConflict)
		case services.ErrEmailExists:
			http.Error(w, "Email already exists", http.StatusConflict)
		default:
			http.Error(w, "Registration failed", http.StatusInternalServerError)
		}
		return
	}

	// Generate auth tokens
	authResponse, err := h.authService.GenerateTokens(user)
	if err != nil {
		h.logger.Error("Failed to generate tokens after registration", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		// Return success but without tokens - user can login separately
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User registered successfully. Please login to get access tokens.",
			"user":    user,
		})
		return
	}

	h.logger.Info("User registered successfully", map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authResponse)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode login request", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validateLoginRequest(&req); err != nil {
		h.logger.Error("Login validation failed", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authenticate user
	authResponse, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		h.logger.Error("User login failed", map[string]interface{}{
			"username": req.Username,
			"email":    req.Email,
			"error":    err.Error(),
		})

		switch err {
		case services.ErrInvalidCredentials:
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		case services.ErrUserNotActive:
			http.Error(w, "Account is not active", http.StatusForbidden)
		default:
			http.Error(w, "Login failed", http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("User logged in successfully", map[string]interface{}{
		"user_id":  authResponse.User.ID,
		"username": authResponse.User.Username,
		"email":    authResponse.User.Email,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	authResponse, err := h.userService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Error("Token refresh failed", map[string]interface{}{
			"error": err.Error(),
		})

		switch err {
		case services.ErrTokenExpired:
			http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		case services.ErrInvalidToken:
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		case services.ErrUserNotActive:
			http.Error(w, "Account is not active", http.StatusForbidden)
		default:
			http.Error(w, "Token refresh failed", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse)
}

// Logout handles user logout (currently just returns success)
// In a more complex implementation, you might maintain a blacklist of tokens
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context for logging
	userID, _ := middleware.GetUserIDFromContext(r.Context())
	
	// Log logout activity
	if userID != "" {
		h.userService.LogUserActivity(
			r.Context(),
			userID,
			"logout",
			"auth",
			"",
			middleware.GetClientIP(r),
			r.UserAgent(),
			nil,
		)

		h.logger.Info("User logged out", map[string]interface{}{
			"user_id": userID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user profile", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile handles profile updates
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validateUpdateProfileRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Failed to update user profile", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})

		switch err {
		case services.ErrEmailExists:
			http.Error(w, "Email already exists", http.StatusConflict)
		default:
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("User profile updated", map[string]interface{}{
		"user_id": userID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ChangePassword handles password changes
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validateChangePasswordRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.userService.ChangePassword(r.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Failed to change password", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})

		switch err {
		case services.ErrInvalidCredentials:
			http.Error(w, "Current password is incorrect", http.StatusBadRequest)
		default:
			http.Error(w, "Failed to change password", http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("User password changed", map[string]interface{}{
		"user_id": userID,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password changed successfully",
	})
}

// Validation helper functions

func (h *AuthHandler) validateRegisterRequest(req *models.RegisterRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return services.ErrInvalidUserID
	}
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return services.ErrInvalidUserID
	}
	if strings.TrimSpace(req.Email) == "" || !strings.Contains(req.Email, "@") {
		return services.ErrEmailExists
	}
	if strings.TrimSpace(req.FullName) == "" {
		return services.ErrInvalidUserID
	}
	if len(req.Password) < 8 {
		return services.ErrInvalidCredentials
	}
	return nil
}

func (h *AuthHandler) validateLoginRequest(req *models.LoginRequest) error {
	if req.Username == "" && req.Email == "" {
		return services.ErrInvalidCredentials
	}
	if req.Password == "" {
		return services.ErrInvalidCredentials
	}
	return nil
}

func (h *AuthHandler) validateUpdateProfileRequest(req *models.UpdateProfileRequest) error {
	if strings.TrimSpace(req.FullName) == "" {
		return services.ErrInvalidUserID
	}
	if strings.TrimSpace(req.Email) == "" || !strings.Contains(req.Email, "@") {
		return services.ErrEmailExists
	}
	return nil
}

func (h *AuthHandler) validateChangePasswordRequest(req *models.ChangePasswordRequest) error {
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return services.ErrInvalidCredentials
	}
	if len(req.NewPassword) < 8 {
		return services.ErrInvalidCredentials
	}
	return nil
}