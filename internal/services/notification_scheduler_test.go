package services

import (
	"testing"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/config"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

func TestNotificationScheduler(t *testing.T) {
	// Setup test environment
	logger := NewLogger("info", true)
	cfg := &config.Config{Port: "8080"}
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}
	
	templateService := NewNotificationTemplateService(logger)
	metricsService := NewMetricsService()
	notificationService := NewNotificationService(cfg, store, templateService, metricsService, logger)
	
	scheduler := NewNotificationScheduler(notificationService, logger)
	defer scheduler.Stop()

	t.Run("ScheduleNotification", func(t *testing.T) {
		channel := &models.NotificationChannel{
			ID:      "scheduler-test-channel",
			Name:    "Scheduler Test Channel",
			Type:    "slack",
			Enabled: true,
			Config:  map[string]string{"token": "test", "channel": "#test"},
		}
		
		incident := &models.Incident{
			ID:          "scheduler-test-incident",
			Title:       "Scheduler Test",
			Description: "Testing notification scheduling",
			Status:      models.IncidentStatusOpen,
			Severity:    models.SeverityHigh,
			CreatedAt:   time.Now(),
		}
		
		scheduledAt := time.Now().Add(5 * time.Second)
		metadata := map[string]interface{}{
			"test": true,
			"priority": "high",
		}
		
		scheduledNotification, err := scheduler.ScheduleNotification(
			incident, 
			channel, 
			"incident_created", 
			scheduledAt, 
			metadata,
		)
		
		if err != nil {
			t.Fatalf("Failed to schedule notification: %v", err)
		}
		
		if scheduledNotification.ID == "" {
			t.Error("Scheduled notification should have an ID")
		}
		
		if scheduledNotification.Status != models.DeliveryStatusPending {
			t.Errorf("Expected status pending, got %s", scheduledNotification.Status)
		}
		
		// Check that it's in the scheduler
		notifications := scheduler.GetScheduledNotifications("", models.DeliveryStatusPending)
		if len(notifications) == 0 {
			t.Error("Scheduled notification should be retrievable")
		}
		
		found := false
		for _, n := range notifications {
			if n.ID == scheduledNotification.ID {
				found = true
				break
			}
		}
		
		if !found {
			t.Error("Scheduled notification not found in list")
		}
	})

	t.Run("CancelScheduledNotification", func(t *testing.T) {
		channel := &models.NotificationChannel{
			ID:   "cancel-test-channel",
			Type: "email",
		}
		
		incident := &models.Incident{
			ID:       "cancel-test-incident",
			Title:    "Cancel Test",
			Severity: models.SeverityLow,
			Status:   models.IncidentStatusOpen,
		}
		
		scheduledAt := time.Now().Add(1 * time.Hour)
		
		scheduledNotification, err := scheduler.ScheduleNotification(
			incident, 
			channel, 
			"incident_resolved", 
			scheduledAt, 
			nil,
		)
		
		if err != nil {
			t.Fatalf("Failed to schedule notification: %v", err)
		}
		
		// Cancel the notification
		err = scheduler.CancelScheduledNotification(scheduledNotification.ID)
		if err != nil {
			t.Errorf("Failed to cancel scheduled notification: %v", err)
		}
		
		// Check that it's no longer in the scheduler
		notifications := scheduler.GetScheduledNotifications("", models.DeliveryStatusPending)
		for _, n := range notifications {
			if n.ID == scheduledNotification.ID {
				t.Error("Cancelled notification should not be in list")
			}
		}
	})

	t.Run("ScheduleRecurringNotification", func(t *testing.T) {
		channel := &models.NotificationChannel{
			ID:   "recurring-test-channel",
			Type: "telegram",
		}
		
		incident := &models.Incident{
			ID:       "recurring-test-incident",
			Title:    "Recurring Test",
			Severity: models.SeverityMedium,
			Status:   models.IncidentStatusOpen,
		}
		
		startTime := time.Now().Add(1 * time.Minute)
		interval := 30 * time.Second
		maxOccurrences := 3
		
		scheduledNotifications, err := scheduler.ScheduleRecurringNotification(
			incident,
			channel,
			"incident_created",
			startTime,
			interval,
			nil, // no end time
			maxOccurrences,
		)
		
		if err != nil {
			t.Fatalf("Failed to schedule recurring notifications: %v", err)
		}
		
		if len(scheduledNotifications) != maxOccurrences {
			t.Errorf("Expected %d scheduled notifications, got %d", maxOccurrences, len(scheduledNotifications))
		}
		
		// Check that they are scheduled at correct intervals
		for i, notification := range scheduledNotifications {
			expectedTime := startTime.Add(time.Duration(i) * interval)
			if !notification.ScheduledAt.Equal(expectedTime) {
				t.Errorf("Notification %d scheduled at wrong time: expected %v, got %v", 
					i, expectedTime, notification.ScheduledAt)
			}
			
			// Check metadata
			if metadata, ok := notification.Metadata["occurrence"].(int); !ok || metadata != i+1 {
				t.Errorf("Notification %d has wrong occurrence metadata", i)
			}
		}
	})

	t.Run("FilterScheduledNotifications", func(t *testing.T) {
		// Create channels
		channel1 := &models.NotificationChannel{ID: "filter-channel-1", Type: "slack"}
		channel2 := &models.NotificationChannel{ID: "filter-channel-2", Type: "email"}
		
		incident := &models.Incident{
			ID:       "filter-test-incident",
			Title:    "Filter Test",
			Severity: models.SeverityHigh,
			Status:   models.IncidentStatusOpen,
		}
		
		scheduledAt := time.Now().Add(2 * time.Hour)
		
		// Schedule notifications on different channels
		_, err := scheduler.ScheduleNotification(incident, channel1, "incident_created", scheduledAt, nil)
		if err != nil {
			t.Fatalf("Failed to schedule notification 1: %v", err)
		}
		
		_, err = scheduler.ScheduleNotification(incident, channel2, "incident_created", scheduledAt.Add(1*time.Minute), nil)
		if err != nil {
			t.Fatalf("Failed to schedule notification 2: %v", err)
		}
		
		// Test channel filtering
		channel1Notifications := scheduler.GetScheduledNotifications("filter-channel-1", "")
		channel1Count := 0
		for _, n := range channel1Notifications {
			if n.Channel.ID == "filter-channel-1" {
				channel1Count++
			}
		}
		
		if channel1Count == 0 {
			t.Error("Should have found notifications for channel 1")
		}
		
		// Test status filtering
		allPending := scheduler.GetScheduledNotifications("", models.DeliveryStatusPending)
		for _, n := range allPending {
			if n.Status != models.DeliveryStatusPending {
				t.Error("Status filter not working correctly")
			}
		}
	})
}