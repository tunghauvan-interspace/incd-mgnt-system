package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

// NotificationBatchProcessor handles batching of notifications for efficient delivery
type NotificationBatchProcessor struct {
	service      *NotificationService
	logger       *Logger
	batches      map[string]*models.NotificationBatch
	mutex        sync.RWMutex
	ticker       *time.Ticker
	stopChan     chan bool
	batchTimeout time.Duration
}

// NewNotificationBatchProcessor creates a new batch processor
func NewNotificationBatchProcessor(service *NotificationService, logger *Logger) *NotificationBatchProcessor {
	processor := &NotificationBatchProcessor{
		service:      service,
		logger:       logger,
		batches:      make(map[string]*models.NotificationBatch),
		stopChan:     make(chan bool),
		batchTimeout: 5 * time.Minute, // default batch timeout
	}
	
	// Start the batch processing ticker
	processor.ticker = time.NewTicker(1 * time.Minute) // check every minute
	go processor.processBatches()
	
	return processor
}

// AddToBatch adds a notification to a batch for the given channel
func (bp *NotificationBatchProcessor) AddToBatch(incident *models.Incident, channel *models.NotificationChannel, notificationType string) error {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()
	
	batchKey := fmt.Sprintf("%s_%s", channel.ID, notificationType)
	
	// Get or create batch
	batch, exists := bp.batches[batchKey]
	if !exists {
		batch = &models.NotificationBatch{
			ID:            uuid.New().String(),
			ChannelID:     channel.ID,
			Type:          notificationType,
			Count:         0,
			Status:        models.DeliveryStatusPending,
			Notifications: make([]string, 0),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		bp.batches[batchKey] = batch
	}
	
	// Create notification history entry for the batch item
	historyID := uuid.New().String()
	history := &models.NotificationHistory{
		ID:         historyID,
		IncidentID: incident.ID,
		ChannelID:  channel.ID,
		Type:       notificationType,
		Channel:    channel.Type,
		Status:     models.DeliveryStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	// Store notification history
	if err := bp.service.storeNotificationHistory(history); err != nil {
		return fmt.Errorf("failed to store notification history: %w", err)
	}
	
	// Add to batch
	batch.Notifications = append(batch.Notifications, historyID)
	batch.Count++
	batch.UpdatedAt = time.Now()
	
	bp.logger.Info("Added notification to batch", map[string]interface{}{
		"batch_id":    batch.ID,
		"channel_id":  channel.ID,
		"type":        notificationType,
		"batch_size":  batch.Count,
	})
	
	// Check if batch should be processed immediately
	maxBatchSize := 10 // default
	if channel.Preferences != nil && channel.Preferences.MaxBatchSize > 0 {
		maxBatchSize = channel.Preferences.MaxBatchSize
	}
	
	if batch.Count >= maxBatchSize {
		return bp.processBatch(batch, channel)
	}
	
	return nil
}

// processBatches periodically processes pending batches
func (bp *NotificationBatchProcessor) processBatches() {
	for {
		select {
		case <-bp.ticker.C:
			bp.processTimedOutBatches()
		case <-bp.stopChan:
			return
		}
	}
}

// processTimedOutBatches processes batches that have timed out
func (bp *NotificationBatchProcessor) processTimedOutBatches() {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()
	
	now := time.Now()
	
	for batchKey, batch := range bp.batches {
		// Check if batch has timed out
		if now.Sub(batch.CreatedAt) > bp.batchTimeout {
			// Get channel info for processing
			channel, err := bp.service.store.GetNotificationChannel(batch.ChannelID)
			if err != nil {
				bp.logger.Error("Failed to get channel for batch processing", map[string]interface{}{
					"batch_id":   batch.ID,
					"channel_id": batch.ChannelID,
					"error":      err.Error(),
				})
				continue
			}
			
			// Process the batch
			if err := bp.processBatch(batch, channel); err != nil {
				bp.logger.Error("Failed to process timed out batch", map[string]interface{}{
					"batch_id": batch.ID,
					"error":    err.Error(),
				})
			} else {
				// Remove processed batch
				delete(bp.batches, batchKey)
			}
		}
	}
}

// processBatch processes a single batch
func (bp *NotificationBatchProcessor) processBatch(batch *models.NotificationBatch, channel *models.NotificationChannel) error {
	bp.logger.Info("Processing notification batch", map[string]interface{}{
		"batch_id":   batch.ID,
		"channel_id": channel.ID,
		"count":      batch.Count,
		"type":       batch.Type,
	})
	
	// Update batch status
	batch.Status = models.DeliveryStatusSent
	now := time.Now()
	batch.ProcessedAt = &now
	batch.UpdatedAt = now
	
	// Get template for batched notifications
	template := bp.service.getTemplateForChannel(channel, batch.Type)
	
	// Create batched content
	incidents := make([]*models.Incident, 0, batch.Count)
	for _, historyID := range batch.Notifications {
		// In a real implementation, we would fetch the incident from the history
		// For now, we'll create a placeholder
		incident := &models.Incident{
			ID:          fmt.Sprintf("incident_%s", historyID),
			Title:       fmt.Sprintf("Batched incident %s", historyID[:8]),
			Description: "Batched notification",
			Severity:    models.SeverityMedium,
			Status:      models.IncidentStatusOpen,
			CreatedAt:   time.Now(),
		}
		incidents = append(incidents, incident)
	}
	
	// Create batched message
	content := bp.createBatchedMessage(incidents, batch.Type, template)
	subject := fmt.Sprintf("Batched %s Notifications (%d)", batch.Type, batch.Count)
	
	// Send batched notification
	var err error
	switch channel.Type {
	case "slack":
		err = bp.service.sendSlackNotificationWithConfig(content, channel.Config)
	case "email":
		err = bp.service.sendEmailNotificationWithConfig(subject, content, channel.Config, nil)
	case "telegram":
		err = bp.service.sendTelegramNotificationWithConfig(content, channel.Config)
	default:
		err = fmt.Errorf("unsupported channel type: %s", channel.Type)
	}
	
	if err != nil {
		batch.Status = models.DeliveryStatusFailed
		bp.logger.Error("Failed to send batched notification", map[string]interface{}{
			"batch_id": batch.ID,
			"error":    err.Error(),
		})
		return err
	}
	
	// Update all individual notification histories
	for _, historyID := range batch.Notifications {
		// In a real implementation, we would update the database
		bp.logger.Info("Updated batched notification status", map[string]interface{}{
			"history_id": historyID,
			"status":     "sent",
			"batch_id":   batch.ID,
		})
	}
	
	bp.logger.Info("Batch processed successfully", map[string]interface{}{
		"batch_id":   batch.ID,
		"channel_id": channel.ID,
		"count":      batch.Count,
	})
	
	return nil
}

// createBatchedMessage creates a message for batched notifications
func (bp *NotificationBatchProcessor) createBatchedMessage(incidents []*models.Incident, notificationType string, template *models.NotificationTemplate) string {
	if len(incidents) == 0 {
		return "No incidents to report"
	}
	
	if len(incidents) == 1 {
		// Single incident, use normal template
		vars := TemplateVariables{
			Incident:  incidents[0],
			Timestamp: time.Now(),
		}
		
		if template != nil {
			_, content, err := bp.service.templateService.RenderTemplate(template, vars)
			if err == nil {
				return content
			}
		}
		
		return bp.service.generateLegacyMessage(incidents[0], notificationType)
	}
	
	// Multiple incidents, create summary
	var content strings.Builder
	
	switch notificationType {
	case "incident_created":
		content.WriteString(fmt.Sprintf("🚨 **%d New Incidents Created**\n\n", len(incidents)))
	case "incident_acknowledged":
		content.WriteString(fmt.Sprintf("✅ **%d Incidents Acknowledged**\n\n", len(incidents)))
	case "incident_resolved":
		content.WriteString(fmt.Sprintf("🎉 **%d Incidents Resolved**\n\n", len(incidents)))
	default:
		content.WriteString(fmt.Sprintf("📋 **%d Incident Updates**\n\n", len(incidents)))
	}
	
	// Add summary of incidents
	for i, incident := range incidents {
		if i >= 10 { // Limit to first 10 incidents
			content.WriteString(fmt.Sprintf("... and %d more incidents\n", len(incidents)-10))
			break
		}
		
		content.WriteString(fmt.Sprintf("• **%s** (%s) - %s\n", 
			incident.Title, 
			incident.Severity, 
			incident.Status))
	}
	
	content.WriteString(fmt.Sprintf("\n*Batched at: %s*", time.Now().Format("2006-01-02 15:04:05")))
	
	return content.String()
}

// Stop stops the batch processor
func (bp *NotificationBatchProcessor) Stop() {
	if bp.ticker != nil {
		bp.ticker.Stop()
	}
	close(bp.stopChan)
}