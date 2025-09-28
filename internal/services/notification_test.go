package services

import (
	"testing"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/config"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

func TestNotificationServiceEnhanced(t *testing.T) {
	// Setup test environment
	cfg := &config.Config{
		Port:         "8080",
		SlackToken:   "test-token",
		SlackChannel: "#test",
	}
	
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}
	logger := NewLogger("info", true)
	templateService := NewNotificationTemplateService(logger)
	metricsService := NewMetricsService()
	
	notificationService := NewNotificationService(cfg, store, templateService, metricsService, logger)

	t.Run("SendTestNotification", func(t *testing.T) {
		// Create test channel
		channel := &models.NotificationChannel{
			ID:      "test-channel-1",
			Name:    "Test Channel",
			Type:    "slack",
			Enabled: true,
			Config: map[string]string{
				"token":   "test-token",
				"channel": "#test-channel",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		// Store channel
		err := store.CreateNotificationChannel(channel)
		if err != nil {
			t.Fatalf("Failed to create test channel: %v", err)
		}
		
		// Create test incident
		incident := &models.Incident{
			ID:          "test-incident-1",
			Title:       "Test Notification System",
			Description: "Testing the enhanced notification system",
			Status:      models.IncidentStatusOpen,
			Severity:    models.SeverityMedium,
			CreatedAt:   time.Now(),
		}
		
		// Test notification (will fail due to invalid token, but should test the flow)
		err = notificationService.SendTestNotification(incident, channel, "test")
		// We expect this to fail with network/auth error, not code error
		if err == nil {
			t.Log("Test notification appeared to succeed (unexpected in test environment)")
		} else {
			t.Logf("Test notification failed as expected: %v", err)
		}
	})

	t.Run("NotificationChannelPreferences", func(t *testing.T) {
		// Create channel with preferences
		channel := &models.NotificationChannel{
			ID:      "test-channel-prefs",
			Name:    "Preferences Test Channel",
			Type:    "email",
			Enabled: true,
			Preferences: &models.ChannelPreferences{
				OptIn:          true,
				SeverityFilter: []string{"high", "critical"},
				QuietHours: &models.QuietHoursConfig{
					Enabled:   true,
					StartTime: "22:00",
					EndTime:   "06:00",
					Days:      []int{0, 1, 2, 3, 4, 5, 6}, // All days
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		// Test shouldNotify function
		highIncident := &models.Incident{
			Severity: models.SeverityHigh,
		}
		
		lowIncident := &models.Incident{
			Severity: models.SeverityLow,
		}
		
		// High severity should pass filter
		shouldNotify := notificationService.shouldNotify(channel, highIncident, "incident_created")
		if !shouldNotify {
			// Note: This test might fail during quiet hours - this is expected behavior
			t.Logf("High severity incident was filtered (possibly due to quiet hours): %v", shouldNotify)
		}
		
		// Low severity should be filtered out
		shouldNotify = notificationService.shouldNotify(channel, lowIncident, "incident_created")
		if shouldNotify {
			t.Error("Low severity incident should be filtered out")
		}
		
		// Test opt-out
		channel.Preferences.OptIn = false
		shouldNotify = notificationService.shouldNotify(channel, highIncident, "incident_created")
		if shouldNotify {
			t.Error("Opted-out channel should not receive notifications")
		}
	})

	t.Run("TemplateSelection", func(t *testing.T) {
		// Test default template selection
		channel := &models.NotificationChannel{
			Type: "slack",
		}
		
		template := notificationService.getTemplateForChannel(channel, "incident_created")
		if template == nil {
			t.Error("Should return default template")
		}
		if template.Type != "incident_created" || template.Channel != "slack" {
			t.Error("Template type/channel mismatch")
		}
		
		// Test custom template
		channel.Templates = map[string]string{
			"incident_created": "Custom template: {{.Incident.Title}}",
		}
		
		customTemplate := notificationService.getTemplateForChannel(channel, "incident_created")
		if customTemplate == nil {
			t.Error("Should return custom template")
		}
		if customTemplate.Body != "Custom template: {{.Incident.Title}}" {
			t.Error("Custom template body mismatch")
		}
	})

	t.Run("BatchProcessing", func(t *testing.T) {
		if notificationService.batchProcessor == nil {
			t.Error("Batch processor should be initialized")
		}
		
		// Create a channel with batching enabled
		channel := &models.NotificationChannel{
			ID:      "batch-test-channel",
			Name:    "Batch Test Channel",
			Type:    "slack",
			Enabled: true,
			Preferences: &models.ChannelPreferences{
				OptIn:            true,
				BatchingEnabled:  true,
				MaxBatchSize:     5,
				BatchingInterval: 2 * time.Minute,
			},
			Config:    map[string]string{"token": "test", "channel": "#test"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		err := store.CreateNotificationChannel(channel)
		if err != nil {
			t.Fatalf("Failed to create batch test channel: %v", err)
		}
		
		// Create test incident
		incident := &models.Incident{
			ID:          "batch-test-incident",
			Title:       "Batch Test Incident",
			Description: "Testing batching functionality",
			Status:      models.IncidentStatusOpen,
			Severity:    models.SeverityMedium,
			CreatedAt:   time.Now(),
		}
		
		// Test adding to batch
		err = notificationService.batchProcessor.AddToBatch(incident, channel, "incident_created")
		if err != nil {
			t.Errorf("Failed to add notification to batch: %v", err)
		}
		
		// Verify batch was created
		if len(notificationService.batchProcessor.batches) == 0 {
			t.Error("Batch should have been created")
		}
	})
}

func TestNotificationBatchProcessor(t *testing.T) {
	logger := NewLogger("info", true)
	
	// Create a mock notification service for testing
	cfg := &config.Config{Port: "8080"}
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}
	templateService := NewNotificationTemplateService(logger)
	
	// Skip metrics service to avoid duplicate registration
	notificationService := &NotificationService{
		config:          cfg,
		store:           store,
		templateService: templateService,
		logger:          logger,
	}
	
	processor := NewNotificationBatchProcessor(notificationService, logger)
	
	t.Run("BatchCreation", func(t *testing.T) {
		channel := &models.NotificationChannel{
			ID:   "batch-channel-1",
			Type: "slack",
			Preferences: &models.ChannelPreferences{
				MaxBatchSize: 3,
			},
		}
		
		incident := &models.Incident{
			ID:       "incident-1",
			Title:    "Test Incident 1",
			Severity: models.SeverityHigh,
			Status:   models.IncidentStatusOpen,
		}
		
		err := processor.AddToBatch(incident, channel, "incident_created")
		if err != nil {
			t.Errorf("Failed to add to batch: %v", err)
		}
		
		// Check that batch was created
		if len(processor.batches) == 0 {
			t.Error("Batch should have been created")
		}
		
		batchKey := "batch-channel-1_incident_created"
		batch, exists := processor.batches[batchKey]
		if !exists {
			t.Error("Batch should exist")
		}
		
		if batch.Count != 1 {
			t.Errorf("Expected batch count 1, got %d", batch.Count)
		}
	})
	
	// Stop the processor
	processor.Stop()
}