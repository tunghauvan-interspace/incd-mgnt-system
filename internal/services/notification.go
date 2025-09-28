package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/config"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/retry"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

// NotificationService handles sending notifications with enhanced features
type NotificationService struct {
	config                   *config.Config
	store                   storage.Store
	templateService         *NotificationTemplateService
	metricsService          *MetricsService
	logger                  *Logger
	retryer                 *retry.Retryer
	batchProcessor          *NotificationBatchProcessor
}

// NewNotificationService creates a new notification service with enhanced features
func NewNotificationService(
	config *config.Config, 
	store storage.Store, 
	templateService *NotificationTemplateService,
	metricsService *MetricsService,
	logger *Logger,
) *NotificationService {
	// Create retry policy for notifications
	retryPolicy := &retry.RetryPolicy{
		MaxAttempts: 3,
		BaseDelay:   2 * time.Second,
		MaxDelay:    30 * time.Second,
		Multiplier:  2.0,
	}
	
	retryer := retry.NewRetryer(retryPolicy, retry.DefaultIsRetryable)
	
	service := &NotificationService{
		config:          config,
		store:           store,
		templateService: templateService,
		metricsService:  metricsService,
		logger:          logger,
		retryer:         retryer,
	}
	
	// Initialize batch processor
	service.batchProcessor = NewNotificationBatchProcessor(service, logger)
	
	return service
}

// NotifyIncidentCreated sends notifications when an incident is created using templates
func (s *NotificationService) NotifyIncidentCreated(incident *models.Incident) error {
	return s.sendTemplatedNotification(incident, "incident_created")
}

// NotifyIncidentAcknowledged sends notifications when an incident is acknowledged using templates
func (s *NotificationService) NotifyIncidentAcknowledged(incident *models.Incident) error {
	return s.sendTemplatedNotification(incident, "incident_acknowledged")
}

// NotifyIncidentResolved sends notifications when an incident is resolved using templates
func (s *NotificationService) NotifyIncidentResolved(incident *models.Incident) error {
	return s.sendTemplatedNotification(incident, "incident_resolved")
}

// sendTemplatedNotification sends notifications using templates and enhanced delivery tracking
func (s *NotificationService) sendTemplatedNotification(incident *models.Incident, notificationType string) error {
	// Get enabled notification channels
	channels, err := s.store.ListNotificationChannels()
	if err != nil {
		s.logger.Error("Failed to get notification channels", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	var errors []string
	
	for _, channel := range channels {
		if !channel.Enabled {
			continue
		}
		
		// Check user preferences if applicable
		if !s.shouldNotify(channel, incident, notificationType) {
			continue
		}
		
		// Check if batching is enabled
		if channel.Preferences != nil && channel.Preferences.BatchingEnabled {
			if err := s.batchProcessor.AddToBatch(incident, channel, notificationType); err != nil {
				s.logger.Error("Failed to add notification to batch", map[string]interface{}{
					"channel_id":        channel.ID,
					"notification_type": notificationType,
					"error":            err.Error(),
				})
				errors = append(errors, fmt.Sprintf("Batching error for %s: %v", channel.Name, err))
			}
			continue
		}
		
		// Send immediately
		if err := s.sendNotificationToChannel(incident, channel, notificationType); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", channel.Name, err))
		}
	}

	// Fallback to legacy config-based notifications if no channels configured
	if len(channels) == 0 {
		return s.sendNotifications(s.generateLegacyMessage(incident, notificationType), incident)
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %s", strings.Join(errors, ", "))
	}

	return nil
}

// sendNotificationToChannel sends a notification to a specific channel with template support
func (s *NotificationService) sendNotificationToChannel(incident *models.Incident, channel *models.NotificationChannel, notificationType string) error {
	// Create notification history entry
	history := &models.NotificationHistory{
		ID:         uuid.New().String(),
		IncidentID: incident.ID,
		ChannelID:  channel.ID,
		Type:       notificationType,
		Channel:    channel.Type,
		Status:     models.DeliveryStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Get template (custom or default)
	template := s.getTemplateForChannel(channel, notificationType)
	if template != nil {
		history.TemplateID = template.ID
	}

	// Store notification history
	if err := s.storeNotificationHistory(history); err != nil {
		s.logger.Error("Failed to store notification history", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Send notification with retry
	ctx := context.Background()
	err := s.retryer.Execute(ctx, func() error {
		return s.deliverNotification(history, template, incident, channel)
	})

	// Update history status
	if err != nil {
		history.Status = models.DeliveryStatusFailed
		history.ErrorMsg = err.Error()
		history.UpdatedAt = time.Now()
		
		s.metricsService.RecordNotificationSent(channel.Type, "failed")
		s.logger.Error("Notification delivery failed", map[string]interface{}{
			"channel_id":        channel.ID,
			"channel_type":      channel.Type,
			"notification_type": notificationType,
			"error":            err.Error(),
		})
	} else {
		history.Status = models.DeliveryStatusSent
		now := time.Now()
		history.SentAt = &now
		history.UpdatedAt = now
		
		s.metricsService.RecordNotificationSent(channel.Type, "sent")
		s.logger.Info("Notification sent successfully", map[string]interface{}{
			"channel_id":        channel.ID,
			"channel_type":      channel.Type,
			"notification_type": notificationType,
		})
	}

	// Update history
	if updateErr := s.updateNotificationHistory(history); updateErr != nil {
		s.logger.Error("Failed to update notification history", map[string]interface{}{
			"error": updateErr.Error(),
		})
	}

	return err
}

// deliverNotification performs the actual notification delivery
func (s *NotificationService) deliverNotification(history *models.NotificationHistory, template *models.NotificationTemplate, incident *models.Incident, channel *models.NotificationChannel) error {
	// Prepare template variables
	vars := TemplateVariables{
		Incident:    incident,
		Timestamp:   time.Now(),
		SystemName:  "Incident Management System",
		SystemURL:   "http://localhost:" + s.config.Port, // Use configured port
		ChannelName: channel.Name,
		Severity:    string(incident.Severity),
		Status:      string(incident.Status),
	}
	
	var subject, content string
	var err error
	
	if template != nil {
		// Use template
		subject, content, err = s.templateService.RenderTemplate(template, vars)
		if err != nil {
			return fmt.Errorf("template rendering failed: %w", err)
		}
	} else {
		// Use legacy format
		content = s.generateLegacyMessage(incident, history.Type)
		subject = fmt.Sprintf("Incident Alert: %s", incident.Title)
	}
	
	// Store rendered content in history
	history.Subject = subject
	history.Content = content
	
	// Send based on channel type
	switch channel.Type {
	case "slack":
		return s.sendSlackNotificationWithConfig(content, channel.Config)
	case "email":
		return s.sendEmailNotificationWithConfig(subject, content, channel.Config, incident)
	case "telegram":
		return s.sendTelegramNotificationWithConfig(content, channel.Config)
	default:
		return fmt.Errorf("unsupported channel type: %s", channel.Type)
	}
}

// sendNotifications sends notifications via all configured channels
func (s *NotificationService) sendNotifications(message string, incident *models.Incident) error {
	var errors []string

	// Send Slack notification
	if s.config.SlackToken != "" && s.config.SlackChannel != "" {
		if err := s.sendSlackNotification(message); err != nil {
			errors = append(errors, fmt.Sprintf("Slack: %v", err))
		}
	}

	// Send Email notification
	if s.config.EmailSMTPHost != "" && s.config.EmailUsername != "" {
		if err := s.sendEmailNotification(message, incident); err != nil {
			errors = append(errors, fmt.Sprintf("Email: %v", err))
		}
	}

	// Send Telegram notification
	if s.config.TelegramBotToken != "" && s.config.TelegramChatID != "" {
		if err := s.sendTelegramNotification(message); err != nil {
			errors = append(errors, fmt.Sprintf("Telegram: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %s", strings.Join(errors, ", "))
	}

	return nil
}

// SlackMessage represents a Slack message payload
type SlackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// sendSlackNotification sends a notification to Slack
func (s *NotificationService) sendSlackNotification(message string) error {
	payload := SlackMessage{
		Channel: s.config.SlackChannel,
		Text:    message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.SlackToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API returned status %d", resp.StatusCode)
	}

	log.Printf("Slack notification sent successfully")
	return nil
}

// sendEmailNotification sends an email notification
func (s *NotificationService) sendEmailNotification(message string, incident *models.Incident) error {
	auth := smtp.PlainAuth("", s.config.EmailUsername, s.config.EmailPassword, s.config.EmailSMTPHost)

	to := []string{"alerts@example.com"} // This should be configurable
	subject := fmt.Sprintf("Incident Alert: %s", incident.Title)
	
	body := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, message)

	addr := fmt.Sprintf("%s:%d", s.config.EmailSMTPHost, s.config.EmailSMTPPort)
	err := smtp.SendMail(addr, auth, s.config.EmailUsername, to, []byte(body))
	if err != nil {
		return err
	}

	log.Printf("Email notification sent successfully")
	return nil
}

// TelegramMessage represents a Telegram message payload
type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// sendTelegramNotification sends a notification to Telegram
func (s *NotificationService) sendTelegramNotification(message string) error {
	payload := TelegramMessage{
		ChatID: s.config.TelegramChatID,
		Text:   message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.config.TelegramBotToken)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	return nil
}

// sendSlackNotificationWithConfig sends a notification to Slack with channel-specific config
func (s *NotificationService) sendSlackNotificationWithConfig(message string, config map[string]string) error {
	token := config["token"]
	channel := config["channel"]
	
	// Fall back to global config if not provided
	if token == "" {
		token = s.config.SlackToken
	}
	if channel == "" {
		channel = s.config.SlackChannel
	}
	
	if token == "" || channel == "" {
		return fmt.Errorf("slack token and channel are required")
	}
	
	payload := SlackMessage{
		Channel: channel,
		Text:    message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API returned status %d", resp.StatusCode)
	}

	return nil
}

// sendEmailNotificationWithConfig sends an email notification with channel-specific config
func (s *NotificationService) sendEmailNotificationWithConfig(subject, message string, config map[string]string, incident *models.Incident) error {
	smtpHost := config["smtp_host"]
	smtpPort := config["smtp_port"]
	username := config["username"]
	password := config["password"]
	from := config["from"]
	to := config["to"]
	
	// Fall back to global config if not provided
	if smtpHost == "" {
		smtpHost = s.config.EmailSMTPHost
	}
	if smtpPort == "" {
		smtpPort = fmt.Sprintf("%d", s.config.EmailSMTPPort)
	}
	if username == "" {
		username = s.config.EmailUsername
	}
	if password == "" {
		password = s.config.EmailPassword
	}
	if from == "" {
		from = s.config.EmailFrom
		if from == "" {
			from = username
		}
	}
	if to == "" {
		to = s.config.EmailTo
		if to == "" {
			to = "alerts@example.com" // default fallback
		}
	}
	
	if smtpHost == "" || username == "" || password == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	auth := smtp.PlainAuth("", username, password, smtpHost)
	
	recipients := strings.Split(to, ",")
	for i, recipient := range recipients {
		recipients[i] = strings.TrimSpace(recipient)
	}
	
	emailBody := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, message)

	port := 587
	if smtpPort != "" {
		fmt.Sscanf(smtpPort, "%d", &port)
	}
	
	addr := fmt.Sprintf("%s:%d", smtpHost, port)
	err := smtp.SendMail(addr, auth, from, recipients, []byte(emailBody))
	if err != nil {
		return err
	}

	return nil
}

// sendTelegramNotificationWithConfig sends a notification to Telegram with channel-specific config
func (s *NotificationService) sendTelegramNotificationWithConfig(message string, config map[string]string) error {
	botToken := config["bot_token"]
	chatID := config["chat_id"]
	
	// Fall back to global config if not provided
	if botToken == "" {
		botToken = s.config.TelegramBotToken
	}
	if chatID == "" {
		chatID = s.config.TelegramChatID
	}
	
	if botToken == "" || chatID == "" {
		return fmt.Errorf("telegram bot token and chat ID are required")
	}
	
	payload := TelegramMessage{
		ChatID: chatID,
		Text:   message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	return nil
}

// shouldNotify checks if a notification should be sent based on preferences
func (s *NotificationService) shouldNotify(channel *models.NotificationChannel, incident *models.Incident, notificationType string) bool {
	if channel.Preferences == nil {
		return true
	}
	
	// Check opt-in status
	if !channel.Preferences.OptIn {
		return false
	}
	
	// Check severity filter
	if len(channel.Preferences.SeverityFilter) > 0 {
		severityMatch := false
		for _, severity := range channel.Preferences.SeverityFilter {
			if strings.EqualFold(severity, string(incident.Severity)) {
				severityMatch = true
				break
			}
		}
		if !severityMatch {
			return false
		}
	}
	
	// Check quiet hours
	if channel.Preferences.QuietHours != nil && channel.Preferences.QuietHours.Enabled {
		if s.isInQuietHours(channel.Preferences.QuietHours) {
			return false
		}
	}
	
	return true
}

// isInQuietHours checks if current time is within quiet hours
func (s *NotificationService) isInQuietHours(config *models.QuietHoursConfig) bool {
	// Simple implementation - can be enhanced with timezone support
	now := time.Now()
	currentHour := now.Hour()
	currentDay := int(now.Weekday())
	
	// Check if current day is in quiet days
	if len(config.Days) > 0 {
		dayMatch := false
		for _, day := range config.Days {
			if day == currentDay {
				dayMatch = true
				break
			}
		}
		if !dayMatch {
			return false
		}
	}
	
	// Parse start and end times (simplified)
	startHour := 0
	endHour := 24
	if config.StartTime != "" {
		fmt.Sscanf(config.StartTime, "%d:", &startHour)
	}
	if config.EndTime != "" {
		fmt.Sscanf(config.EndTime, "%d:", &endHour)
	}
	
	// Check if current time is within quiet hours
	if startHour <= endHour {
		return currentHour >= startHour && currentHour < endHour
	} else {
		// Overnight quiet hours (e.g., 22:00 - 06:00)
		return currentHour >= startHour || currentHour < endHour
	}
}

// getTemplateForChannel gets the appropriate template for a channel and notification type
func (s *NotificationService) getTemplateForChannel(channel *models.NotificationChannel, notificationType string) *models.NotificationTemplate {
	// Check if channel has custom templates
	if channel.Templates != nil {
		if templateContent, exists := channel.Templates[notificationType]; exists && templateContent != "" {
			return &models.NotificationTemplate{
				ID:      fmt.Sprintf("%s_%s_%s", channel.ID, notificationType, "custom"),
				Name:    fmt.Sprintf("Custom %s - %s", notificationType, channel.Name),
				Type:    notificationType,
				Channel: channel.Type,
				Body:    templateContent,
			}
		}
	}
	
	// Fall back to default template
	return s.templateService.GetDefaultTemplate(notificationType, channel.Type)
}

// generateLegacyMessage generates a legacy-format message for backward compatibility
func (s *NotificationService) generateLegacyMessage(incident *models.Incident, notificationType string) string {
	switch notificationType {
	case "incident_created":
		return fmt.Sprintf("🚨 New Incident Created\n\nTitle: %s\nSeverity: %s\nStatus: %s\nCreated: %s\n\nDescription: %s",
			incident.Title,
			incident.Severity,
			incident.Status,
			incident.CreatedAt.Format(time.RFC3339),
			incident.Description)
	case "incident_acknowledged":
		ackedTime := incident.CreatedAt.Format(time.RFC3339) // fallback
		if incident.AckedAt != nil {
			ackedTime = incident.AckedAt.Format(time.RFC3339)
		}
		return fmt.Sprintf("✅ Incident Acknowledged\n\nTitle: %s\nStatus: %s\nAcknowledged: %s\nAssignee: %s",
			incident.Title,
			incident.Status,
			ackedTime,
			incident.AssigneeID)
	case "incident_resolved":
		duration := ""
		resolvedTime := incident.CreatedAt.Format(time.RFC3339) // fallback
		if incident.ResolvedAt != nil {
			resolvedTime = incident.ResolvedAt.Format(time.RFC3339)
			duration = incident.ResolvedAt.Sub(incident.CreatedAt).String()
		}
		return fmt.Sprintf("🎉 Incident Resolved\n\nTitle: %s\nStatus: %s\nResolved: %s\nDuration: %s",
			incident.Title,
			incident.Status,
			resolvedTime,
			duration)
	default:
		return fmt.Sprintf("Incident Update: %s\nTitle: %s\nStatus: %s\nSeverity: %s",
			notificationType, incident.Title, incident.Status, incident.Severity)
	}
}

// storeNotificationHistory stores notification history (placeholder implementation)
func (s *NotificationService) storeNotificationHistory(history *models.NotificationHistory) error {
	// In a real implementation, this would store to the database
	// For now, we'll just log it
	s.logger.Info("Storing notification history", map[string]interface{}{
		"id":         history.ID,
		"incident_id": history.IncidentID,
		"channel_id": history.ChannelID,
		"type":       history.Type,
		"status":     history.Status,
	})
	return nil
}

// updateNotificationHistory updates notification history (placeholder implementation)
func (s *NotificationService) updateNotificationHistory(history *models.NotificationHistory) error {
	// In a real implementation, this would update the database record
	s.logger.Info("Updating notification history", map[string]interface{}{
		"id":     history.ID,
		"status": history.Status,
		"error":  history.ErrorMsg,
	})
	return nil
}

// SendTestNotification sends a test notification to verify channel configuration
func (s *NotificationService) SendTestNotification(incident *models.Incident, channel *models.NotificationChannel, notificationType string) error {
	// Create test notification history entry
	history := &models.NotificationHistory{
		ID:         uuid.New().String(),
		IncidentID: incident.ID,
		ChannelID:  channel.ID,
		Type:       notificationType,
		Channel:    channel.Type,
		Status:     models.DeliveryStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Get template for test
	template := s.getTemplateForChannel(channel, "incident_created") // Use incident_created template for test
	
	// Deliver the test notification
	err := s.deliverNotification(history, template, incident, channel)
	
	if err != nil {
		s.logger.Error("Test notification failed", map[string]interface{}{
			"channel_id":   channel.ID,
			"channel_type": channel.Type,
			"error":        err.Error(),
		})
		return err
	}

	s.logger.Info("Test notification sent successfully", map[string]interface{}{
		"channel_id":   channel.ID,
		"channel_type": channel.Type,
	})

	return nil
}