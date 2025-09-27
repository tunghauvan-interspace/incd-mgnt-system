package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

var (
	ErrUsernameExists    = errors.New("username already exists")
	ErrEmailExists       = errors.New("email already exists")
	ErrInvalidUserID     = errors.New("invalid user ID")
)

// UserService handles user management operations
type UserService struct {
	store       storage.Store
	authService *AuthService
	logger      *Logger
}

// NewUserService creates a new user service
func NewUserService(store storage.Store, authService *AuthService, logger *Logger) *UserService {
	return &UserService{
		store:       store,
		authService: authService,
		logger:      logger,
	}
}

// Register creates a new user account
func (s *UserService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// Check if username exists
	if _, err := s.GetUserByUsername(ctx, req.Username); err == nil {
		return nil, ErrUsernameExists
	} else if !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}

	// Check if email exists
	if _, err := s.GetUserByEmail(ctx, req.Email); err == nil {
		return nil, ErrEmailExists
	} else if !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	// Hash password
	hashedPassword, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Username:  strings.ToLower(strings.TrimSpace(req.Username)),
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		FullName:  strings.TrimSpace(req.FullName),
		Password:  hashedPassword,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign default viewer role to new users
	if err := s.AssignRole(ctx, user.ID, "viewer"); err != nil {
		s.logger.Error("Failed to assign default role to new user", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		// Don't fail registration if role assignment fails
	}

	// Load user with roles for response
	userWithRoles, err := s.GetUserByID(ctx, user.ID)
	if err != nil {
		return user, nil // Return basic user if loading with roles fails
	}

	s.logger.Info("User registered successfully", map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	})

	return userWithRoles, nil
}

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Find user by username or email
	var user *models.User
	var err error

	if req.Email != "" {
		user, err = s.GetUserByEmail(ctx, req.Email)
	} else {
		user, err = s.GetUserByUsername(ctx, req.Username)
	}

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Validate password
	if err := s.authService.ValidatePassword(req.Password, user.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	if err := s.UpdateLastLogin(ctx, user.ID); err != nil {
		s.logger.Error("Failed to update last login", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
	}

	// Generate tokens
	authResponse, err := s.authService.GenerateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Log user activity
	s.LogUserActivity(ctx, user.ID, "login", "auth", "", GetIPAddress(ctx), GetUserAgent(ctx), nil)

	s.logger.Info("User logged in successfully", map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	})

	return authResponse, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	// Parse refresh token to get user ID
	claims, err := s.authService.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Get user with roles
	user, err := s.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Generate new tokens
	return s.authService.RefreshToken(refreshToken, user)
}

// GetUserByID retrieves a user by ID with roles and permissions
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	// This would typically be implemented in a repository pattern
	// For now, we'll need to add this method to the storage interface
	user, err := s.getUserFromStorage(ctx, "id", userID)
	if err != nil {
		return nil, err
	}
	
	// Load roles and permissions
	if err := s.loadUserRoles(ctx, user); err != nil {
		s.logger.Error("Failed to load user roles", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
	}
	
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.getUserFromStorage(ctx, "username", strings.ToLower(username))
	if err != nil {
		return nil, err
	}
	
	// Load roles and permissions
	if err := s.loadUserRoles(ctx, user); err != nil {
		s.logger.Error("Failed to load user roles", map[string]interface{}{
			"username": username,
			"error":    err.Error(),
		})
	}
	
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.getUserFromStorage(ctx, "email", strings.ToLower(email))
	if err != nil {
		return nil, err
	}
	
	// Load roles and permissions
	if err := s.loadUserRoles(ctx, user); err != nil {
		s.logger.Error("Failed to load user roles", map[string]interface{}{
			"email": email,
			"error": err.Error(),
		})
	}
	
	return user, nil
}

// CreateUser creates a new user in the database
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.store.CreateUser(user)
}

// UpdateProfile updates a user's profile information
func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if email is changing and if new email already exists
	if strings.ToLower(req.Email) != strings.ToLower(user.Email) {
		if _, err := s.GetUserByEmail(ctx, req.Email); err == nil {
			return nil, ErrEmailExists
		} else if !errors.Is(err, ErrUserNotFound) {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
	}

	// Update user fields
	user.FullName = strings.TrimSpace(req.FullName)
	user.Email = strings.ToLower(strings.TrimSpace(req.Email))
	user.UpdatedAt = time.Now()

	if err := s.updateUserInStorage(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Log user activity
	s.LogUserActivity(ctx, userID, "update_profile", "user", userID, GetIPAddress(ctx), GetUserAgent(ctx), nil)

	s.logger.Info("User profile updated", map[string]interface{}{
		"user_id": userID,
		"email":   user.Email,
	})

	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, userID string, req *models.ChangePasswordRequest) error {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Validate current password
	if err := s.authService.ValidatePassword(req.CurrentPassword, user.Password); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := s.authService.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	if err := s.updateUserInStorage(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Log user activity
	s.LogUserActivity(ctx, userID, "change_password", "user", userID, GetIPAddress(ctx), GetUserAgent(ctx), nil)

	s.logger.Info("User password changed", map[string]interface{}{
		"user_id": userID,
	})

	return nil
}

// AssignRole assigns a role to a user
func (s *UserService) AssignRole(ctx context.Context, userID, roleName string) error {
	// Get role by name
	role, err := s.store.GetRoleByName(roleName)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	return s.store.AssignRoleToUser(userID, role.ID)
}

// UpdateLastLogin updates the user's last login timestamp
func (s *UserService) UpdateLastLogin(ctx context.Context, userID string) error {
	return s.store.UpdateLastLogin(userID, time.Now())
}

// LogUserActivity logs user activity for audit trails
func (s *UserService) LogUserActivity(ctx context.Context, userID, action, resource, resourceID, ipAddress, userAgent string, metadata map[string]interface{}) {
	activity := &models.UserActivity{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Metadata:   metadata,
		CreatedAt:  time.Now(),
	}

	// Log activity asynchronously to avoid blocking operations
	go func() {
		if err := s.logActivityToStorage(ctx, activity); err != nil {
			s.logger.Error("Failed to log user activity", map[string]interface{}{
				"user_id": userID,
				"action":  action,
				"error":   err.Error(),
			})
		}
	}()
}

// Helper methods

func (s *UserService) getUserFromStorage(ctx context.Context, field, value string) (*models.User, error) {
	switch field {
	case "id":
		return s.store.GetUser(value)
	case "username":
		return s.store.GetUserByUsername(value)
	case "email":
		return s.store.GetUserByEmail(value)
	default:
		return nil, ErrUserNotFound
	}
}

func (s *UserService) loadUserRoles(ctx context.Context, user *models.User) error {
	roles, err := s.store.GetUserRoles(user.ID)
	if err != nil {
		return err
	}
	user.Roles = roles
	return nil
}

func (s *UserService) updateUserInStorage(ctx context.Context, user *models.User) error {
	return s.store.UpdateUser(user)
}

func (s *UserService) logActivityToStorage(ctx context.Context, activity *models.UserActivity) error {
	return s.store.LogUserActivity(activity)
}

// Helper functions for extracting request context

func GetIPAddress(ctx context.Context) string {
	if ip, ok := ctx.Value("ip_address").(string); ok {
		return ip
	}
	return "unknown"
}

func GetUserAgent(ctx context.Context) string {
	if ua, ok := ctx.Value("user_agent").(string); ok {
		return ua
	}
	return "unknown"
}