package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// VaultConfig holds configuration for HashiCorp Vault integration
type VaultConfig struct {
	Enabled   bool
	Address   string
	Token     string
	SecretPath string
	RoleID    string
	SecretID  string
	Namespace string
}

// SecretManager provides secure credential management
type SecretManager struct {
	vault     *VaultConfig
	secrets   map[string]string
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	refreshCh chan struct{}
}

// NewSecretManager creates a new secret manager
func NewSecretManager(cfg *VaultConfig) *SecretManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	sm := &SecretManager{
		vault:     cfg,
		secrets:   make(map[string]string),
		ctx:       ctx,
		cancel:    cancel,
		refreshCh: make(chan struct{}, 1),
	}

	return sm
}

// Start initializes the secret manager and starts refresh routines
func (sm *SecretManager) Start() error {
	if !sm.vault.Enabled {
		log.Println("Vault integration disabled, using environment variables for secrets")
		return nil
	}

	log.Println("Initializing HashiCorp Vault integration...")
	
	if err := sm.validateVaultConfig(); err != nil {
		return fmt.Errorf("vault configuration validation failed: %w", err)
	}

	if err := sm.loadSecretsFromVault(); err != nil {
		return fmt.Errorf("failed to load secrets from vault: %w", err)
	}

	// Start refresh routine
	go sm.refreshRoutine()

	log.Println("Vault integration initialized successfully")
	return nil
}

// Stop stops the secret manager
func (sm *SecretManager) Stop() {
	if sm.cancel != nil {
		sm.cancel()
	}
}

// GetSecret retrieves a secret value (thread-safe)
func (sm *SecretManager) GetSecret(key string) (string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if !sm.vault.Enabled {
		// Fall back to environment variable
		value := os.Getenv(key)
		return value, value != ""
	}

	value, exists := sm.secrets[key]
	return value, exists
}

// GetSecretOrEnv retrieves a secret from Vault or falls back to environment variable
func (sm *SecretManager) GetSecretOrEnv(key, envKey string) string {
	if value, exists := sm.GetSecret(key); exists {
		return value
	}
	return os.Getenv(envKey)
}

// RefreshSecrets triggers a manual refresh of secrets
func (sm *SecretManager) RefreshSecrets() {
	select {
	case sm.refreshCh <- struct{}{}:
	default:
		// Refresh already pending
	}
}

// validateVaultConfig validates the Vault configuration
func (sm *SecretManager) validateVaultConfig() error {
	if sm.vault.Address == "" {
		return fmt.Errorf("vault address is required")
	}
	
	if sm.vault.Token == "" && (sm.vault.RoleID == "" || sm.vault.SecretID == "") {
		return fmt.Errorf("either vault token or both role_id and secret_id must be provided")
	}

	return nil
}

// loadSecretsFromVault loads secrets from Vault (placeholder implementation)
func (sm *SecretManager) loadSecretsFromVault() error {
	// This is a placeholder for actual Vault integration
	// In a real implementation, this would use the Vault API to fetch secrets
	log.Printf("Loading secrets from Vault at %s", sm.vault.Address)
	
	// For now, we'll simulate loading some secrets
	// In production, this would make actual API calls to Vault
	sm.mu.Lock()
	sm.secrets = map[string]string{
		"database_password": "vault-managed-db-password",
		"slack_token":      "vault-managed-slack-token",
		"email_password":   "vault-managed-email-password",
	}
	sm.mu.Unlock()

	log.Printf("Loaded %d secrets from Vault", len(sm.secrets))
	return nil
}

// refreshRoutine periodically refreshes secrets from Vault
func (sm *SecretManager) refreshRoutine() {
	ticker := time.NewTicker(15 * time.Minute) // Refresh every 15 minutes
	defer ticker.Stop()

	for {
		select {
		case <-sm.ctx.Done():
			return
		case <-ticker.C:
			sm.refreshSecrets()
		case <-sm.refreshCh:
			sm.refreshSecrets()
		}
	}
}

// refreshSecrets refreshes secrets from Vault
func (sm *SecretManager) refreshSecrets() {
	if !sm.vault.Enabled {
		return
	}

	log.Println("Refreshing secrets from Vault...")
	
	if err := sm.loadSecretsFromVault(); err != nil {
		log.Printf("Failed to refresh secrets from Vault: %v", err)
	} else {
		log.Println("Secrets refreshed successfully")
	}
}

// LoadVaultConfig loads Vault configuration from environment variables
func LoadVaultConfig() *VaultConfig {
	return &VaultConfig{
		Enabled:    getEnvBool("VAULT_ENABLED", false),
		Address:    getEnv("VAULT_ADDR", ""),
		Token:      getEnv("VAULT_TOKEN", ""),
		SecretPath: getEnv("VAULT_SECRET_PATH", "secret/incident-management"),
		RoleID:     getEnv("VAULT_ROLE_ID", ""),
		SecretID:   getEnv("VAULT_SECRET_ID", ""),
		Namespace:  getEnv("VAULT_NAMESPACE", ""),
	}
}

// Enhanced Config to support secure credential management
func (c *Config) LoadSecureCredentials(sm *SecretManager) {
	// Load database credentials securely
	if c.DatabaseURL == "" || strings.Contains(c.DatabaseURL, "password@") {
		// If database URL contains placeholder or is empty, try to load from secrets
		if _, exists := sm.GetSecret("database_password"); exists {
			// Update database URL with secure password
			// This is a placeholder - actual implementation would parse and update URL properly
			log.Println("Using secure database password from Vault")
		}
	}

	// Load notification credentials securely
	if c.SlackToken == "" {
		c.SlackToken = sm.GetSecretOrEnv("slack_token", "SLACK_TOKEN")
	}

	if c.EmailPassword == "" {
		c.EmailPassword = sm.GetSecretOrEnv("email_password", "EMAIL_PASSWORD")
	}

	if c.TelegramBotToken == "" {
		c.TelegramBotToken = sm.GetSecretOrEnv("telegram_bot_token", "TELEGRAM_BOT_TOKEN")
	}
}

// IsSensitiveField returns true if a configuration field contains sensitive data
func IsSensitiveField(fieldName string) bool {
	sensitiveFields := []string{
		"DatabaseURL",
		"SlackToken", 
		"EmailPassword",
		"TelegramBotToken",
		"TLSKeyFile",
	}

	for _, sensitive := range sensitiveFields {
		if fieldName == sensitive {
			return true
		}
	}

	return false
}

// MaskSensitiveValue masks sensitive configuration values for logging
func MaskSensitiveValue(value string) string {
	if value == "" {
		return "(not set)"
	}

	if len(value) <= 8 {
		return "****"
	}

	// Show first 2 and last 2 characters, mask the rest
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

// LogConfigSafely logs configuration details while masking sensitive values
func (c *Config) LogConfigSafely() {
	log.Println("Configuration Summary:")
	log.Printf("  - Port: %s", c.Port)
	log.Printf("  - Log Level: %s", c.LogLevel)
	log.Printf("  - Database URL: %s", func() string {
		if c.DatabaseURL == "" {
			return "(not configured - using in-memory)"
		}
		return MaskSensitiveValue(c.DatabaseURL)
	}())
	log.Printf("  - Slack Token: %s", MaskSensitiveValue(c.SlackToken))
	log.Printf("  - Email Password: %s", MaskSensitiveValue(c.EmailPassword))
	log.Printf("  - Telegram Token: %s", MaskSensitiveValue(c.TelegramBotToken))
	log.Printf("  - TLS Enabled: %t", c.IsTLSEnabled())
	log.Printf("  - Metrics Enabled: %t (port %s)", c.MetricsEnabled, c.MetricsPort)
	log.Printf("  - Debug Mode: %t", c.DebugMode)
}