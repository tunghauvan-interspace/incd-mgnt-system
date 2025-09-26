package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/config"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

// NotificationService handles sending notifications
type NotificationService struct {
	config *config.Config
}

// NewNotificationService creates a new notification service
func NewNotificationService(config *config.Config) *NotificationService {
	return &NotificationService{
		config: config,
	}
}

// NotifyIncidentCreated sends notifications when an incident is created
func (s *NotificationService) NotifyIncidentCreated(incident *models.Incident) error {
	message := fmt.Sprintf("🚨 New Incident Created\n\nTitle: %s\nSeverity: %s\nStatus: %s\nCreated: %s\n\nDescription: %s",
		incident.Title,
		incident.Severity,
		incident.Status,
		incident.CreatedAt.Format(time.RFC3339),
		incident.Description)

	return s.sendNotifications(message, incident)
}

// NotifyIncidentAcknowledged sends notifications when an incident is acknowledged
func (s *NotificationService) NotifyIncidentAcknowledged(incident *models.Incident) error {
	message := fmt.Sprintf("✅ Incident Acknowledged\n\nTitle: %s\nStatus: %s\nAcknowledged: %s\nAssignee: %s",
		incident.Title,
		incident.Status,
		incident.AckedAt.Format(time.RFC3339),
		incident.AssigneeID)

	return s.sendNotifications(message, incident)
}

// NotifyIncidentResolved sends notifications when an incident is resolved
func (s *NotificationService) NotifyIncidentResolved(incident *models.Incident) error {
	duration := ""
	if incident.ResolvedAt != nil {
		duration = incident.ResolvedAt.Sub(incident.CreatedAt).String()
	}

	message := fmt.Sprintf("🎉 Incident Resolved\n\nTitle: %s\nStatus: %s\nResolved: %s\nDuration: %s",
		incident.Title,
		incident.Status,
		incident.ResolvedAt.Format(time.RFC3339),
		duration)

	return s.sendNotifications(message, incident)
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

	log.Printf("Telegram notification sent successfully")
	return nil
}