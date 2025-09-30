package services

import (
	"sync"
	"testing"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/config"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

// Shared services to avoid metrics registration conflicts
var (
	testMetricsService *MetricsService
	testLogger *Logger
	testSetupOnce sync.Once
)

func setupTestServices() (*MetricsService, *Logger) {
	testSetupOnce.Do(func() {
		testMetricsService = NewMetricsService()
		testLogger = NewLogger("debug", true)
	})
	return testMetricsService, testLogger
}

func TestNotificationServiceTracking(t *testing.T) {
	// Create memory store
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}

	// Create services
	metricsService, logger := setupTestServices()
	
	templateService := NewNotificationTemplateService(logger)
	
	cfg := &config.Config{
		SlackToken:   "test-token",
		SlackChannel: "#test",
	}
	
	notificationService := NewNotificationService(cfg, store, templateService, metricsService, logger)

	// Create a test incident
	incident := &models.Incident{
		ID:          "test-incident-1",
		Title:       "Test Incident",
		Description: "Test incident description",
		Status:      models.IncidentStatusOpen,
		Severity:    models.SeverityHigh,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create a test notification channel
	channel := &models.NotificationChannel{
		ID:      "test-channel-1",
		Name:    "Test Slack Channel",
		Type:    "slack",
		Enabled: true,
		Config: map[string]string{
			"token":   "test-token",
			"channel": "#test-channel",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store the channel
	err = store.CreateNotificationChannel(channel)
	if err != nil {
		t.Fatalf("Failed to create notification channel: %v", err)
	}

	// Send notification (this will fail due to invalid Slack token, but we're testing tracking)
	err = notificationService.NotifyIncidentCreated(incident)
	// Error is expected due to invalid credentials

	// Check that notification history was created
	histories, err := store.ListNotificationHistory(incident.ID)
	if err != nil {
		t.Fatalf("Failed to list notification history: %v", err)
	}

	if len(histories) == 0 {
		t.Error("Expected at least one notification history entry")
	}

	// Verify history details
	history := histories[0]
	if history.IncidentID != incident.ID {
		t.Errorf("Expected incident ID %s, got %s", incident.ID, history.IncidentID)
	}

	if history.ChannelID != channel.ID {
		t.Errorf("Expected channel ID %s, got %s", channel.ID, history.ChannelID)
	}

	if history.Type != "incident_created" {
		t.Errorf("Expected type 'incident_created', got %s", history.Type)
	}

	if history.Channel != "slack" {
		t.Errorf("Expected channel 'slack', got %s", history.Channel)
	}

	if history.Recipient != "#test-channel" {
		t.Errorf("Expected recipient '#test-channel', got %s", history.Recipient)
	}

	if history.Status != models.DeliveryStatusSent && history.Status != models.DeliveryStatusFailed && history.Status != models.DeliveryStatusRetrying {
		t.Errorf("Expected status to be 'sent', 'failed', or 'retrying', got %s", history.Status)
	}

	if history.Content == "" {
		t.Error("Expected content to be set")
	}

	t.Logf("Notification history created successfully: ID=%s, Status=%s, Recipient=%s", 
		history.ID, history.Status, history.Recipient)
}

func TestNotificationTemplateFromDatabase(t *testing.T) {
	// Create memory store
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}

	// Create custom template
	template := &models.NotificationTemplate{
		ID:        "custom-template-1",
		Name:      "Custom Slack Template",
		Type:      "incident_created",
		Channel:   "slack",
		Subject:   "Custom Alert",
		Body:      "🔥 Custom Incident: {{.Incident.Title}} - {{.Incident.Severity}}",
		IsDefault: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = store.CreateNotificationTemplate(template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Create services
	metricsService, logger := setupTestServices()
	
	templateService := NewNotificationTemplateService(logger)
	
	cfg := &config.Config{
		SlackToken:   "test-token",
		SlackChannel: "#test",
	}
	
	notificationService := NewNotificationService(cfg, store, templateService, metricsService, logger)

	// Create a test notification channel
	channel := &models.NotificationChannel{
		ID:      "test-channel-1",
		Name:    "Test Slack Channel",
		Type:    "slack",
		Enabled: true,
		Config: map[string]string{
			"token":   "test-token",
			"channel": "#test-custom",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test getting template for channel
	retrievedTemplate := notificationService.getTemplateForChannel(channel, "incident_created")
	if retrievedTemplate == nil {
		t.Fatal("Expected to get template from database")
	}

	if retrievedTemplate.ID != template.ID {
		t.Errorf("Expected template ID %s, got %s", template.ID, retrievedTemplate.ID)
	}

	if retrievedTemplate.Body != template.Body {
		t.Errorf("Expected template body %s, got %s", template.Body, retrievedTemplate.Body)
	}

	t.Logf("Database template retrieved successfully: ID=%s, Body=%s", 
		retrievedTemplate.ID, retrievedTemplate.Body)
}