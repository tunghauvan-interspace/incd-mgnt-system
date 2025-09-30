package storage

import (
	"testing"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

func TestNotificationTemplateOperations(t *testing.T) {
	store, err := NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}

	// Test template creation
	template := &models.NotificationTemplate{
		ID:        "template-1",
		Name:      "Test Template",
		Type:      "incident_created",
		Channel:   "slack",
		Body:      "Test message: {{.Title}}",
		IsDefault: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = store.CreateNotificationTemplate(template)
	if err != nil {
		t.Fatalf("Failed to create notification template: %v", err)
	}

	// Test template retrieval
	retrieved, err := store.GetNotificationTemplate("template-1")
	if err != nil {
		t.Fatalf("Failed to get notification template: %v", err)
	}

	if retrieved.Name != template.Name {
		t.Errorf("Expected name %s, got %s", template.Name, retrieved.Name)
	}

	// Test get by type
	byType, err := store.GetNotificationTemplateByType("incident_created", "slack")
	if err != nil {
		t.Fatalf("Failed to get template by type: %v", err)
	}

	if byType.ID != template.ID {
		t.Errorf("Expected ID %s, got %s", template.ID, byType.ID)
	}

	// Test list templates
	templates, err := store.ListNotificationTemplates()
	if err != nil {
		t.Fatalf("Failed to list templates: %v", err)
	}

	if len(templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(templates))
	}

	// Test template update
	template.Name = "Updated Template"
	template.UpdatedAt = time.Now()
	
	err = store.UpdateNotificationTemplate(template)
	if err != nil {
		t.Fatalf("Failed to update notification template: %v", err)
	}

	updated, err := store.GetNotificationTemplate("template-1")
	if err != nil {
		t.Fatalf("Failed to get updated template: %v", err)
	}

	if updated.Name != "Updated Template" {
		t.Errorf("Expected updated name, got %s", updated.Name)
	}

	// Test template deletion
	err = store.DeleteNotificationTemplate("template-1")
	if err != nil {
		t.Fatalf("Failed to delete notification template: %v", err)
	}

	_, err = store.GetNotificationTemplate("template-1")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestNotificationHistoryOperations(t *testing.T) {
	store, err := NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}

	// Test history creation
	history := &models.NotificationHistory{
		ID:         "history-1",
		IncidentID: "incident-1",
		ChannelID:  "channel-1",
		Type:       "incident_created",
		Channel:    "slack",
		Recipient:  "user@example.com",
		Content:    "Test notification content",
		Status:     models.DeliveryStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = store.CreateNotificationHistory(history)
	if err != nil {
		t.Fatalf("Failed to create notification history: %v", err)
	}

	// Test history retrieval
	retrieved, err := store.GetNotificationHistory("history-1")
	if err != nil {
		t.Fatalf("Failed to get notification history: %v", err)
	}

	if retrieved.IncidentID != history.IncidentID {
		t.Errorf("Expected incident ID %s, got %s", history.IncidentID, retrieved.IncidentID)
	}

	// Test list history for incident
	histories, err := store.ListNotificationHistory("incident-1")
	if err != nil {
		t.Fatalf("Failed to list notification history: %v", err)
	}

	if len(histories) != 1 {
		t.Errorf("Expected 1 history entry, got %d", len(histories))
	}

	// Test list all history
	allHistories, err := store.ListNotificationHistory("")
	if err != nil {
		t.Fatalf("Failed to list all notification history: %v", err)
	}

	if len(allHistories) != 1 {
		t.Errorf("Expected 1 history entry, got %d", len(allHistories))
	}

	// Test history update
	history.Status = models.DeliveryStatusSent
	sentTime := time.Now()
	history.SentAt = &sentTime
	history.UpdatedAt = time.Now()

	err = store.UpdateNotificationHistory(history)
	if err != nil {
		t.Fatalf("Failed to update notification history: %v", err)
	}

	updated, err := store.GetNotificationHistory("history-1")
	if err != nil {
		t.Fatalf("Failed to get updated history: %v", err)
	}

	if updated.Status != models.DeliveryStatusSent {
		t.Errorf("Expected status %s, got %s", models.DeliveryStatusSent, updated.Status)
	}

	if updated.SentAt == nil {
		t.Error("Expected SentAt to be set")
	}

	// Test history deletion
	err = store.DeleteNotificationHistory("history-1")
	if err != nil {
		t.Fatalf("Failed to delete notification history: %v", err)
	}

	_, err = store.GetNotificationHistory("history-1")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestNotificationBatchOperations(t *testing.T) {
	store, err := NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}

	// Test batch creation
	batch := &models.NotificationBatch{
		ID:            "batch-1",
		ChannelID:     "channel-1",
		Type:          "incident_created",
		Count:         2,
		Status:        models.DeliveryStatusPending,
		Notifications: []string{"history-1", "history-2"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = store.CreateNotificationBatch(batch)
	if err != nil {
		t.Fatalf("Failed to create notification batch: %v", err)
	}

	// Test batch retrieval
	retrieved, err := store.GetNotificationBatch("batch-1")
	if err != nil {
		t.Fatalf("Failed to get notification batch: %v", err)
	}

	if retrieved.Count != batch.Count {
		t.Errorf("Expected count %d, got %d", batch.Count, retrieved.Count)
	}

	if len(retrieved.Notifications) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(retrieved.Notifications))
	}

	// Test list batches for channel
	batches, err := store.ListNotificationBatches("channel-1")
	if err != nil {
		t.Fatalf("Failed to list notification batches: %v", err)
	}

	if len(batches) != 1 {
		t.Errorf("Expected 1 batch, got %d", len(batches))
	}

	// Test list all batches
	allBatches, err := store.ListNotificationBatches("")
	if err != nil {
		t.Fatalf("Failed to list all notification batches: %v", err)
	}

	if len(allBatches) != 1 {
		t.Errorf("Expected 1 batch, got %d", len(allBatches))
	}

	// Test batch update
	batch.Status = models.DeliveryStatusSent
	processedTime := time.Now()
	batch.ProcessedAt = &processedTime
	batch.UpdatedAt = time.Now()

	err = store.UpdateNotificationBatch(batch)
	if err != nil {
		t.Fatalf("Failed to update notification batch: %v", err)
	}

	updated, err := store.GetNotificationBatch("batch-1")
	if err != nil {
		t.Fatalf("Failed to get updated batch: %v", err)
	}

	if updated.Status != models.DeliveryStatusSent {
		t.Errorf("Expected status %s, got %s", models.DeliveryStatusSent, updated.Status)
	}

	if updated.ProcessedAt == nil {
		t.Error("Expected ProcessedAt to be set")
	}

	// Test batch deletion
	err = store.DeleteNotificationBatch("batch-1")
	if err != nil {
		t.Fatalf("Failed to delete notification batch: %v", err)
	}

	_, err = store.GetNotificationBatch("batch-1")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}