package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Port                string
	DatabaseURL         string
	LogLevel            string
	
	// Database connection pooling settings
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
	TelegramBotToken    string
	TelegramChatID      string

	// Alertmanager settings
	AlertmanagerURL     string
	AlertmanagerTimeout int

	// Metrics settings
	MetricsEnabled      bool
	MetricsPort         string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		Port:                getEnv("PORT", "8080"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		
		// Database connection pooling settings
		DBMaxOpenConns:      getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:      getEnvInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime:   getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		
		SlackToken:          getEnv("SLACK_TOKEN", ""),
		SlackChannel:        getEnv("SLACK_CHANNEL", ""),
		EmailSMTPHost:       getEnv("EMAIL_SMTP_HOST", ""),
		EmailSMTPPort:       getEnvInt("EMAIL_SMTP_PORT", 587),
		EmailUsername:       getEnv("EMAIL_USERNAME", ""),
		EmailPassword:       getEnv("EMAIL_PASSWORD", ""),
		TelegramBotToken:    getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:      getEnv("TELEGRAM_CHAT_ID", ""),
		
		AlertmanagerURL:     getEnv("ALERTMANAGER_URL", "http://localhost:9093"),
		AlertmanagerTimeout: getEnvInt("ALERTMANAGER_TIMEOUT", 30),
		
		MetricsEnabled:      getEnvBool("METRICS_ENABLED", true),
		MetricsPort:         getEnv("METRICS_PORT", "9090"),
	}

	return cfg
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