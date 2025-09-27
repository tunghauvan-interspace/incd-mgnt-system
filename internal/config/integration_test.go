package config

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// Integration test setup for configuration management
func TestIntegration_ConfigurationManagement(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Check if we're running in CI or have database available
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		testDBURL = "postgres://user:password@localhost:5432/incidentdb_test?sslmode=disable"
		
		// Try to connect to see if database is available
		db, err := sql.Open("postgres", testDBURL)
		if err != nil {
			t.Skipf("Database not available for integration tests: %v", err)
		}
		defer db.Close()
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		if err := db.PingContext(ctx); err != nil {
			t.Skipf("Database not responding for integration tests: %v", err)
		}
	}

	t.Run("LoadConfigurationWithDatabase", func(t *testing.T) {
		testLoadConfigurationWithDatabase(t, testDBURL)
	})

	t.Run("ConfigurationValidationIntegration", func(t *testing.T) {
		testConfigurationValidationIntegration(t)
	})

	t.Run("HotReloadIntegration", func(t *testing.T) {
		testHotReloadIntegration(t)
	})

	t.Run("VaultIntegrationPreparation", func(t *testing.T) {
		testVaultIntegrationPreparation(t)
	})
}

func testLoadConfigurationWithDatabase(t *testing.T, dbURL string) {
	// Set environment variables for the test
	oldEnvVars := map[string]string{
		"DATABASE_URL":         os.Getenv("DATABASE_URL"),
		"DB_MAX_OPEN_CONNS":    os.Getenv("DB_MAX_OPEN_CONNS"),
		"DB_MAX_IDLE_CONNS":    os.Getenv("DB_MAX_IDLE_CONNS"),
		"DB_CONN_MAX_LIFETIME": os.Getenv("DB_CONN_MAX_LIFETIME"),
		"SLACK_TOKEN":          os.Getenv("SLACK_TOKEN"),
		"SLACK_CHANNEL":        os.Getenv("SLACK_CHANNEL"),
	}

	// Cleanup function to restore environment
	cleanup := func() {
		for key, value := range oldEnvVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}
	defer cleanup()

	// Set test configuration
	os.Setenv("DATABASE_URL", dbURL)
	os.Setenv("DB_MAX_OPEN_CONNS", "10")
	os.Setenv("DB_MAX_IDLE_CONNS", "2")
	os.Setenv("DB_CONN_MAX_LIFETIME", "1m")
	os.Setenv("SLACK_TOKEN", "xoxb-integration-test")
	os.Setenv("SLACK_CHANNEL", "#integration-test")

	// Load and validate configuration
	cfg, err := LoadAndValidateConfig()
	if err != nil {
		t.Fatalf("Failed to load and validate configuration: %v", err)
	}

	// Verify database configuration
	if cfg.DatabaseURL != dbURL {
		t.Errorf("Expected DatabaseURL %s, got %s", dbURL, cfg.DatabaseURL)
	}
	if cfg.DBMaxOpenConns != 10 {
		t.Errorf("Expected DBMaxOpenConns 10, got %d", cfg.DBMaxOpenConns)
	}
	if cfg.DBMaxIdleConns != 2 {
		t.Errorf("Expected DBMaxIdleConns 2, got %d", cfg.DBMaxIdleConns)
	}
	if cfg.DBConnMaxLifetime != time.Minute {
		t.Errorf("Expected DBConnMaxLifetime 1m, got %v", cfg.DBConnMaxLifetime)
	}

	// Verify notification configuration
	if cfg.SlackToken != "xoxb-integration-test" {
		t.Errorf("Expected SlackToken xoxb-integration-test, got %s", cfg.SlackToken)
	}
	if cfg.SlackChannel != "#integration-test" {
		t.Errorf("Expected SlackChannel #integration-test, got %s", cfg.SlackChannel)
	}

	// Test configuration methods
	if !cfg.HasNotificationConfigured() {
		t.Error("Expected HasNotificationConfigured() to return true")
	}
	if cfg.IsTLSEnabled() {
		t.Error("Expected IsTLSEnabled() to return false")
	}
}

func testConfigurationValidationIntegration(t *testing.T) {
	testCases := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		errorField  string
	}{
		{
			name: "Valid complete configuration",
			envVars: map[string]string{
				"PORT":                "8080",
				"LOG_LEVEL":           "info",
				"DATABASE_URL":        "postgres://user:pass@localhost:5432/db",
				"SLACK_TOKEN":         "xoxb-valid-token",
				"SLACK_CHANNEL":       "#alerts",
				"EMAIL_SMTP_HOST":     "smtp.test.com",
				"EMAIL_USERNAME":      "test@test.com",
				"EMAIL_PASSWORD":      "password",
				"METRICS_ENABLED":     "true",
				"METRICS_PORT":        "9090",
			},
			expectError: false,
		},
		{
			name: "Invalid port configuration",
			envVars: map[string]string{
				"PORT":     "99999",
				"LOG_LEVEL": "info",
			},
			expectError: true,
			errorField:  "PORT",
		},
		{
			name: "Incomplete Slack configuration",
			envVars: map[string]string{
				"PORT":         "8080",
				"LOG_LEVEL":    "info",
				"SLACK_TOKEN":  "xoxb-token",
				// Missing SLACK_CHANNEL
			},
			expectError: true,
			errorField:  "SLACK_CHANNEL",
		},
		{
			name: "Invalid database connection pool",
			envVars: map[string]string{
				"PORT":              "8080",
				"LOG_LEVEL":         "info",
				"DB_MAX_OPEN_CONNS": "5",
				"DB_MAX_IDLE_CONNS": "10", // Invalid: idle > open
			},
			expectError: true,
			errorField:  "DB_MAX_IDLE_CONNS",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save current environment
			savedVars := make(map[string]string)
			for key := range tc.envVars {
				savedVars[key] = os.Getenv(key)
			}

			// Cleanup function
			cleanup := func() {
				for key, value := range savedVars {
					if value == "" {
						os.Unsetenv(key)
					} else {
						os.Setenv(key, value)
					}
				}
			}
			defer cleanup()

			// Clear all config env vars first
			configKeys := []string{
				"PORT", "LOG_LEVEL", "DATABASE_URL", "DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS",
				"SLACK_TOKEN", "SLACK_CHANNEL", "EMAIL_SMTP_HOST", "EMAIL_USERNAME", "EMAIL_PASSWORD",
				"METRICS_ENABLED", "METRICS_PORT",
			}
			for _, key := range configKeys {
				os.Unsetenv(key)
			}

			// Set test environment variables
			for key, value := range tc.envVars {
				os.Setenv(key, value)
			}

			// Load and validate configuration
			_, err := LoadAndValidateConfig()

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected validation error for %s, but got none", tc.name)
					return
				}

				if validationErrs, ok := err.(ValidationErrors); ok {
					found := false
					for _, vErr := range validationErrs {
						if vErr.Field == tc.errorField {
							found = true
							break
						}
					}
					if !found && tc.errorField != "" {
						t.Errorf("Expected validation error for field %s, but it was not found in errors: %v", tc.errorField, err)
					}
				} else {
					t.Errorf("Expected ValidationErrors type, got %T: %v", err, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error for %s, but got: %v", tc.name, err)
				}
			}
		})
	}
}

func testHotReloadIntegration(t *testing.T) {
	// Save current environment
	savedVars := map[string]string{
		"LOG_LEVEL":    os.Getenv("LOG_LEVEL"),
		"DEBUG_MODE":   os.Getenv("DEBUG_MODE"),
		"ENABLE_CORS":  os.Getenv("ENABLE_CORS"),
	}

	cleanup := func() {
		for key, value := range savedVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}
	defer cleanup()

	// Set initial configuration
	os.Setenv("PORT", "8080")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("DEBUG_MODE", "false")
	os.Setenv("ENABLE_CORS", "true")

	// Create watcher with fast polling for test
	watcher := NewWatcher(50 * time.Millisecond)
	err := watcher.Start()
	if err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}
	defer watcher.Stop()

	// Get initial configuration
	initialCfg := watcher.GetConfig()
	if initialCfg.LogLevel != "info" {
		t.Errorf("Expected initial log level 'info', got %s", initialCfg.LogLevel)
	}
	if initialCfg.DebugMode != false {
		t.Errorf("Expected initial debug mode false, got %v", initialCfg.DebugMode)
	}

	// Change reloadable configuration
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("DEBUG_MODE", "true")
	os.Setenv("ENABLE_CORS", "false")

	// Wait for reload
	select {
	case newCfg := <-watcher.ReloadChannel():
		if newCfg.LogLevel != "debug" {
			t.Errorf("Expected reloaded log level 'debug', got %s", newCfg.LogLevel)
		}
		if newCfg.DebugMode != true {
			t.Errorf("Expected reloaded debug mode true, got %v", newCfg.DebugMode)
		}
		if newCfg.EnableCORS != false {
			t.Errorf("Expected reloaded CORS false, got %v", newCfg.EnableCORS)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Configuration reload did not occur within expected time")
	}

	// Verify watcher has the updated configuration
	updatedCfg := watcher.GetConfig()
	if updatedCfg.LogLevel != "debug" {
		t.Errorf("Expected watcher log level 'debug', got %s", updatedCfg.LogLevel)
	}
	if updatedCfg.DebugMode != true {
		t.Errorf("Expected watcher debug mode true, got %v", updatedCfg.DebugMode)
	}

	// Test that sensitive changes don't trigger reload
	initialPort := updatedCfg.Port
	os.Setenv("PORT", "9090")

	// Wait a bit to see if configuration changes (it shouldn't for sensitive changes)
	select {
	case <-watcher.ReloadChannel():
		t.Error("Configuration should not reload for sensitive changes like PORT")
	case <-time.After(200 * time.Millisecond):
		// This is expected - no reload for sensitive changes
	}

	// Verify port didn't change in watcher (sensitive change)
	currentCfg := watcher.GetConfig()
	if currentCfg.Port != initialPort {
		t.Errorf("Port should not change via hot reload, expected %s, got %s", initialPort, currentCfg.Port)
	}
}

func testVaultIntegrationPreparation(t *testing.T) {
	// Save current environment
	savedVars := map[string]string{
		"VAULT_ENABLED":    os.Getenv("VAULT_ENABLED"),
		"VAULT_ADDR":       os.Getenv("VAULT_ADDR"),
		"VAULT_TOKEN":      os.Getenv("VAULT_TOKEN"),
		"VAULT_SECRET_PATH": os.Getenv("VAULT_SECRET_PATH"),
	}

	cleanup := func() {
		for key, value := range savedVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}
	defer cleanup()

	// Test Vault configuration loading
	os.Setenv("VAULT_ENABLED", "false")
	vaultCfg := LoadVaultConfig()
	if vaultCfg.Enabled != false {
		t.Errorf("Expected VAULT_ENABLED false, got %v", vaultCfg.Enabled)
	}

	os.Setenv("VAULT_ENABLED", "true")
	os.Setenv("VAULT_ADDR", "https://vault.test.com:8200")
	os.Setenv("VAULT_TOKEN", "test-token")
	os.Setenv("VAULT_SECRET_PATH", "secret/test-app")

	vaultCfg = LoadVaultConfig()
	if vaultCfg.Enabled != true {
		t.Errorf("Expected VAULT_ENABLED true, got %v", vaultCfg.Enabled)
	}
	if vaultCfg.Address != "https://vault.test.com:8200" {
		t.Errorf("Expected VAULT_ADDR 'https://vault.test.com:8200', got %s", vaultCfg.Address)
	}
	if vaultCfg.Token != "test-token" {
		t.Errorf("Expected VAULT_TOKEN 'test-token', got %s", vaultCfg.Token)
	}
	if vaultCfg.SecretPath != "secret/test-app" {
		t.Errorf("Expected VAULT_SECRET_PATH 'secret/test-app', got %s", vaultCfg.SecretPath)
	}

	// Test SecretManager creation and initialization
	sm := NewSecretManager(vaultCfg)
	if sm == nil {
		t.Fatal("Expected SecretManager to be created")
	}

	// Test disabled vault (should not fail)
	vaultCfg.Enabled = false
	err := sm.Start()
	if err != nil {
		t.Errorf("Expected no error when vault is disabled, got: %v", err)
	}
	sm.Stop()

	// Test secret fallback to environment variables
	os.Setenv("TEST_SECRET", "env-value")
	defer os.Unsetenv("TEST_SECRET")
	
	value := sm.GetSecretOrEnv("non-existent-key", "TEST_SECRET")
	if value != "env-value" {
		t.Errorf("Expected GetSecretOrEnv to fallback to env var 'env-value', got %s", value)
	}

	// Test sensitive field identification
	if !IsSensitiveField("SlackToken") {
		t.Error("Expected SlackToken to be identified as sensitive field")
	}
	if !IsSensitiveField("DatabaseURL") {
		t.Error("Expected DatabaseURL to be identified as sensitive field")
	}
	if IsSensitiveField("Port") {
		t.Error("Expected Port to not be identified as sensitive field")
	}

	// Test sensitive value masking
	masked := MaskSensitiveValue("xoxb-1234567890-abcdefghijk")
	if masked == "xoxb-1234567890-abcdefghijk" {
		t.Error("Expected sensitive value to be masked")
	}
	if len(masked) == 0 {
		t.Error("Expected masked value to not be empty")
	}

	emptyMasked := MaskSensitiveValue("")
	if emptyMasked != "(not set)" {
		t.Errorf("Expected empty value to be marked as '(not set)', got %s", emptyMasked)
	}
}

// TestIntegration_ConfigurationEndToEnd tests the complete configuration lifecycle
func TestIntegration_ConfigurationEndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Save current environment
	envKeys := []string{
		"PORT", "LOG_LEVEL", "DATABASE_URL", "SLACK_TOKEN", "SLACK_CHANNEL",
		"EMAIL_SMTP_HOST", "EMAIL_USERNAME", "EMAIL_PASSWORD", "DEBUG_MODE",
		"METRICS_ENABLED", "VAULT_ENABLED",
	}
	
	savedVars := make(map[string]string)
	for _, key := range envKeys {
		savedVars[key] = os.Getenv(key)
		os.Unsetenv(key) // Clear for clean test
	}

	cleanup := func() {
		for key, value := range savedVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}
	defer cleanup()

	// Test complete configuration lifecycle
	t.Log("Testing complete configuration management lifecycle...")

	// 1. Load minimal configuration
	t.Log("Step 1: Loading minimal configuration...")
	os.Setenv("PORT", "8080")
	os.Setenv("LOG_LEVEL", "info")
	
	cfg1, err := LoadAndValidateConfig()
	if err != nil {
		t.Fatalf("Failed to load minimal configuration: %v", err)
	}
	if cfg1.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", cfg1.Port)
	}

	// 2. Add database configuration
	t.Log("Step 2: Adding database configuration...")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	os.Setenv("DB_MAX_OPEN_CONNS", "20")
	os.Setenv("DB_MAX_IDLE_CONNS", "5")
	
	cfg2, err := LoadAndValidateConfig()
	if err != nil {
		t.Fatalf("Failed to load database configuration: %v", err)
	}
	if cfg2.DatabaseURL != "postgres://test:test@localhost:5432/testdb" {
		t.Errorf("Expected database URL to be set")
	}

	// 3. Add notification configuration
	t.Log("Step 3: Adding notification configuration...")
	os.Setenv("SLACK_TOKEN", "xoxb-test-integration")
	os.Setenv("SLACK_CHANNEL", "#integration-test")
	os.Setenv("EMAIL_SMTP_HOST", "smtp.integration.test")
	os.Setenv("EMAIL_USERNAME", "test@integration.com")
	os.Setenv("EMAIL_PASSWORD", "integration-password")

	cfg3, err := LoadAndValidateConfig()
	if err != nil {
		t.Fatalf("Failed to load notification configuration: %v", err)
	}
	if !cfg3.HasNotificationConfigured() {
		t.Error("Expected notifications to be configured")
	}

	// 4. Enable advanced features
	t.Log("Step 4: Enabling advanced features...")
	os.Setenv("DEBUG_MODE", "true")
	os.Setenv("METRICS_ENABLED", "true")
	os.Setenv("VAULT_ENABLED", "false") // Keep Vault disabled for integration test

	cfg4, err := LoadAndValidateConfig()
	if err != nil {
		t.Fatalf("Failed to load advanced configuration: %v", err)
	}
	if !cfg4.DebugMode {
		t.Error("Expected debug mode to be enabled")
	}
	if !cfg4.MetricsEnabled {
		t.Error("Expected metrics to be enabled")
	}

	// 5. Test configuration methods
	t.Log("Step 5: Testing configuration methods...")
	reloadable := cfg4.GetReloadableConfig()
	if reloadable.LogLevel != "info" {
		t.Errorf("Expected reloadable log level 'info', got %s", reloadable.LogLevel)
	}
	if reloadable.DebugMode != true {
		t.Errorf("Expected reloadable debug mode true, got %v", reloadable.DebugMode)
	}

	// Test safe logging
	t.Log("Step 6: Testing safe configuration logging...")
	// This should not panic and should mask sensitive values
	cfg4.LogConfigSafely()

	t.Log("✅ End-to-end configuration integration test completed successfully")
}