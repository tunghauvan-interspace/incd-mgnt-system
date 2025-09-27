package config

import (
	"os"
	"testing"
	"time"
)

func TestWatcher_Start(t *testing.T) {
	// Set valid environment variables
	setEnvVars(map[string]string{
		"PORT":      "8080",
		"LOG_LEVEL": "info",
	})
	defer clearEnvVars([]string{"PORT", "LOG_LEVEL"})

	watcher := NewWatcher(100 * time.Millisecond)

	err := watcher.Start()
	if err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}
	defer watcher.Stop()

	cfg := watcher.GetConfig()
	if cfg == nil {
		t.Error("Expected configuration to be loaded")
	}
	if cfg.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", cfg.Port)
	}
}

func TestWatcher_GetConfig(t *testing.T) {
	// Set valid environment variables
	setEnvVars(map[string]string{
		"PORT":      "8080",
		"LOG_LEVEL": "info",
	})
	defer clearEnvVars([]string{"PORT", "LOG_LEVEL"})

	watcher := NewWatcher(100 * time.Millisecond)
	err := watcher.Start()
	if err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}
	defer watcher.Stop()

	cfg1 := watcher.GetConfig()
	cfg2 := watcher.GetConfig()

	// Should return same configuration instance
	if cfg1 != cfg2 {
		t.Error("Expected same configuration instance")
	}
}

func TestWatcher_HotReload(t *testing.T) {
	// Set initial environment variables
	setEnvVars(map[string]string{
		"PORT":      "8080",
		"LOG_LEVEL": "info",
		"DEBUG_MODE": "false",
	})
	defer clearEnvVars([]string{"PORT", "LOG_LEVEL", "DEBUG_MODE"})

	watcher := NewWatcher(50 * time.Millisecond)
	err := watcher.Start()
	if err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}
	defer watcher.Stop()

	// Get initial configuration
	initialCfg := watcher.GetConfig()
	if initialCfg.LogLevel != "info" {
		t.Errorf("Expected initial log level info, got %s", initialCfg.LogLevel)
	}
	if initialCfg.DebugMode != false {
		t.Errorf("Expected initial debug mode false, got %v", initialCfg.DebugMode)
	}

	// Change configuration
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("DEBUG_MODE", "true")

	// Wait for reload
	select {
	case newCfg := <-watcher.ReloadChannel():
		if newCfg.LogLevel != "debug" {
			t.Errorf("Expected reloaded log level debug, got %s", newCfg.LogLevel)
		}
		if newCfg.DebugMode != true {
			t.Errorf("Expected reloaded debug mode true, got %v", newCfg.DebugMode)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("Configuration reload did not occur within expected time")
	}

	// Verify watcher has updated configuration
	updatedCfg := watcher.GetConfig()
	if updatedCfg.LogLevel != "debug" {
		t.Errorf("Expected watcher log level debug, got %s", updatedCfg.LogLevel)
	}
	if updatedCfg.DebugMode != true {
		t.Errorf("Expected watcher debug mode true, got %v", updatedCfg.DebugMode)
	}
}

func TestWatcher_InvalidConfigurationChange(t *testing.T) {
	// Set valid initial environment variables
	setEnvVars(map[string]string{
		"PORT":      "8080",
		"LOG_LEVEL": "info",
	})
	defer clearEnvVars([]string{"PORT", "LOG_LEVEL"})

	watcher := NewWatcher(50 * time.Millisecond)
	err := watcher.Start()
	if err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}
	defer watcher.Stop()

	// Get initial configuration
	initialCfg := watcher.GetConfig()
	if initialCfg.LogLevel != "info" {
		t.Errorf("Expected initial log level info, got %s", initialCfg.LogLevel)
	}

	// Change to invalid configuration
	os.Setenv("LOG_LEVEL", "invalid")

	// Wait to see if configuration changes (it shouldn't)
	select {
	case <-watcher.ReloadChannel():
		t.Error("Configuration should not have reloaded with invalid settings")
	case <-time.After(200 * time.Millisecond):
		// This is expected - no reload should occur
	}

	// Verify watcher still has original configuration
	currentCfg := watcher.GetConfig()
	if currentCfg.LogLevel != "info" {
		t.Errorf("Expected configuration to remain unchanged, got log level: %s", currentCfg.LogLevel)
	}
}

func TestWatcher_NoChangeNoReload(t *testing.T) {
	// Set environment variables
	setEnvVars(map[string]string{
		"PORT":      "8080",
		"LOG_LEVEL": "info",
	})
	defer clearEnvVars([]string{"PORT", "LOG_LEVEL"})

	watcher := NewWatcher(50 * time.Millisecond)
	err := watcher.Start()
	if err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}
	defer watcher.Stop()

	// Wait to see if configuration changes (it shouldn't)
	select {
	case <-watcher.ReloadChannel():
		t.Error("Configuration should not reload when no changes occur")
	case <-time.After(200 * time.Millisecond):
		// This is expected - no reload should occur
	}
}

func TestWatcher_SensitiveChangeNoReload(t *testing.T) {
	// Set initial environment variables
	setEnvVars(map[string]string{
		"PORT":        "8080",
		"LOG_LEVEL":   "info",
		"DATABASE_URL": "postgres://old-url",
	})
	defer clearEnvVars([]string{"PORT", "LOG_LEVEL", "DATABASE_URL"})

	watcher := NewWatcher(50 * time.Millisecond)
	err := watcher.Start()
	if err != nil {
		t.Fatalf("Failed to start watcher: %v", err)
	}
	defer watcher.Stop()

	// Change sensitive configuration (database URL)
	os.Setenv("DATABASE_URL", "postgres://new-url")

	// Wait to see if configuration reloads (it shouldn't for sensitive changes)
	select {
	case <-watcher.ReloadChannel():
		t.Error("Configuration should not reload for sensitive changes like DATABASE_URL")
	case <-time.After(200 * time.Millisecond):
		// This is expected - no reload should occur for sensitive changes
	}
}

func TestHasNonSensitiveChanges(t *testing.T) {
	watcher := NewWatcher(time.Second)

	tests := []struct {
		name     string
		old      *Config
		new      *Config
		expected bool
	}{
		{
			name: "No changes",
			old: &Config{
				LogLevel:   "info",
				DebugMode:  false,
				EnableCORS: true,
			},
			new: &Config{
				LogLevel:   "info",
				DebugMode:  false,
				EnableCORS: true,
			},
			expected: false,
		},
		{
			name: "Log level change",
			old: &Config{
				LogLevel:   "info",
				DebugMode:  false,
				EnableCORS: true,
			},
			new: &Config{
				LogLevel:   "debug",
				DebugMode:  false,
				EnableCORS: true,
			},
			expected: true,
		},
		{
			name: "Debug mode change",
			old: &Config{
				LogLevel:   "info",
				DebugMode:  false,
				EnableCORS: true,
			},
			new: &Config{
				LogLevel:   "info",
				DebugMode:  true,
				EnableCORS: true,
			},
			expected: true,
		},
		{
			name: "CORS change",
			old: &Config{
				LogLevel:   "info",
				DebugMode:  false,
				EnableCORS: true,
			},
			new: &Config{
				LogLevel:   "info",
				DebugMode:  false,
				EnableCORS: false,
			},
			expected: true,
		},
		{
			name: "Sensitive change only (should return false)",
			old: &Config{
				LogLevel:    "info",
				DatabaseURL: "postgres://old",
				Port:        "8080",
			},
			new: &Config{
				LogLevel:    "info",
				DatabaseURL: "postgres://new",
				Port:        "9090",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := watcher.hasNonSensitiveChanges(tt.old, tt.new)
			if result != tt.expected {
				t.Errorf("Expected hasNonSensitiveChanges() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetReloadableConfig(t *testing.T) {
	cfg := &Config{
		LogLevel:            "debug",
		Port:                "8080", // This should not be included
		AlertmanagerTimeout: 45,
		DatabaseURL:         "postgres://secret", // This should not be included
		MetricsEnabled:      false,
		WebhookTimeout:      30 * time.Second,
		EnableCORS:          false,
		DebugMode:           true,
	}

	reloadable := cfg.GetReloadableConfig()

	if reloadable.LogLevel != "debug" {
		t.Errorf("Expected LogLevel debug, got %s", reloadable.LogLevel)
	}
	if reloadable.AlertmanagerTimeout != 45 {
		t.Errorf("Expected AlertmanagerTimeout 45, got %d", reloadable.AlertmanagerTimeout)
	}
	if reloadable.MetricsEnabled != false {
		t.Errorf("Expected MetricsEnabled false, got %v", reloadable.MetricsEnabled)
	}
	if reloadable.WebhookTimeout != 30*time.Second {
		t.Errorf("Expected WebhookTimeout 30s, got %v", reloadable.WebhookTimeout)
	}
	if reloadable.EnableCORS != false {
		t.Errorf("Expected EnableCORS false, got %v", reloadable.EnableCORS)
	}
	if reloadable.DebugMode != true {
		t.Errorf("Expected DebugMode true, got %v", reloadable.DebugMode)
	}
}

func TestApplyReloadableConfig(t *testing.T) {
	cfg := &Config{
		LogLevel:            "info",
		Port:                "8080", // This should remain unchanged
		AlertmanagerTimeout: 30,
		DatabaseURL:         "postgres://db", // This should remain unchanged
		MetricsEnabled:      true,
		WebhookTimeout:      15 * time.Second,
		EnableCORS:          true,
		DebugMode:           false,
	}

	reloadable := ReloadableConfig{
		LogLevel:            "debug",
		AlertmanagerTimeout: 60,
		MetricsEnabled:      false,
		WebhookTimeout:      45 * time.Second,
		EnableCORS:          false,
		DebugMode:           true,
	}

	cfg.ApplyReloadableConfig(reloadable)

	// Check reloadable fields were updated
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel debug, got %s", cfg.LogLevel)
	}
	if cfg.AlertmanagerTimeout != 60 {
		t.Errorf("Expected AlertmanagerTimeout 60, got %d", cfg.AlertmanagerTimeout)
	}
	if cfg.MetricsEnabled != false {
		t.Errorf("Expected MetricsEnabled false, got %v", cfg.MetricsEnabled)
	}
	if cfg.EnableCORS != false {
		t.Errorf("Expected EnableCORS false, got %v", cfg.EnableCORS)
	}

	// Check non-reloadable fields were not changed
	if cfg.Port != "8080" {
		t.Errorf("Expected Port to remain 8080, got %s", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://db" {
		t.Errorf("Expected DatabaseURL to remain unchanged, got %s", cfg.DatabaseURL)
	}
}