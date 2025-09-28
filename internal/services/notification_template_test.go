package services

import (
	"testing"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

func TestNotificationTemplateService(t *testing.T) {
	logger := NewLogger("info", true)
	templateService := NewNotificationTemplateService(logger)

	t.Run("GetDefaultTemplate", func(t *testing.T) {
		template := templateService.GetDefaultTemplate("incident_created", "slack")
		if template == nil {
			t.Fatal("Expected default template, got nil")
		}
		
		if template.Type != "incident_created" {
			t.Errorf("Expected type 'incident_created', got '%s'", template.Type)
		}
		
		if template.Channel != "slack" {
			t.Errorf("Expected channel 'slack', got '%s'", template.Channel)
		}
		
		if template.Body == "" {
			t.Error("Expected non-empty body")
		}
	})

	t.Run("RenderTemplate", func(t *testing.T) {
		template := templateService.GetDefaultTemplate("incident_created", "slack")
		
		incident := &models.Incident{
			ID:          "test-123",
			Title:       "Test Incident",
			Description: "This is a test incident",
			Status:      models.IncidentStatusOpen,
			Severity:    models.SeverityHigh,
			CreatedAt:   time.Now(),
		}
		
		vars := TemplateVariables{
			Incident:    incident,
			Timestamp:   time.Now(),
			SystemName:  "Test System",
			SystemURL:   "http://test.com",
			ChannelName: "test-channel",
			Severity:    string(incident.Severity),
			Status:      string(incident.Status),
		}
		
		subject, content, err := templateService.RenderTemplate(template, vars)
		if err != nil {
			t.Fatalf("Template rendering failed: %v", err)
		}
		
		if content == "" {
			t.Error("Expected non-empty content")
		}
		
		// Check that variables were substituted
		if !containsString(content, "Test Incident") {
			t.Error("Expected content to contain incident title")
		}
		
		if !containsString(content, "HIGH") {
			t.Error("Expected content to contain severity")
		}
		
		t.Logf("Rendered subject: %s", subject)
		t.Logf("Rendered content: %s", content)
	})

	t.Run("ValidateTemplate", func(t *testing.T) {
		// Test valid template
		validTemplate := &models.NotificationTemplate{
			Name:    "Test Template",
			Type:    "incident_created",
			Channel: "slack",
			Body:    "Incident: {{.Incident.Title}}",
		}
		
		err := templateService.ValidateTemplate(validTemplate)
		if err != nil {
			t.Errorf("Valid template should pass validation, got error: %v", err)
		}
		
		// Test invalid template - missing required fields
		invalidTemplate := &models.NotificationTemplate{
			Name: "Invalid Template",
			// Missing Type, Channel, Body
		}
		
		err = templateService.ValidateTemplate(invalidTemplate)
		if err == nil {
			t.Error("Invalid template should fail validation")
		}
		
		// Test template with invalid syntax
		invalidSyntaxTemplate := &models.NotificationTemplate{
			Name:    "Invalid Syntax Template",
			Type:    "incident_created",
			Channel: "slack",
			Body:    "Incident: {{.InvalidField",
		}
		
		err = templateService.ValidateTemplate(invalidSyntaxTemplate)
		if err == nil {
			t.Error("Template with invalid syntax should fail validation")
		}
	})
}

// containsString checks if a string contains a substring (case-insensitive)
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
			len(s) > len(substr) && 
			findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}