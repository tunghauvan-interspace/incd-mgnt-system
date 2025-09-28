package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

// TemplateHandlers handles template-related HTTP endpoints
type TemplateHandlers struct {
	store           storage.Store
	templateService *services.NotificationTemplateService
	logger          *services.Logger
}

// NewTemplateHandlers creates new template handlers
func NewTemplateHandlers(
	store storage.Store,
	templateService *services.NotificationTemplateService,
	logger *services.Logger,
) *TemplateHandlers {
	return &TemplateHandlers{
		store:           store,
		templateService: templateService,
		logger:          logger,
	}
}

// GetTemplates returns all templates with optional filtering
func (h *TemplateHandlers) GetTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	templateType := r.URL.Query().Get("type")
	channel := r.URL.Query().Get("channel")

	// For now, return built-in templates as we don't have database storage yet
	builtinTemplates := h.getBuiltinTemplatesForAPI()
	
	// Filter if requested
	var filteredTemplates []models.NotificationTemplate
	for _, template := range builtinTemplates {
		if templateType != "" && template.Type != templateType {
			continue
		}
		if channel != "" && template.Channel != channel {
			continue
		}
		filteredTemplates = append(filteredTemplates, template)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredTemplates)
}

// GetTemplate returns a specific template
func (h *TemplateHandlers) GetTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract template ID from URL path
	templateID := r.URL.Path[len("/api/templates/"):]
	if templateID == "" {
		http.Error(w, "Template ID is required", http.StatusBadRequest)
		return
	}

	// For now, search in built-in templates
	builtinTemplates := h.getBuiltinTemplatesForAPI()
	for _, template := range builtinTemplates {
		if template.ID == templateID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(template)
			return
		}
	}

	http.Error(w, "Template not found", http.StatusNotFound)
}

// CreateTemplate creates a new custom template
func (h *TemplateHandlers) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var template models.NotificationTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Set defaults
	if template.ID == "" {
		template.ID = uuid.New().String()
	}
	
	now := time.Now()
	template.CreatedAt = now
	template.UpdatedAt = now

	// Validate template
	if err := h.templateService.ValidateTemplate(&template); err != nil {
		http.Error(w, fmt.Sprintf("Template validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// In a real implementation, we would store this in the database
	h.logger.Info("Template created", map[string]interface{}{
		"template_id": template.ID,
		"name":        template.Name,
		"type":        template.Type,
		"channel":     template.Channel,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(template)
}

// UpdateTemplate updates an existing template
func (h *TemplateHandlers) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract template ID from URL path
	templateID := r.URL.Path[len("/api/templates/"):]
	if templateID == "" {
		http.Error(w, "Template ID is required", http.StatusBadRequest)
		return
	}

	var template models.NotificationTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Ensure ID matches URL
	template.ID = templateID
	template.UpdatedAt = time.Now()

	// Validate template
	if err := h.templateService.ValidateTemplate(&template); err != nil {
		http.Error(w, fmt.Sprintf("Template validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// In a real implementation, we would update this in the database
	h.logger.Info("Template updated", map[string]interface{}{
		"template_id": template.ID,
		"name":        template.Name,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

// DeleteTemplate deletes a template
func (h *TemplateHandlers) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract template ID from URL path
	templateID := r.URL.Path[len("/api/templates/"):]
	if templateID == "" {
		http.Error(w, "Template ID is required", http.StatusBadRequest)
		return
	}

	// Check if it's a built-in template (can't delete)
	builtinTemplates := h.getBuiltinTemplatesForAPI()
	for _, template := range builtinTemplates {
		if template.ID == templateID && template.IsDefault {
			http.Error(w, "Cannot delete built-in templates", http.StatusBadRequest)
			return
		}
	}

	// In a real implementation, we would delete from database
	h.logger.Info("Template deleted", map[string]interface{}{
		"template_id": templateID,
	})

	w.WriteHeader(http.StatusNoContent)
}

// PreviewTemplate previews a template with sample data
func (h *TemplateHandlers) PreviewTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Template models.NotificationTemplate `json:"template"`
		SampleData map[string]interface{} `json:"sample_data,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create sample incident data if not provided
	sampleIncident := &models.Incident{
		ID:          "sample-incident-123",
		Title:       "Database Connection Issues",
		Description: "Multiple database connection failures detected",
		Status:      models.IncidentStatusOpen,
		Severity:    models.SeverityHigh,
		CreatedAt:   time.Now(),
		Labels: map[string]string{
			"service":     "database",
			"environment": "production",
		},
	}

	// Create template variables
	vars := services.TemplateVariables{
		Incident:    sampleIncident,
		Timestamp:   time.Now(),
		SystemName:  "Incident Management System",
		SystemURL:   "https://incidents.example.com",
		ChannelName: "general",
		Severity:    string(sampleIncident.Severity),
		Status:      string(sampleIncident.Status),
	}

	// Render template
	subject, content, err := h.templateService.RenderTemplate(&request.Template, vars)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template rendering failed: %v", err), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"subject": subject,
		"content": content,
		"variables_used": map[string]interface{}{
			"incident_id":    sampleIncident.ID,
			"incident_title": sampleIncident.Title,
			"severity":       sampleIncident.Severity,
			"status":         sampleIncident.Status,
			"created_at":     sampleIncident.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ValidateTemplate validates a template without saving it
func (h *TemplateHandlers) ValidateTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var template models.NotificationTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate template
	err := h.templateService.ValidateTemplate(&template)
	
	response := map[string]interface{}{
		"valid": err == nil,
	}
	
	if err != nil {
		response["error"] = err.Error()
	} else {
		response["message"] = "Template is valid"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getBuiltinTemplatesForAPI returns built-in templates for API responses
func (h *TemplateHandlers) getBuiltinTemplatesForAPI() []models.NotificationTemplate {
	templateMap := map[string]*models.NotificationTemplate{
		"incident_created_slack": h.templateService.GetDefaultTemplate("incident_created", "slack"),
		"incident_created_email": h.templateService.GetDefaultTemplate("incident_created", "email"),
		"incident_created_telegram": h.templateService.GetDefaultTemplate("incident_created", "telegram"),
		"incident_acknowledged_slack": h.templateService.GetDefaultTemplate("incident_acknowledged", "slack"),
		"incident_resolved_slack": h.templateService.GetDefaultTemplate("incident_resolved", "slack"),
	}
	
	var templates []models.NotificationTemplate
	for _, template := range templateMap {
		if template != nil {
			templates = append(templates, *template)
		}
	}
	
	return templates
}