package config

import (
	"log"
	"sync"
	"time"
)

// Watcher provides hot-reloading capability for configuration
type Watcher struct {
	config       *Config
	mu           sync.RWMutex
	stopCh       chan struct{}
	reloadCh     chan *Config
	interval     time.Duration
	lastModTime  time.Time
}

// NewWatcher creates a new configuration watcher
func NewWatcher(interval time.Duration) *Watcher {
	return &Watcher{
		stopCh:   make(chan struct{}),
		reloadCh: make(chan *Config, 1),
		interval: interval,
	}
}

// Start begins watching for configuration changes
func (w *Watcher) Start() error {
	// Load initial configuration
	cfg, err := LoadAndValidateConfig()
	if err != nil {
		return err
	}

	w.mu.Lock()
	w.config = cfg
	w.mu.Unlock()

	// Start watching in a goroutine
	go w.watch()

	return nil
}

// Stop stops the configuration watcher
func (w *Watcher) Stop() {
	close(w.stopCh)
}

// GetConfig returns the current configuration (thread-safe)
func (w *Watcher) GetConfig() *Config {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.config
}

// ReloadChannel returns a channel that receives new configuration on reload
func (w *Watcher) ReloadChannel() <-chan *Config {
	return w.reloadCh
}

// watch monitors for configuration changes
func (w *Watcher) watch() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.checkForChanges()
		}
	}
}

// checkForChanges checks if configuration has changed and reloads if needed
func (w *Watcher) checkForChanges() {
	// Load new configuration
	newCfg, err := LoadAndValidateConfig()
	if err != nil {
		log.Printf("Configuration reload failed: %v", err)
		return
	}

	w.mu.Lock()
	oldCfg := w.config

	// Check if non-sensitive configuration has changed
	if w.hasNonSensitiveChanges(oldCfg, newCfg) {
		log.Println("Configuration changes detected, reloading...")
		w.config = newCfg
		w.mu.Unlock()

		// Send new config to reload channel (non-blocking)
		select {
		case w.reloadCh <- newCfg:
		default:
			// Channel full, skip this reload notification
		}

		log.Println("Configuration reloaded successfully")
	} else {
		w.mu.Unlock()
	}
}

// hasNonSensitiveChanges checks if non-sensitive configuration values have changed
func (w *Watcher) hasNonSensitiveChanges(old, new *Config) bool {
	// Only reload for safe, non-sensitive configuration changes
	// Sensitive changes (like database connections) require restart
	return old.LogLevel != new.LogLevel ||
		old.AlertmanagerTimeout != new.AlertmanagerTimeout ||
		old.MetricsEnabled != new.MetricsEnabled ||
		old.WebhookTimeout != new.WebhookTimeout ||
		old.NotificationTimeout != new.NotificationTimeout ||
		old.MaxIncidentAge != new.MaxIncidentAge ||
		old.EnableCORS != new.EnableCORS ||
		old.CORSOrigin != new.CORSOrigin ||
		old.DebugMode != new.DebugMode
}

// ReloadableConfig represents configuration that can be safely reloaded
type ReloadableConfig struct {
	LogLevel            string
	AlertmanagerTimeout int
	MetricsEnabled      bool
	WebhookTimeout      time.Duration
	NotificationTimeout time.Duration
	MaxIncidentAge      time.Duration
	EnableCORS          bool
	CORSOrigin          string
	DebugMode           bool
}

// GetReloadableConfig extracts only the safely reloadable configuration
func (c *Config) GetReloadableConfig() ReloadableConfig {
	return ReloadableConfig{
		LogLevel:            c.LogLevel,
		AlertmanagerTimeout: c.AlertmanagerTimeout,
		MetricsEnabled:      c.MetricsEnabled,
		WebhookTimeout:      c.WebhookTimeout,
		NotificationTimeout: c.NotificationTimeout,
		MaxIncidentAge:      c.MaxIncidentAge,
		EnableCORS:          c.EnableCORS,
		CORSOrigin:          c.CORSOrigin,
		DebugMode:           c.DebugMode,
	}
}

// ApplyReloadableConfig applies reloadable configuration changes
func (c *Config) ApplyReloadableConfig(reloadable ReloadableConfig) {
	c.LogLevel = reloadable.LogLevel
	c.AlertmanagerTimeout = reloadable.AlertmanagerTimeout
	c.MetricsEnabled = reloadable.MetricsEnabled
	c.WebhookTimeout = reloadable.WebhookTimeout
	c.NotificationTimeout = reloadable.NotificationTimeout
	c.MaxIncidentAge = reloadable.MaxIncidentAge
	c.EnableCORS = reloadable.EnableCORS
	c.CORSOrigin = reloadable.CORSOrigin
	c.DebugMode = reloadable.DebugMode
}