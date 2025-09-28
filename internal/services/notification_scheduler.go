package services

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

// NotificationScheduler handles scheduling and time-based delivery of notifications
type NotificationScheduler struct {
	service           *NotificationService
	logger            *Logger
	scheduledNotifications map[string]*ScheduledNotification
	mutex             sync.RWMutex
	ticker            *time.Ticker
	stopChan          chan bool
}

// ScheduledNotification represents a notification scheduled for future delivery
type ScheduledNotification struct {
	ID           string                 `json:"id"`
	Incident     *models.Incident       `json:"incident"`
	Channel      *models.NotificationChannel `json:"channel"`
	Type         string                 `json:"type"`
	ScheduledAt  time.Time              `json:"scheduled_at"`
	CreatedAt    time.Time              `json:"created_at"`
	Status       models.NotificationDeliveryStatus `json:"status"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// NewNotificationScheduler creates a new notification scheduler
func NewNotificationScheduler(service *NotificationService, logger *Logger) *NotificationScheduler {
	scheduler := &NotificationScheduler{
		service:                service,
		logger:                 logger,
		scheduledNotifications: make(map[string]*ScheduledNotification),
		stopChan:               make(chan bool),
	}
	
	// Start the scheduler ticker
	scheduler.ticker = time.NewTicker(30 * time.Second) // Check every 30 seconds
	go scheduler.processScheduledNotifications()
	
	return scheduler
}

// ScheduleNotification schedules a notification for future delivery
func (s *NotificationScheduler) ScheduleNotification(
	incident *models.Incident,
	channel *models.NotificationChannel,
	notificationType string,
	scheduledAt time.Time,
	metadata map[string]interface{},
) (*ScheduledNotification, error) {
	
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	notification := &ScheduledNotification{
		ID:          uuid.New().String(),
		Incident:    incident,
		Channel:     channel,
		Type:        notificationType,
		ScheduledAt: scheduledAt,
		CreatedAt:   time.Now(),
		Status:      models.DeliveryStatusPending,
		Metadata:    metadata,
	}
	
	s.scheduledNotifications[notification.ID] = notification
	
	s.logger.Info("Notification scheduled", map[string]interface{}{
		"notification_id": notification.ID,
		"channel_id":      channel.ID,
		"incident_id":     incident.ID,
		"type":            notificationType,
		"scheduled_at":    scheduledAt.Format(time.RFC3339),
	})
	
	return notification, nil
}

// CancelScheduledNotification cancels a scheduled notification
func (s *NotificationScheduler) CancelScheduledNotification(notificationID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if notification, exists := s.scheduledNotifications[notificationID]; exists {
		delete(s.scheduledNotifications, notificationID)
		
		s.logger.Info("Scheduled notification cancelled", map[string]interface{}{
			"notification_id": notificationID,
			"channel_id":      notification.Channel.ID,
		})
		
		return nil
	}
	
	return ErrNotFound
}

// GetScheduledNotifications returns all scheduled notifications, optionally filtered
func (s *NotificationScheduler) GetScheduledNotifications(channelID string, status models.NotificationDeliveryStatus) []*ScheduledNotification {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var notifications []*ScheduledNotification
	
	for _, notification := range s.scheduledNotifications {
		// Apply filters
		if channelID != "" && notification.Channel.ID != channelID {
			continue
		}
		
		if status != "" && notification.Status != status {
			continue
		}
		
		notifications = append(notifications, notification)
	}
	
	// Sort by scheduled time
	sort.Slice(notifications, func(i, j int) bool {
		return notifications[i].ScheduledAt.Before(notifications[j].ScheduledAt)
	})
	
	return notifications
}

// processScheduledNotifications processes notifications that are due for delivery
func (s *NotificationScheduler) processScheduledNotifications() {
	for {
		select {
		case <-s.ticker.C:
			s.processDueNotifications()
		case <-s.stopChan:
			return
		}
	}
}

// processDueNotifications processes notifications that are due for delivery
func (s *NotificationScheduler) processDueNotifications() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	now := time.Now()
	var toProcess []*ScheduledNotification
	
	// Find notifications that are due
	for id, notification := range s.scheduledNotifications {
		if notification.Status == models.DeliveryStatusPending && notification.ScheduledAt.Before(now) {
			toProcess = append(toProcess, notification)
			delete(s.scheduledNotifications, id)
		}
	}
	
	// Process due notifications (outside of the main mutex to avoid blocking)
	if len(toProcess) > 0 {
		s.mutex.Unlock()
		s.processNotifications(toProcess)
		s.mutex.Lock()
	}
}

// processNotifications processes a batch of due notifications
func (s *NotificationScheduler) processNotifications(notifications []*ScheduledNotification) {
	for _, notification := range notifications {
		s.logger.Info("Processing scheduled notification", map[string]interface{}{
			"notification_id": notification.ID,
			"channel_id":      notification.Channel.ID,
			"type":            notification.Type,
		})
		
		// Send the notification
		err := s.service.sendNotificationToChannel(notification.Incident, notification.Channel, notification.Type)
		
		if err != nil {
			s.logger.Error("Scheduled notification failed", map[string]interface{}{
				"notification_id": notification.ID,
				"channel_id":      notification.Channel.ID,
				"error":           err.Error(),
			})
			
			// Could implement retry logic here for scheduled notifications
		} else {
			s.logger.Info("Scheduled notification sent successfully", map[string]interface{}{
				"notification_id": notification.ID,
				"channel_id":      notification.Channel.ID,
			})
		}
	}
}

// ScheduleRecurringNotification schedules recurring notifications (basic implementation)
func (s *NotificationScheduler) ScheduleRecurringNotification(
	incident *models.Incident,
	channel *models.NotificationChannel,
	notificationType string,
	startTime time.Time,
	interval time.Duration,
	endTime *time.Time,
	maxOccurrences int,
) ([]*ScheduledNotification, error) {
	
	var scheduledNotifications []*ScheduledNotification
	current := startTime
	count := 0
	
	for {
		// Check limits
		if maxOccurrences > 0 && count >= maxOccurrences {
			break
		}
		
		if endTime != nil && current.After(*endTime) {
			break
		}
		
		// Schedule notification
		metadata := map[string]interface{}{
			"recurring":    true,
			"interval":     interval.String(),
			"occurrence":   count + 1,
		}
		
		notification, err := s.ScheduleNotification(incident, channel, notificationType, current, metadata)
		if err != nil {
			return scheduledNotifications, err
		}
		
		scheduledNotifications = append(scheduledNotifications, notification)
		
		// Move to next occurrence
		current = current.Add(interval)
		count++
		
		// Safety check to avoid infinite loops
		if count > 1000 {
			s.logger.Warn("Recurring notification limit reached", map[string]interface{}{
				"channel_id": channel.ID,
				"count":      count,
			})
			break
		}
	}
	
	s.logger.Info("Recurring notifications scheduled", map[string]interface{}{
		"channel_id":      channel.ID,
		"incident_id":     incident.ID,
		"type":            notificationType,
		"count":           len(scheduledNotifications),
		"interval":        interval.String(),
	})
	
	return scheduledNotifications, nil
}

// Stop stops the notification scheduler
func (s *NotificationScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.stopChan)
}

// ErrNotFound represents a not found error
var ErrNotFound = fmt.Errorf("not found")