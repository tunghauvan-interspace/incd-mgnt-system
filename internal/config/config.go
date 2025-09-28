package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Core application settings
	Port                string
	LogLevel            string
	AlertmanagerURL     string
	AlertmanagerTimeout int

	// Database settings
	DatabaseURL         string
	DBMaxOpenConns      int
	DBMaxIdleConns      int
	DBConnMaxLifetime   time.Duration

	// Notification settings
	SlackToken          string
	SlackChannel        string
	EmailSMTPHost       string
	EmailSMTPPort       int
	EmailUsername       string
	EmailPassword       string
	EmailFrom           string
	EmailTo             string
	TelegramBotToken    string
	TelegramChatID      string

	// Metrics settings
	MetricsEnabled      bool
	MetricsPort         string

	// Security settings
	ServerReadTimeout   time.Duration
	ServerWriteTimeout  time.Duration
	ServerIdleTimeout   time.Duration
	TLSCertFile         string
	TLSKeyFile          string

	// JWT Authentication settings
	JWTSecret           string
	JWTExpiration       time.Duration
	RefreshExpiration   time.Duration

	// Advanced settings
	WebhookTimeout      time.Duration
	NotificationTimeout time.Duration
	MaxIncidentAge      time.Duration
	EnableCORS          bool
	CORSOrigin          string

	// Development settings
	DebugMode           bool
	TestDatabaseURL     string
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("config validation failed for %s: %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no validation errors"
	}
	
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("configuration validation failed:\n  - %s", strings.Join(messages, "\n  - "))
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		// Core application settings
		Port:                getEnv("PORT", "8080"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		AlertmanagerURL:     getEnv("ALERTMANAGER_URL", "http://localhost:9093"),
		AlertmanagerTimeout: getEnvInt("ALERTMANAGER_TIMEOUT", 30),

		// Database settings
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		DBMaxOpenConns:      getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:      getEnvInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime:   getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),

		// Notification settings
		SlackToken:          getEnv("SLACK_TOKEN", ""),
		SlackChannel:        getEnv("SLACK_CHANNEL", ""),
		EmailSMTPHost:       getEnv("EMAIL_SMTP_HOST", ""),
		EmailSMTPPort:       getEnvInt("EMAIL_SMTP_PORT", 587),
		EmailUsername:       getEnv("EMAIL_USERNAME", ""),
		EmailPassword:       getEnv("EMAIL_PASSWORD", ""),
		EmailFrom:           getEnv("EMAIL_FROM", ""),
		EmailTo:             getEnv("EMAIL_TO", ""),
		TelegramBotToken:    getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:      getEnv("TELEGRAM_CHAT_ID", ""),

		// Metrics settings
		MetricsEnabled:      getEnvBool("METRICS_ENABLED", true),
		MetricsPort:         getEnv("METRICS_PORT", "9090"),

		// Security settings
		ServerReadTimeout:   getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
		ServerWriteTimeout:  getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		ServerIdleTimeout:   getEnvDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
		TLSCertFile:         getEnv("TLS_CERT_FILE", ""),
		TLSKeyFile:          getEnv("TLS_KEY_FILE", ""),

		// JWT Authentication settings
		JWTSecret:           getEnv("JWT_SECRET", generateDefaultJWTSecret()),
		JWTExpiration:       getEnvDuration("JWT_EXPIRATION", 1*time.Hour),
		RefreshExpiration:   getEnvDuration("REFRESH_EXPIRATION", 24*time.Hour),

		// Advanced settings
		WebhookTimeout:      getEnvDuration("WEBHOOK_TIMEOUT", 30*time.Second),
		NotificationTimeout: getEnvDuration("NOTIFICATION_TIMEOUT", 15*time.Second),
		MaxIncidentAge:      getEnvDuration("MAX_INCIDENT_AGE", 24*time.Hour),
		EnableCORS:          getEnvBool("ENABLE_CORS", true),
		CORSOrigin:          getEnv("CORS_ORIGIN", "*"),

		// Development settings
		DebugMode:           getEnvBool("DEBUG_MODE", false),
		TestDatabaseURL:     getEnv("TEST_DATABASE_URL", ""),
	}

	return cfg
}

// LoadAndValidateConfig loads and validates configuration
func LoadAndValidateConfig() (*Config, error) {
	cfg := LoadConfig()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate validates the configuration and returns any errors
func (c *Config) Validate() error {
	var errors ValidationErrors

	// Validate core settings
	if err := c.validatePort(c.Port, "PORT"); err != nil {
		errors = append(errors, *err)
	}
	if err := c.validateLogLevel(c.LogLevel); err != nil {
		errors = append(errors, *err)
	}
	if err := c.validatePort(c.MetricsPort, "METRICS_PORT"); err != nil {
		errors = append(errors, *err)
	}

	// Validate database settings
	if c.DBMaxOpenConns <= 0 {
		errors = append(errors, ValidationError{
			Field:   "DB_MAX_OPEN_CONNS",
			Message: "must be greater than 0",
		})
	}
	if c.DBMaxIdleConns < 0 {
		errors = append(errors, ValidationError{
			Field:   "DB_MAX_IDLE_CONNS", 
			Message: "must be greater than or equal to 0",
		})
	}
	if c.DBMaxIdleConns > c.DBMaxOpenConns {
		errors = append(errors, ValidationError{
			Field:   "DB_MAX_IDLE_CONNS",
			Message: "cannot be greater than DB_MAX_OPEN_CONNS",
		})
	}

	// Validate timeouts
	if c.AlertmanagerTimeout <= 0 {
		errors = append(errors, ValidationError{
			Field:   "ALERTMANAGER_TIMEOUT",
			Message: "must be greater than 0",
		})
	}

	// Validate notification settings
	if err := c.validateSlackConfig(); err != nil {
		errors = append(errors, *err)
	}
	if err := c.validateEmailConfig(); err != nil {
		errors = append(errors, *err)
	}
	if err := c.validateTelegramConfig(); err != nil {
		errors = append(errors, *err)
	}

	// Validate TLS settings
	if err := c.validateTLSConfig(); err != nil {
		errors = append(errors, *err)
	}

	// Validate JWT settings
	if err := c.validateJWTConfig(); err != nil {
		errors = append(errors, *err)
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// validatePort validates that a port string is a valid port number
func (c *Config) validatePort(port, fieldName string) *ValidationError {
	if port == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "cannot be empty",
		}
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be a valid number",
		}
	}

	if portNum < 1 || portNum > 65535 {
		return &ValidationError{
			Field:   fieldName,
			Message: "must be between 1 and 65535",
		}
	}

	return nil
}

// validateLogLevel validates the log level
func (c *Config) validateLogLevel(level string) *ValidationError {
	validLevels := []string{"debug", "info", "warn", "error"}
	for _, valid := range validLevels {
		if strings.ToLower(level) == valid {
			return nil
		}
	}

	return &ValidationError{
		Field:   "LOG_LEVEL",
		Message: fmt.Sprintf("must be one of: %s", strings.Join(validLevels, ", ")),
	}
}

// validateSlackConfig validates Slack configuration
func (c *Config) validateSlackConfig() *ValidationError {
	if c.SlackToken == "" && c.SlackChannel == "" {
		return nil // Both empty is OK (Slack disabled)
	}

	if c.SlackToken != "" && c.SlackChannel == "" {
		return &ValidationError{
			Field:   "SLACK_CHANNEL",
			Message: "required when SLACK_TOKEN is provided",
		}
	}

	if c.SlackToken == "" && c.SlackChannel != "" {
		return &ValidationError{
			Field:   "SLACK_TOKEN",
			Message: "required when SLACK_CHANNEL is provided",
		}
	}

	if c.SlackToken != "" && !strings.HasPrefix(c.SlackToken, "xoxb-") {
		return &ValidationError{
			Field:   "SLACK_TOKEN",
			Message: "must be a bot token starting with 'xoxb-'",
		}
	}

	if c.SlackChannel != "" && !strings.HasPrefix(c.SlackChannel, "#") {
		return &ValidationError{
			Field:   "SLACK_CHANNEL",
			Message: "must start with '#' (e.g., #alerts)",
		}
	}

	return nil
}

// validateEmailConfig validates email configuration
func (c *Config) validateEmailConfig() *ValidationError {
	if c.EmailSMTPHost == "" && c.EmailUsername == "" && c.EmailPassword == "" {
		return nil // All empty is OK (Email disabled)
	}

	if c.EmailSMTPHost != "" && c.EmailUsername == "" {
		return &ValidationError{
			Field:   "EMAIL_USERNAME",
			Message: "required when EMAIL_SMTP_HOST is provided",
		}
	}

	if c.EmailSMTPHost != "" && c.EmailPassword == "" {
		return &ValidationError{
			Field:   "EMAIL_PASSWORD",
			Message: "required when EMAIL_SMTP_HOST is provided",
		}
	}

	if c.EmailSMTPPort <= 0 || c.EmailSMTPPort > 65535 {
		return &ValidationError{
			Field:   "EMAIL_SMTP_PORT",
			Message: "must be between 1 and 65535",
		}
	}

	return nil
}

// validateTelegramConfig validates Telegram configuration
func (c *Config) validateTelegramConfig() *ValidationError {
	if c.TelegramBotToken == "" && c.TelegramChatID == "" {
		return nil // Both empty is OK (Telegram disabled)
	}

	if c.TelegramBotToken != "" && c.TelegramChatID == "" {
		return &ValidationError{
			Field:   "TELEGRAM_CHAT_ID",
			Message: "required when TELEGRAM_BOT_TOKEN is provided",
		}
	}

	if c.TelegramBotToken == "" && c.TelegramChatID != "" {
		return &ValidationError{
			Field:   "TELEGRAM_BOT_TOKEN",
			Message: "required when TELEGRAM_CHAT_ID is provided",
		}
	}

	return nil
}

// validateTLSConfig validates TLS configuration
func (c *Config) validateTLSConfig() *ValidationError {
	if c.TLSCertFile == "" && c.TLSKeyFile == "" {
		return nil // Both empty is OK (TLS disabled)
	}

	if c.TLSCertFile != "" && c.TLSKeyFile == "" {
		return &ValidationError{
			Field:   "TLS_KEY_FILE",
			Message: "required when TLS_CERT_FILE is provided",
		}
	}

	if c.TLSCertFile == "" && c.TLSKeyFile != "" {
		return &ValidationError{
			Field:   "TLS_CERT_FILE",
			Message: "required when TLS_KEY_FILE is provided",
		}
	}

	return nil
}

// validateJWTConfig validates JWT configuration
func (c *Config) validateJWTConfig() *ValidationError {
	if c.JWTSecret == "" {
		return &ValidationError{
			Field:   "JWT_SECRET",
			Message: "cannot be empty",
		}
	}

	if len(c.JWTSecret) < 32 {
		return &ValidationError{
			Field:   "JWT_SECRET",
			Message: "must be at least 32 characters long for security",
		}
	}

	if c.JWTExpiration <= 0 {
		return &ValidationError{
			Field:   "JWT_EXPIRATION",
			Message: "must be greater than 0",
		}
	}

	if c.RefreshExpiration <= 0 {
		return &ValidationError{
			Field:   "REFRESH_EXPIRATION",
			Message: "must be greater than 0",
		}
	}

	if c.JWTExpiration >= c.RefreshExpiration {
		return &ValidationError{
			Field:   "JWT_EXPIRATION",
			Message: "should be less than REFRESH_EXPIRATION",
		}
	}

	return nil
}

// HasNotificationConfigured returns true if at least one notification method is configured
func (c *Config) HasNotificationConfigured() bool {
	return (c.SlackToken != "" && c.SlackChannel != "") ||
		(c.EmailSMTPHost != "" && c.EmailUsername != "" && c.EmailPassword != "") ||
		(c.TelegramBotToken != "" && c.TelegramChatID != "")
}

// IsTLSEnabled returns true if TLS is configured
func (c *Config) IsTLSEnabled() bool {
	return c.TLSCertFile != "" && c.TLSKeyFile != ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// generateDefaultJWTSecret generates a default JWT secret if none is provided
func generateDefaultJWTSecret() string {
	return "default-jwt-secret-please-change-in-production"
}