package config

import (
	"os"
	"testing"
	"time"
)

// Helper function to set environment variables for testing
func setEnvVars(vars map[string]string) {
	for key, value := range vars {
		os.Setenv(key, value)
	}
}

// Helper function to clear environment variables for testing
func clearEnvVars(keys []string) {
	for _, key := range keys {
		os.Unsetenv(key)
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear any existing env vars
	envVars := []string{
		"PORT", "LOG_LEVEL", "ALERTMANAGER_URL", "DATABASE_URL",
		"SLACK_TOKEN", "SLACK_CHANNEL", "EMAIL_SMTP_HOST",
		"TELEGRAM_BOT_TOKEN", "TELEGRAM_CHAT_ID", "METRICS_ENABLED",
	}
	clearEnvVars(envVars)

	cfg := LoadConfig()

	// Test default values
	if cfg.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Port)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("Expected default log level info, got %s", cfg.LogLevel)
	}
	if cfg.AlertmanagerURL != "http://localhost:9093" {
		t.Errorf("Expected default alertmanager URL, got %s", cfg.AlertmanagerURL)
	}
	if cfg.DBMaxOpenConns != 25 {
		t.Errorf("Expected default DB max open conns 25, got %d", cfg.DBMaxOpenConns)
	}
	if cfg.MetricsEnabled != true {
		t.Errorf("Expected metrics enabled by default, got %v", cfg.MetricsEnabled)
	}
}

func TestLoadConfig_CustomValues(t *testing.T) {
	// Set custom environment variables
	setEnvVars(map[string]string{
		"PORT":                "9000",
		"LOG_LEVEL":           "debug",
		"ALERTMANAGER_URL":    "http://custom-alertmanager:9093",
		"DATABASE_URL":        "postgres://user:pass@localhost:5432/testdb",
		"DB_MAX_OPEN_CONNS":   "50",
		"DB_MAX_IDLE_CONNS":   "10",
		"SLACK_TOKEN":         "xoxb-test-token",
		"SLACK_CHANNEL":       "#test-alerts",
		"EMAIL_SMTP_HOST":     "smtp.test.com",
		"EMAIL_SMTP_PORT":     "465",
		"METRICS_ENABLED":     "false",
		"METRICS_PORT":        "8090",
		"DEBUG_MODE":          "true",
	})

	cfg := LoadConfig()

	// Test custom values
	if cfg.Port != "9000" {
		t.Errorf("Expected port 9000, got %s", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected log level debug, got %s", cfg.LogLevel)
	}
	if cfg.AlertmanagerURL != "http://custom-alertmanager:9093" {
		t.Errorf("Expected custom alertmanager URL, got %s", cfg.AlertmanagerURL)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/testdb" {
		t.Errorf("Expected custom database URL, got %s", cfg.DatabaseURL)
	}
	if cfg.DBMaxOpenConns != 50 {
		t.Errorf("Expected DB max open conns 50, got %d", cfg.DBMaxOpenConns)
	}
	if cfg.SlackToken != "xoxb-test-token" {
		t.Errorf("Expected slack token, got %s", cfg.SlackToken)
	}
	if cfg.MetricsEnabled != false {
		t.Errorf("Expected metrics disabled, got %v", cfg.MetricsEnabled)
	}
	if cfg.DebugMode != true {
		t.Errorf("Expected debug mode enabled, got %v", cfg.DebugMode)
	}

	// Cleanup
	clearEnvVars([]string{
		"PORT", "LOG_LEVEL", "ALERTMANAGER_URL", "DATABASE_URL",
		"DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS", "SLACK_TOKEN", "SLACK_CHANNEL",
		"EMAIL_SMTP_HOST", "EMAIL_SMTP_PORT", "METRICS_ENABLED", "METRICS_PORT", "DEBUG_MODE",
	})
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := &Config{
		Port:                "8080",
		LogLevel:            "info",
		MetricsPort:         "9090",
		DBMaxOpenConns:      25,
		DBMaxIdleConns:      5,
		AlertmanagerTimeout: 30,
		EmailSMTPPort:       587,
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	cfg := &Config{
		Port:                "70000", // Invalid port
		LogLevel:            "info",
		MetricsPort:         "9090",
		DBMaxOpenConns:      25,
		DBMaxIdleConns:      5,
		AlertmanagerTimeout: 30,
		EmailSMTPPort:       587,
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for invalid port")
	}

	validationErrs, ok := err.(ValidationErrors)
	if !ok {
		t.Errorf("Expected ValidationErrors, got %T", err)
		return
	}

	found := false
	for _, vErr := range validationErrs {
		if vErr.Field == "PORT" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected validation error for PORT field")
	}
}

func TestValidate_InvalidLogLevel(t *testing.T) {
	cfg := &Config{
		Port:                "8080",
		LogLevel:            "invalid",
		MetricsPort:         "9090",
		DBMaxOpenConns:      25,
		DBMaxIdleConns:      5,
		AlertmanagerTimeout: 30,
		EmailSMTPPort:       587,
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for invalid log level")
	}

	validationErrs, ok := err.(ValidationErrors)
	if !ok {
		t.Errorf("Expected ValidationErrors, got %T", err)
		return
	}

	found := false
	for _, vErr := range validationErrs {
		if vErr.Field == "LOG_LEVEL" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected validation error for LOG_LEVEL field")
	}
}

func TestValidate_DatabaseConnections(t *testing.T) {
	cfg := &Config{
		Port:                "8080",
		LogLevel:            "info",
		MetricsPort:         "9090",
		DBMaxOpenConns:      5,
		DBMaxIdleConns:      10, // Invalid: idle > open
		AlertmanagerTimeout: 30,
		EmailSMTPPort:       587,
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected validation error for idle connections > open connections")
	}
}

func TestValidate_SlackConfig(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		channel     string
		expectError bool
		errorField  string
	}{
		{
			name:        "Both empty (valid - Slack disabled)",
			token:       "",
			channel:     "",
			expectError: false,
		},
		{
			name:        "Valid Slack config",
			token:       "xoxb-test-token",
			channel:     "#alerts",
			expectError: false,
		},
		{
			name:        "Token without channel",
			token:       "xoxb-test-token",
			channel:     "",
			expectError: true,
			errorField:  "SLACK_CHANNEL",
		},
		{
			name:        "Channel without token",
			token:       "",
			channel:     "#alerts",
			expectError: true,
			errorField:  "SLACK_TOKEN",
		},
		{
			name:        "Invalid token format",
			token:       "invalid-token",
			channel:     "#alerts",
			expectError: true,
			errorField:  "SLACK_TOKEN",
		},
		{
			name:        "Invalid channel format",
			token:       "xoxb-test-token",
			channel:     "alerts",
			expectError: true,
			errorField:  "SLACK_CHANNEL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Port:                "8080",
				LogLevel:            "info",
				MetricsPort:         "9090",
				DBMaxOpenConns:      25,
				DBMaxIdleConns:      5,
				AlertmanagerTimeout: 30,
				EmailSMTPPort:       587,
				SlackToken:          tt.token,
				SlackChannel:        tt.channel,
			}

			err := cfg.Validate()
			if tt.expectError {
				if err == nil {
					t.Error("Expected validation error")
					return
				}

				validationErrs, ok := err.(ValidationErrors)
				if !ok {
					t.Errorf("Expected ValidationErrors, got %T", err)
					return
				}

				found := false
				for _, vErr := range validationErrs {
					if vErr.Field == tt.errorField {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected validation error for field %s", tt.errorField)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error, got: %v", err)
				}
			}
		})
	}
}

func TestValidate_TLSConfig(t *testing.T) {
	tests := []struct {
		name        string
		certFile    string
		keyFile     string
		expectError bool
		errorField  string
	}{
		{
			name:        "Both empty (valid - TLS disabled)",
			certFile:    "",
			keyFile:     "",
			expectError: false,
		},
		{
			name:        "Valid TLS config",
			certFile:    "/path/to/cert.pem",
			keyFile:     "/path/to/key.pem",
			expectError: false,
		},
		{
			name:        "Cert without key",
			certFile:    "/path/to/cert.pem",
			keyFile:     "",
			expectError: true,
			errorField:  "TLS_KEY_FILE",
		},
		{
			name:        "Key without cert",
			certFile:    "",
			keyFile:     "/path/to/key.pem",
			expectError: true,
			errorField:  "TLS_CERT_FILE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Port:                "8080",
				LogLevel:            "info",
				MetricsPort:         "9090",
				DBMaxOpenConns:      25,
				DBMaxIdleConns:      5,
				AlertmanagerTimeout: 30,
				EmailSMTPPort:       587,
				TLSCertFile:         tt.certFile,
				TLSKeyFile:          tt.keyFile,
			}

			err := cfg.Validate()
			if tt.expectError {
				if err == nil {
					t.Error("Expected validation error")
					return
				}

				validationErrs, ok := err.(ValidationErrors)
				if !ok {
					t.Errorf("Expected ValidationErrors, got %T", err)
					return
				}

				found := false
				for _, vErr := range validationErrs {
					if vErr.Field == tt.errorField {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected validation error for field %s", tt.errorField)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error, got: %v", err)
				}
			}
		})
	}
}

func TestHasNotificationConfigured(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected bool
	}{
		{
			name: "No notifications configured",
			config: &Config{
				SlackToken:       "",
				SlackChannel:     "",
				EmailSMTPHost:    "",
				TelegramBotToken: "",
			},
			expected: false,
		},
		{
			name: "Slack configured",
			config: &Config{
				SlackToken:   "xoxb-token",
				SlackChannel: "#alerts",
			},
			expected: true,
		},
		{
			name: "Email configured",
			config: &Config{
				EmailSMTPHost: "smtp.test.com",
				EmailUsername: "test@test.com",
				EmailPassword: "password",
			},
			expected: true,
		},
		{
			name: "Telegram configured",
			config: &Config{
				TelegramBotToken: "123:token",
				TelegramChatID:   "456",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.HasNotificationConfigured()
			if result != tt.expected {
				t.Errorf("Expected HasNotificationConfigured() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsTLSEnabled(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected bool
	}{
		{
			name: "TLS disabled",
			config: &Config{
				TLSCertFile: "",
				TLSKeyFile:  "",
			},
			expected: false,
		},
		{
			name: "TLS enabled",
			config: &Config{
				TLSCertFile: "/path/to/cert.pem",
				TLSKeyFile:  "/path/to/key.pem",
			},
			expected: true,
		},
		{
			name: "Only cert file (incomplete)",
			config: &Config{
				TLSCertFile: "/path/to/cert.pem",
				TLSKeyFile:  "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsTLSEnabled()
			if result != tt.expected {
				t.Errorf("Expected IsTLSEnabled() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLoadAndValidateConfig(t *testing.T) {
	// Set valid environment variables
	setEnvVars(map[string]string{
		"PORT":                "8080",
		"LOG_LEVEL":           "info",
		"METRICS_PORT":        "9090",
		"DB_MAX_OPEN_CONNS":   "25",
		"DB_MAX_IDLE_CONNS":   "5",
		"ALERTMANAGER_TIMEOUT": "30",
		"EMAIL_SMTP_PORT":     "587",
	})

	cfg, err := LoadAndValidateConfig()
	if err != nil {
		t.Errorf("Expected valid config to load and validate successfully, got error: %v", err)
	}
	if cfg == nil {
		t.Error("Expected config to be returned")
	}

	// Cleanup
	clearEnvVars([]string{
		"PORT", "LOG_LEVEL", "METRICS_PORT", "DB_MAX_OPEN_CONNS",
		"DB_MAX_IDLE_CONNS", "ALERTMANAGER_TIMEOUT", "EMAIL_SMTP_PORT",
	})
}

func TestDurationParsing(t *testing.T) {
	setEnvVars(map[string]string{
		"SERVER_READ_TIMEOUT":  "45s",
		"SERVER_WRITE_TIMEOUT": "1m30s", 
		"MAX_INCIDENT_AGE":     "48h",
	})

	cfg := LoadConfig()

	if cfg.ServerReadTimeout != 45*time.Second {
		t.Errorf("Expected ServerReadTimeout 45s, got %v", cfg.ServerReadTimeout)
	}
	if cfg.ServerWriteTimeout != 90*time.Second {
		t.Errorf("Expected ServerWriteTimeout 90s, got %v", cfg.ServerWriteTimeout)
	}
	if cfg.MaxIncidentAge != 48*time.Hour {
		t.Errorf("Expected MaxIncidentAge 48h, got %v", cfg.MaxIncidentAge)
	}

	// Cleanup
	clearEnvVars([]string{"SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "MAX_INCIDENT_AGE"})
}