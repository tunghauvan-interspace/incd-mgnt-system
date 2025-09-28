package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

// IncidentService handles incident operations
type IncidentService struct {
	store          storage.Store
	metricsService *MetricsService
}

// NewIncidentService creates a new incident service
func NewIncidentService(store storage.Store, metricsService *MetricsService) *IncidentService {
	return &IncidentService{
		store:          store,
		metricsService: metricsService,
	}
}

// CreateIncident creates a new incident
func (s *IncidentService) CreateIncident(title, description string, severity models.IncidentSeverity, alertIDs []string) (*models.Incident, error) {
	incident := &models.Incident{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      models.IncidentStatusOpen,
		Severity:    severity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AlertIDs:    alertIDs,
		Labels:      make(map[string]string),
	}

	start := time.Now()
	err := s.store.CreateIncident(incident)
	if s.metricsService != nil {
		s.metricsService.RecordDBQuery("CREATE", "incidents", time.Since(start))
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to create incident: %w", err)
	}

	// Record metrics
	if s.metricsService != nil {
		s.metricsService.RecordIncidentCreated(string(severity), string(incident.Status))
	}

	return incident, nil
}

// GetIncident retrieves an incident by ID
func (s *IncidentService) GetIncident(id string) (*models.Incident, error) {
	return s.store.GetIncident(id)
}

// ListIncidents retrieves all incidents
func (s *IncidentService) ListIncidents() ([]*models.Incident, error) {
	return s.store.ListIncidents()
}

// AcknowledgeIncident acknowledges an incident
func (s *IncidentService) AcknowledgeIncident(id, assigneeID string) error {
	incident, err := s.store.GetIncident(id)
	if err != nil {
		return err
	}

	now := time.Now()
	incident.Status = models.IncidentStatusAcknowledged
	incident.AckedAt = &now
	incident.UpdatedAt = now
	incident.AssigneeID = assigneeID

	return s.store.UpdateIncident(incident)
}

// ResolveIncident resolves an incident
func (s *IncidentService) ResolveIncident(id string) error {
	incident, err := s.store.GetIncident(id)
	if err != nil {
		return err
	}

	now := time.Now()
	incident.Status = models.IncidentStatusResolved
	incident.ResolvedAt = &now
	incident.UpdatedAt = now

	return s.store.UpdateIncident(incident)
}

// UpdateIncident updates an incident
func (s *IncidentService) UpdateIncident(incident *models.Incident) error {
	incident.UpdatedAt = time.Now()
	return s.store.UpdateIncident(incident)
}

// DeleteIncident deletes an incident
func (s *IncidentService) DeleteIncident(id string) error {
	return s.store.DeleteIncident(id)
}

// CalculateMetrics calculates incident metrics
func (s *IncidentService) CalculateMetrics() (*models.Metrics, error) {
	incidents, err := s.store.ListIncidents()
	if err != nil {
		return nil, err
	}

	metrics := &models.Metrics{
		IncidentsByStatus:   make(map[string]int),
		IncidentsBySeverity: make(map[string]int),
	}

	var totalAckTime time.Duration
	var totalResolveTime time.Duration
	var ackCount, resolveCount int

	for _, incident := range incidents {
		metrics.TotalIncidents++

		// Count by status
		metrics.IncidentsByStatus[string(incident.Status)]++
		switch incident.Status {
		case models.IncidentStatusOpen:
			metrics.OpenIncidents++
		case models.IncidentStatusResolved:
			metrics.ResolvedIncidents++
		}

		// Count by severity
		metrics.IncidentsBySeverity[string(incident.Severity)]++

		// Calculate MTTA
		if incident.AckedAt != nil {
			totalAckTime += incident.AckedAt.Sub(incident.CreatedAt)
			ackCount++
		}

		// Calculate MTTR
		if incident.ResolvedAt != nil {
			totalResolveTime += incident.ResolvedAt.Sub(incident.CreatedAt)
			resolveCount++
		}
	}

	// Calculate averages
	if ackCount > 0 {
		metrics.MTTA = totalAckTime / time.Duration(ackCount)
	}
	if resolveCount > 0 {
		metrics.MTTR = totalResolveTime / time.Duration(resolveCount)
	}

	return metrics, nil
}

// UpdatePrometheusMetrics updates Prometheus metrics with current incident data
func (s *IncidentService) UpdatePrometheusMetrics() error {
	if s.metricsService == nil {
		return nil
	}

	start := time.Now()
	metrics, err := s.CalculateMetrics()
	s.metricsService.RecordDBQuery("SELECT", "incidents", time.Since(start))
	
	if err != nil {
		return err
	}

	// Update MTTA and MTTR metrics
	s.metricsService.UpdateMTTA(metrics.MTTA)
	s.metricsService.UpdateMTTR(metrics.MTTR)

	// Update incidents by status and severity
	for status, count := range metrics.IncidentsByStatus {
		for severity, severityCount := range metrics.IncidentsBySeverity {
			// This gives a rough approximation - in a real system you'd want more granular data
			proportion := float64(severityCount) / float64(metrics.TotalIncidents)
			estimatedCount := float64(count) * proportion
			s.metricsService.UpdateIncidentsByStatus(status, severity, estimatedCount)
		}
	}

	return nil
}

// Enhanced Incident Features - Comments and Timeline

// AddComment adds a comment to an incident timeline
func (s *IncidentService) AddComment(incidentID, userID, content string, commentType models.IncidentCommentType, metadata map[string]interface{}) (*models.IncidentComment, error) {
	// Verify incident exists
	_, err := s.store.GetIncident(incidentID)
	if err != nil {
		return nil, fmt.Errorf("incident not found: %w", err)
	}

	comment := &models.IncidentComment{
		ID:          uuid.New().String(),
		IncidentID:  incidentID,
		UserID:      &userID,
		Content:     content,
		CommentType: commentType,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
	}

	if err := s.store.CreateIncidentComment(comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

// GetComments retrieves comments for an incident
func (s *IncidentService) GetComments(incidentID string) ([]*models.IncidentComment, error) {
	return s.store.GetIncidentComments(incidentID)
}

// GetTimeline retrieves the complete timeline for an incident (comments + system events)
func (s *IncidentService) GetTimeline(incidentID string) ([]*models.IncidentComment, error) {
	return s.store.GetIncidentTimeline(incidentID)
}

// Enhanced Incident Features - Tags

// AddTags adds tags to an incident
func (s *IncidentService) AddTags(incidentID, userID string, tags []models.TemplateTag) error {
	// Verify incident exists
	_, err := s.store.GetIncident(incidentID)
	if err != nil {
		return fmt.Errorf("incident not found: %w", err)
	}

	for _, tag := range tags {
		incidentTag := &models.IncidentTag{
			ID:         uuid.New().String(),
			IncidentID: incidentID,
			TagName:    tag.Name,
			TagValue:   &tag.Value,
			Color:      tag.Color,
			CreatedBy:  &userID,
			CreatedAt:  time.Now(),
		}

		if err := s.store.CreateIncidentTag(incidentTag); err != nil {
			return fmt.Errorf("failed to create tag %s: %w", tag.Name, err)
		}

		// Add timeline entry
		metadata := map[string]interface{}{
			"tag_name":  tag.Name,
			"tag_value": tag.Value,
			"color":     tag.Color,
		}
		_, _ = s.AddComment(incidentID, userID, fmt.Sprintf("Added tag: %s", tag.Name), models.CommentTypeTagAdded, metadata)
	}

	return nil
}

// RemoveTags removes tags from an incident
func (s *IncidentService) RemoveTags(incidentID, userID string, tagNames []string) error {
	for _, tagName := range tagNames {
		if err := s.store.DeleteIncidentTag(incidentID, tagName); err != nil {
			return fmt.Errorf("failed to remove tag %s: %w", tagName, err)
		}

		// Add timeline entry
		metadata := map[string]interface{}{
			"tag_name": tagName,
		}
		_, _ = s.AddComment(incidentID, userID, fmt.Sprintf("Removed tag: %s", tagName), models.CommentTypeTagRemoved, metadata)
	}

	return nil
}

// GetTags retrieves all tags for an incident
func (s *IncidentService) GetTags(incidentID string) ([]*models.IncidentTag, error) {
	return s.store.GetIncidentTags(incidentID)
}

// Enhanced Incident Features - Templates

// CreateTemplate creates a new incident template
func (s *IncidentService) CreateTemplate(template *models.IncidentTemplate) error {
	template.ID = uuid.New().String()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	template.IsActive = true

	return s.store.CreateIncidentTemplate(template)
}

// ListTemplates retrieves all active incident templates
func (s *IncidentService) ListTemplates() ([]*models.IncidentTemplate, error) {
	return s.store.ListIncidentTemplates(true) // only active templates
}

// GetTemplate retrieves a specific incident template
func (s *IncidentService) GetTemplate(templateID string) (*models.IncidentTemplate, error) {
	return s.store.GetIncidentTemplate(templateID)
}

// UseTemplate creates an incident from a template
func (s *IncidentService) UseTemplate(req *models.CreateIncidentFromTemplateRequest, userID string) (*models.Incident, error) {
	template, err := s.store.GetIncidentTemplate(req.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	if !template.IsActive {
		return nil, fmt.Errorf("template is not active")
	}

	// Replace variables in title and description
	title := s.replaceVariables(template.TitleTemplate, req.Variables)
	description := s.replaceVariables(template.DescriptionTemplate, req.Variables)

	// Create incident
	incident, err := s.CreateIncident(title, description, template.Severity, []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to create incident from template: %w", err)
	}

	// Assign if specified
	if req.AssigneeID != nil {
		if err := s.AssignIncident(incident.ID, *req.AssigneeID, userID); err != nil {
			// Log error but don't fail the creation
			fmt.Printf("Failed to assign incident: %v\n", err)
		}
	}

	// Add default tags
	allTags := append(template.DefaultTags, req.AdditionalTags...)
	if len(allTags) > 0 {
		if err := s.AddTags(incident.ID, userID, allTags); err != nil {
			// Log error but don't fail the creation
			fmt.Printf("Failed to add template tags: %v\n", err)
		}
	}

	return incident, nil
}

// replaceVariables replaces {{variable}} placeholders in text
func (s *IncidentService) replaceVariables(text string, variables map[string]string) string {
	result := text
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// Enhanced Incident Features - Attachments

// AttachFile attaches a file to an incident
func (s *IncidentService) AttachFile(attachment *models.IncidentAttachment, userID string) error {
	// Verify incident exists
	_, err := s.store.GetIncident(attachment.IncidentID)
	if err != nil {
		return fmt.Errorf("incident not found: %w", err)
	}

	attachment.ID = uuid.New().String()
	attachment.UploadedBy = &userID
	attachment.CreatedAt = time.Now()

	if err := s.store.CreateIncidentAttachment(attachment); err != nil {
		return fmt.Errorf("failed to attach file: %w", err)
	}

	// Add timeline entry
	metadata := map[string]interface{}{
		"file_name":       attachment.OriginalName,
		"file_size":       attachment.FileSize,
		"attachment_type": attachment.AttachmentType,
		"attachment_id":   attachment.ID,
	}
	_, _ = s.AddComment(attachment.IncidentID, userID, fmt.Sprintf("Attached file: %s", attachment.OriginalName), models.CommentTypeAttachmentAdded, metadata)

	return nil
}

// GetAttachments retrieves all attachments for an incident
func (s *IncidentService) GetAttachments(incidentID string) ([]*models.IncidentAttachment, error) {
	return s.store.GetIncidentAttachments(incidentID)
}

// Enhanced Incident Features - Search and Filtering

// SearchIncidents performs full-text search and filtering on incidents
func (s *IncidentService) SearchIncidents(req *models.IncidentSearchRequest) (*models.IncidentSearchResponse, error) {
	incidents, total, err := s.store.SearchIncidents(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search incidents: %w", err)
	}

	totalPages := (total + req.Limit - 1) / req.Limit

	return &models.IncidentSearchResponse{
		Incidents:  incidents,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// Enhanced Incident Features - Bulk Operations

// BulkAcknowledge acknowledges multiple incidents
func (s *IncidentService) BulkAcknowledge(incidentIDs []string, assigneeID, userID string) (*models.BulkOperationResponse, error) {
	return s.performBulkOperation(incidentIDs, func(incidentID string) error {
		return s.AcknowledgeIncident(incidentID, assigneeID)
	})
}

// BulkUpdateStatus updates status for multiple incidents
func (s *IncidentService) BulkUpdateStatus(incidentIDs []string, status models.IncidentStatus, userID string) (*models.BulkOperationResponse, error) {
	return s.performBulkOperation(incidentIDs, func(incidentID string) error {
		incident, err := s.store.GetIncident(incidentID)
		if err != nil {
			return err
		}

		oldStatus := incident.Status
		incident.Status = status
		incident.UpdatedAt = time.Now()

		if status == models.IncidentStatusResolved {
			now := time.Now()
			incident.ResolvedAt = &now
		}

		if err := s.store.UpdateIncident(incident); err != nil {
			return err
		}

		// Add timeline entry
		metadata := map[string]interface{}{
			"old_status": oldStatus,
			"new_status": status,
		}
		_, _ = s.AddComment(incidentID, userID, fmt.Sprintf("Status changed from %s to %s", oldStatus, status), models.CommentTypeStatusChange, metadata)

		return nil
	})
}

// performBulkOperation executes a bulk operation on incidents
func (s *IncidentService) performBulkOperation(incidentIDs []string, operation func(string) error) (*models.BulkOperationResponse, error) {
	response := &models.BulkOperationResponse{}

	for _, incidentID := range incidentIDs {
		if err := operation(incidentID); err != nil {
			response.FailedCount++
			response.Failures = append(response.Failures, models.BulkOperationFailure{
				IncidentID: incidentID,
				Error:      err.Error(),
			})
		} else {
			response.ProcessedCount++
		}
	}

	return response, nil
}

// Enhanced Incident Features - Assignment

// AssignIncident assigns an incident to a user
func (s *IncidentService) AssignIncident(incidentID, assigneeID, userID string) error {
	incident, err := s.store.GetIncident(incidentID)
	if err != nil {
		return fmt.Errorf("incident not found: %w", err)
	}

	oldAssigneeID := incident.AssigneeID
	incident.AssigneeID = assigneeID
	incident.UpdatedAt = time.Now()

	if err := s.store.UpdateIncident(incident); err != nil {
		return fmt.Errorf("failed to assign incident: %w", err)
	}

	// Add timeline entry
	metadata := map[string]interface{}{
		"old_assignee": oldAssigneeID,
		"new_assignee": assigneeID,
	}
	action := "assigned"
	if oldAssigneeID != "" {
		action = "reassigned"
	}
	_, _ = s.AddComment(incidentID, userID, fmt.Sprintf("Incident %s to user %s", action, assigneeID), models.CommentTypeAssignment, metadata)

	return nil
}

// ReassignIncident reassigns an incident to a different user
func (s *IncidentService) ReassignIncident(incidentID, newAssigneeID, userID string) error {
	return s.AssignIncident(incidentID, newAssigneeID, userID)
}