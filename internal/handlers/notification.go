package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

// NotificationHandlers handles notification-related HTTP endpoints
type NotificationHandlers struct {
	store                    storage.Store
	notificationService      *services.NotificationService
	templateService         *services.NotificationTemplateService
	logger                  *services.Logger
}

// NewNotificationHandlers creates new notification handlers
func NewNotificationHandlers(
	store storage.Store,
	notificationService *services.NotificationService,
	templateService *services.NotificationTemplateService,
	logger *services.Logger,
) *NotificationHandlers {
	return &NotificationHandlers{
		store:               store,
		notificationService: notificationService,
		templateService:    templateService,
		logger:             logger,
	}
}

// CreateNotificationChannel creates a new notification channel
func (h *NotificationHandlers) CreateNotificationChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var channel models.NotificationChannel
	if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if channel.Name == "" {
		http.Error(w, "Channel name is required", http.StatusBadRequest)
		return
	}
	
	if channel.Type == "" {
		http.Error(w, "Channel type is required", http.StatusBadRequest)
		return
	}

	// Set defaults
	if channel.ID == "" {
		channel.ID = uuid.New().String()
	}
	
	now := time.Now()
	channel.CreatedAt = now
	channel.UpdatedAt = now
	
	if channel.Config == nil {
		channel.Config = make(map[string]string)
	}

	// Validate channel type
	validTypes := map[string]bool{"slack": true, "email": true, "telegram": true}
	if !validTypes[channel.Type] {
		http.Error(w, "Invalid channel type. Must be one of: slack, email, telegram", http.StatusBadRequest)
		return
	}

	// Create channel
	if err := h.store.CreateNotificationChannel(&channel); err != nil {
		h.logger.Error("Failed to create notification channel", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(channel)
}

// GetNotificationChannels returns all notification channels
func (h *NotificationHandlers) GetNotificationChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	channels, err := h.store.ListNotificationChannels()
	if err != nil {
		h.logger.Error("Failed to list notification channels", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channels)
}

// GetNotificationChannel returns a specific notification channel
func (h *NotificationHandlers) GetNotificationChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract channel ID from URL path
	channelID := r.URL.Path[len("/api/notification-channels/"):]
	if channelID == "" {
		http.Error(w, "Channel ID is required", http.StatusBadRequest)
		return
	}

	channel, err := h.store.GetNotificationChannel(channelID)
	if err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, "Channel not found", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to get notification channel", map[string]interface{}{
				"channel_id": channelID,
				"error":      err.Error(),
			})
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channel)
}

// UpdateNotificationChannel updates a notification channel
func (h *NotificationHandlers) UpdateNotificationChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract channel ID from URL path
	channelID := r.URL.Path[len("/api/notification-channels/"):]
	if channelID == "" {
		http.Error(w, "Channel ID is required", http.StatusBadRequest)
		return
	}

	var channel models.NotificationChannel
	if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Ensure ID matches URL
	channel.ID = channelID
	channel.UpdatedAt = time.Now()

	if err := h.store.UpdateNotificationChannel(&channel); err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, "Channel not found", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to update notification channel", map[string]interface{}{
				"channel_id": channelID,
				"error":      err.Error(),
			})
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channel)
}

// DeleteNotificationChannel deletes a notification channel
func (h *NotificationHandlers) DeleteNotificationChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract channel ID from URL path
	channelID := r.URL.Path[len("/api/notification-channels/"):]
	if channelID == "" {
		http.Error(w, "Channel ID is required", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteNotificationChannel(channelID); err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, "Channel not found", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to delete notification channel", map[string]interface{}{
				"channel_id": channelID,
				"error":      err.Error(),
			})
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// TestNotificationChannel tests a notification channel
func (h *NotificationHandlers) TestNotificationChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract channel ID from URL path
	channelID := r.URL.Path[len("/api/notification-channels/"):]
	channelID = channelID[:len(channelID)-5] // Remove "/test" suffix
	
	if channelID == "" {
		http.Error(w, "Channel ID is required", http.StatusBadRequest)
		return
	}

	// Get the channel
	channel, err := h.store.GetNotificationChannel(channelID)
	if err != nil {
		if err == storage.ErrNotFound {
			http.Error(w, "Channel not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Create a test incident
	testIncident := &models.Incident{
		ID:          "test-incident",
		Title:       "Test Notification",
		Description: "This is a test notification to verify channel configuration",
		Status:      models.IncidentStatusOpen,
		Severity:    models.SeverityLow,
		CreatedAt:   time.Now(),
	}

	// Send test notification
	err = h.notificationService.SendTestNotification(testIncident, channel, "test")
	if err != nil {
		h.logger.Error("Test notification failed", map[string]interface{}{
			"channel_id": channelID,
			"error":      err.Error(),
		})
		
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Test notification sent successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetNotificationHistory returns notification history with optional filtering
func (h *NotificationHandlers) GetNotificationHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	incidentID := r.URL.Query().Get("incident_id")
	channelID := r.URL.Query().Get("channel_id")
	_ = r.URL.Query().Get("status") // status filter - unused in this implementation

	// In a real implementation, this would query the database with filters
	// For now, return a placeholder response
	history := []map[string]interface{}{
		{
			"id":           "history-1",
			"incident_id":  incidentID,
			"channel_id":   channelID,
			"type":         "incident_created",
			"channel":      "slack",
			"status":       "sent",
			"sent_at":      time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"created_at":   time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}