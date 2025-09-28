package services

import (
	"fmt"
	"strings"
	"text/template"
	"time"
	"bytes"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

// NotificationTemplateService handles notification template management and rendering
type NotificationTemplateService struct {
	logger *Logger
}

// NewNotificationTemplateService creates a new notification template service
func NewNotificationTemplateService(logger *Logger) *NotificationTemplateService {
	return &NotificationTemplateService{
		logger: logger,
	}
}

// TemplateVariables holds the variables available for template rendering
type TemplateVariables struct {
	Incident    *models.Incident
	Timestamp   time.Time
	SystemName  string
	SystemURL   string
	ChannelName string
	Severity    string
	Status      string
	Duration    string
}

// GetDefaultTemplate returns the default template for a given type and channel
func (s *NotificationTemplateService) GetDefaultTemplate(notificationType, channel string) *models.NotificationTemplate {
	templates := s.getBuiltinTemplates()
	key := fmt.Sprintf("%s_%s", notificationType, channel)
	if tmpl, exists := templates[key]; exists {
		return tmpl
	}
	
	// Fallback to generic template
	return s.getGenericTemplate(notificationType, channel)
}

// RenderTemplate renders a notification template with the provided variables
func (s *NotificationTemplateService) RenderTemplate(tmpl *models.NotificationTemplate, vars TemplateVariables) (subject, content string, err error) {
	// Create template functions
	funcMap := template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05 MST")
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
		"duration": func(start, end time.Time) string {
			if end.IsZero() {
				return "N/A"
			}
			return end.Sub(start).String()
		},
	}

	// Render subject if present
	if tmpl.Subject != "" {
		subjTmpl, err := template.New("subject").Funcs(funcMap).Parse(tmpl.Subject)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse subject template: %w", err)
		}
		
		var subjBuf bytes.Buffer
		if err := subjTmpl.Execute(&subjBuf, vars); err != nil {
			return "", "", fmt.Errorf("failed to render subject template: %w", err)
		}
		subject = subjBuf.String()
	}

	// Render body
	bodyTmpl, err := template.New("body").Funcs(funcMap).Parse(tmpl.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse body template: %w", err)
	}
	
	var bodyBuf bytes.Buffer
	if err := bodyTmpl.Execute(&bodyBuf, vars); err != nil {
		return "", "", fmt.Errorf("failed to render body template: %w", err)
	}
	content = bodyBuf.String()

	return subject, content, nil
}

// ValidateTemplate validates a notification template
func (s *NotificationTemplateService) ValidateTemplate(tmpl *models.NotificationTemplate) error {
	if tmpl.Name == "" {
		return fmt.Errorf("template name is required")
	}
	
	if tmpl.Type == "" {
		return fmt.Errorf("template type is required")
	}
	
	if tmpl.Channel == "" {
		return fmt.Errorf("template channel is required")
	}
	
	if tmpl.Body == "" {
		return fmt.Errorf("template body is required")
	}

	// Test template compilation
	_, _, err := s.RenderTemplate(tmpl, TemplateVariables{
		Incident: &models.Incident{
			Title: "Test Incident",
			Severity: models.SeverityHigh,
			Status: models.IncidentStatusOpen,
			CreatedAt: time.Now(),
		},
		Timestamp: time.Now(),
		SystemName: "Test System",
		SystemURL: "https://example.com",
		ChannelName: "test-channel",
	})
	
	if err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	return nil
}

// getBuiltinTemplates returns the built-in default templates
func (s *NotificationTemplateService) getBuiltinTemplates() map[string]*models.NotificationTemplate {
	now := time.Now()
	
	return map[string]*models.NotificationTemplate{
		"incident_created_slack": {
			ID:        "default_incident_created_slack",
			Name:      "Default Incident Created - Slack",
			Type:      "incident_created",
			Channel:   "slack",
			Subject:   "",
			Body:      "🚨 *New Incident Created*\n\n*Title:* {{.Incident.Title}}\n*Severity:* {{.Incident.Severity | upper}}\n*Status:* {{.Incident.Status}}\n*Created:* {{formatTime .Incident.CreatedAt}}\n\n*Description:* {{.Incident.Description}}",
			IsDefault: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		"incident_created_email": {
			ID:        "default_incident_created_email",
			Name:      "Default Incident Created - Email",
			Type:      "incident_created",
			Channel:   "email",
			Subject:   "🚨 New Incident: {{.Incident.Title}}",
			Body:      "A new incident has been created in {{.SystemName}}.\n\nTitle: {{.Incident.Title}}\nSeverity: {{.Incident.Severity | upper}}\nStatus: {{.Incident.Status}}\nCreated: {{formatTime .Incident.CreatedAt}}\n\nDescription:\n{{.Incident.Description}}\n\nView incident: {{.SystemURL}}/incidents/{{.Incident.ID}}",
			IsDefault: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		"incident_created_telegram": {
			ID:        "default_incident_created_telegram",
			Name:      "Default Incident Created - Telegram",
			Type:      "incident_created",
			Channel:   "telegram",
			Subject:   "",
			Body:      "🚨 <b>New Incident Created</b>\n\n<b>Title:</b> {{.Incident.Title}}\n<b>Severity:</b> {{.Incident.Severity | upper}}\n<b>Status:</b> {{.Incident.Status}}\n<b>Created:</b> {{formatTime .Incident.CreatedAt}}\n\n<b>Description:</b> {{.Incident.Description}}",
			IsDefault: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		"incident_acknowledged_slack": {
			ID:        "default_incident_acknowledged_slack",
			Name:      "Default Incident Acknowledged - Slack",
			Type:      "incident_acknowledged",
			Channel:   "slack",
			Subject:   "",
			Body:      "✅ *Incident Acknowledged*\n\n*Title:* {{.Incident.Title}}\n*Status:* {{.Incident.Status}}\n*Acknowledged:* {{formatTime .Incident.AckedAt}}\n*Assignee:* {{.Incident.AssigneeID}}",
			IsDefault: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		"incident_resolved_slack": {
			ID:        "default_incident_resolved_slack",
			Name:      "Default Incident Resolved - Slack",
			Type:      "incident_resolved",
			Channel:   "slack",
			Subject:   "",
			Body:      "🎉 *Incident Resolved*\n\n*Title:* {{.Incident.Title}}\n*Status:* {{.Incident.Status}}\n*Resolved:* {{formatTime .Incident.ResolvedAt}}\n*Duration:* {{duration .Incident.CreatedAt .Incident.ResolvedAt}}",
			IsDefault: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// getGenericTemplate returns a generic fallback template
func (s *NotificationTemplateService) getGenericTemplate(notificationType, channel string) *models.NotificationTemplate {
	now := time.Now()
	
	return &models.NotificationTemplate{
		ID:        fmt.Sprintf("generic_%s_%s", notificationType, channel),
		Name:      fmt.Sprintf("Generic %s - %s", notificationType, channel),
		Type:      notificationType,
		Channel:   channel,
		Subject:   "Incident Alert: {{.Incident.Title}}",
		Body:      "Incident: {{.Incident.Title}}\nSeverity: {{.Incident.Severity}}\nStatus: {{.Incident.Status}}\nTime: {{formatTime .Timestamp}}",
		IsDefault: true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}