package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

func setupTestAuthHandler(t *testing.T) (*AuthHandler, storage.Store) {
	// Setup test store
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}

	// Create default roles and permissions for testing
	setupTestRolesAndPermissions(t, store)

	// Setup test logger
	logger := services.NewLogger("debug", false)

	// Setup auth service with test config
	authService := services.NewAuthService("test-jwt-secret-32-characters-long!", 1*time.Hour, 24*time.Hour)

	// Setup user service
	userService := services.NewUserService(store, authService, logger)

	// Create auth handler
	authHandler := NewAuthHandler(userService, authService, logger)

	return authHandler, store
}

func setupTestRolesAndPermissions(t *testing.T, store storage.Store) {
	// Create default permissions
	permissions := []*models.Permission{
		{ID: "1", Name: "incidents.read", Resource: "incidents", Action: "read", Description: "View incidents"},
		{ID: "2", Name: "incidents.create", Resource: "incidents", Action: "create", Description: "Create incidents"},
		{ID: "3", Name: "alerts.read", Resource: "alerts", Action: "read", Description: "View alerts"},
	}

	// Note: MemoryStore doesn't have CreatePermission method, so we'll add it manually to the internal data
	// For testing purposes, we'll create roles directly

	// Create default roles
	viewerRole := &models.Role{
		ID:          "viewer-role-id",
		Name:        "viewer",
		DisplayName: "Viewer",
		Description: "Read-only access",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Permissions: []*models.Permission{permissions[0], permissions[2]}, // incidents.read, alerts.read
	}

	adminRole := &models.Role{
		ID:          "admin-role-id",
		Name:        "admin", 
		DisplayName: "Administrator",
		Description: "Full access",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Permissions: permissions, // All permissions
	}

	if err := store.CreateRole(viewerRole); err != nil {
		t.Fatalf("Failed to create viewer role: %v", err)
	}

	if err := store.CreateRole(adminRole); err != nil {
		t.Fatalf("Failed to create admin role: %v", err)
	}
}

func TestAuthHandler_Register(t *testing.T) {
	authHandler, _ := setupTestAuthHandler(t)

	// Test valid registration
	t.Run("Valid Registration", func(t *testing.T) {
		reqBody := models.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			FullName: "Test User",
			Password: "password123",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		authHandler.Register(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		var response models.AuthResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Token == "" {
			t.Error("Expected token in response")
		}

		if response.User.Username != reqBody.Username {
			t.Errorf("Expected username %s, got %s", reqBody.Username, response.User.Username)
		}
	})

	// Test invalid registration (missing fields)
	t.Run("Invalid Registration - Missing Fields", func(t *testing.T) {
		reqBody := models.RegisterRequest{
			Username: "testuser2",
			Email:    "", // Missing email
			FullName: "Test User 2",
			Password: "password123",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		authHandler.Register(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func TestAuthHandler_Login(t *testing.T) {
	authHandler, store := setupTestAuthHandler(t)

	// Create a test user first
	testUser := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		FullName:  "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Hash the password
	authService := services.NewAuthService("test-jwt-secret-32-characters-long!", 1*time.Hour, 24*time.Hour)
	hashedPassword, err := authService.HashPassword("password123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	testUser.Password = hashedPassword

	if err := store.CreateUser(testUser); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test valid login
	t.Run("Valid Login", func(t *testing.T) {
		reqBody := models.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		authHandler.Login(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response models.AuthResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Token == "" {
			t.Error("Expected token in response")
		}

		if response.User.Username != reqBody.Username {
			t.Errorf("Expected username %s, got %s", reqBody.Username, response.User.Username)
		}
	})

	// Test invalid login (wrong password)
	t.Run("Invalid Login - Wrong Password", func(t *testing.T) {
		reqBody := models.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		authHandler.Login(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})

	// Test login with email
	t.Run("Valid Login with Email", func(t *testing.T) {
		reqBody := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		authHandler.Login(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	authHandler, _ := setupTestAuthHandler(t)

	// First, register a user to get tokens
	reqBody := models.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	authHandler.Register(w, req)

	var regResponse models.AuthResponse
	if err := json.NewDecoder(w.Body).Decode(&regResponse); err != nil {
		t.Fatalf("Failed to decode registration response: %v", err)
	}

	// Test token refresh
	t.Run("Valid Token Refresh", func(t *testing.T) {
		refreshReq := struct {
			RefreshToken string `json:"refresh_token"`
		}{
			RefreshToken: regResponse.RefreshToken,
		}

		body, _ := json.Marshal(refreshReq)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		authHandler.RefreshToken(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response models.AuthResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Token == "" {
			t.Error("Expected new token in response")
		}

		if response.Token == regResponse.Token {
			t.Error("Expected new token to be different from original")
		}
	})

	// Test invalid refresh token
	t.Run("Invalid Refresh Token", func(t *testing.T) {
		refreshReq := struct {
			RefreshToken string `json:"refresh_token"`
		}{
			RefreshToken: "invalid_token",
		}

		body, _ := json.Marshal(refreshReq)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		authHandler.RefreshToken(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusUnauthorized, w.Code, w.Body.String())
		}
	})
}

func TestAuthService_TokenValidation(t *testing.T) {
	authService := services.NewAuthService("test-jwt-secret-32-characters-long!", 1*time.Hour, 24*time.Hour)

	// Create a test user with roles
	testUser := &models.User{
		ID:       "test-user-id",
		Username: "testuser",
		Email:    "test@example.com",
		Roles: []*models.Role{
			{
				ID:   "role1",
				Name: "admin",
				Permissions: []*models.Permission{
					{ID: "perm1", Name: "incidents.read", Resource: "incidents", Action: "read"},
					{ID: "perm2", Name: "incidents.create", Resource: "incidents", Action: "create"},
				},
			},
		},
	}

	// Generate tokens
	authResponse, err := authService.GenerateTokens(testUser)
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	// Test token validation
	t.Run("Valid Token Validation", func(t *testing.T) {
		claims, err := authService.ValidateToken(authResponse.Token)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}

		if claims.UserID != testUser.ID {
			t.Errorf("Expected user ID %s, got %s", testUser.ID, claims.UserID)
		}

		if claims.Username != testUser.Username {
			t.Errorf("Expected username %s, got %s", testUser.Username, claims.Username)
		}

		if len(claims.Roles) != 1 || claims.Roles[0] != "admin" {
			t.Errorf("Expected roles [admin], got %v", claims.Roles)
		}

		if len(claims.Permissions) != 2 {
			t.Errorf("Expected 2 permissions, got %d", len(claims.Permissions))
		}
	})

	// Test permission checking
	t.Run("Permission Checking", func(t *testing.T) {
		claims, _ := authService.ValidateToken(authResponse.Token)

		if !authService.HasPermission(claims, "incidents.read") {
			t.Error("Expected user to have incidents.read permission")
		}

		if !authService.HasPermission(claims, "incidents.create") {
			t.Error("Expected user to have incidents.create permission")
		}

		if authService.HasPermission(claims, "incidents.delete") {
			t.Error("Expected user not to have incidents.delete permission")
		}
	})

	// Test role checking
	t.Run("Role Checking", func(t *testing.T) {
		claims, _ := authService.ValidateToken(authResponse.Token)

		if !authService.HasRole(claims, "admin") {
			t.Error("Expected user to have admin role")
		}

		if authService.HasRole(claims, "viewer") {
			t.Error("Expected user not to have viewer role")
		}

		if !authService.HasAnyRole(claims, []string{"admin", "responder"}) {
			t.Error("Expected user to have at least one of admin or responder roles")
		}
	})

	// Test invalid token
	t.Run("Invalid Token Validation", func(t *testing.T) {
		_, err := authService.ValidateToken("invalid_token")
		if err == nil {
			t.Error("Expected error for invalid token")
		}
	})
}

func TestPasswordHashing(t *testing.T) {
	authService := services.NewAuthService("test-jwt-secret-32-characters-long!", 1*time.Hour, 24*time.Hour)

	password := "testpassword123"

	// Test password hashing
	t.Run("Password Hashing", func(t *testing.T) {
		hash, err := authService.HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}

		if hash == password {
			t.Error("Expected hash to be different from original password")
		}

		if len(hash) == 0 {
			t.Error("Expected hash to be non-empty")
		}
	})

	// Test password validation
	t.Run("Password Validation", func(t *testing.T) {
		hash, _ := authService.HashPassword(password)

		// Valid password
		if err := authService.ValidatePassword(password, hash); err != nil {
			t.Errorf("Failed to validate correct password: %v", err)
		}

		// Invalid password
		if err := authService.ValidatePassword("wrongpassword", hash); err == nil {
			t.Error("Expected error for incorrect password")
		}
	})
}