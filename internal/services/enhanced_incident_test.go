package services

import (
	"testing"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

func TestEnhancedIncidentFeatures(t *testing.T) {
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}
	defer store.Close()

	metricsService := NewMetricsService()
	incidentService := NewIncidentService(store, metricsService)

	// Create a test incident
	incident, err := incidentService.CreateIncident(
		"Test Incident",
		"Test Description",
		models.SeverityHigh,
		[]string{},
	)
	if err != nil {
		t.Fatalf("Failed to create incident: %v", err)
	}

	t.Run("Comments and Timeline", func(t *testing.T) {
		// Test adding comments
		comment, err := incidentService.AddComment(
			incident.ID,
			"test-user-1",
			"This is a test comment",
			models.CommentTypeComment,
			nil,
		)
		if err != nil {
			t.Fatalf("Failed to add comment: %v", err)
		}

		if comment.Content != "This is a test comment" {
			t.Errorf("Expected comment content 'This is a test comment', got '%s'", comment.Content)
		}

		// Test retrieving comments
		comments, err := incidentService.GetComments(incident.ID)
		if err != nil {
			t.Fatalf("Failed to get comments: %v", err)
		}

		if len(comments) != 1 {
			t.Errorf("Expected 1 comment, got %d", len(comments))
		}

		// Test timeline
		timeline, err := incidentService.GetTimeline(incident.ID)
		if err != nil {
			t.Fatalf("Failed to get timeline: %v", err)
		}

		if len(timeline) != 1 {
			t.Errorf("Expected 1 timeline event, got %d", len(timeline))
		}
	})

	t.Run("Tags", func(t *testing.T) {
		// Test adding tags
		tags := []models.TemplateTag{
			{Name: "environment", Value: "production", Color: "#ff0000"},
			{Name: "service", Value: "api", Color: "#00ff00"},
		}

		err := incidentService.AddTags(incident.ID, "test-user-1", tags)
		if err != nil {
			t.Fatalf("Failed to add tags: %v", err)
		}

		// Test retrieving tags
		retrievedTags, err := incidentService.GetTags(incident.ID)
		if err != nil {
			t.Fatalf("Failed to get tags: %v", err)
		}

		if len(retrievedTags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(retrievedTags))
		}

		// Test removing tags
		err = incidentService.RemoveTags(incident.ID, "test-user-1", []string{"environment"})
		if err != nil {
			t.Fatalf("Failed to remove tag: %v", err)
		}

		// Verify tag removal
		remainingTags, err := incidentService.GetTags(incident.ID)
		if err != nil {
			t.Fatalf("Failed to get remaining tags: %v", err)
		}

		if len(remainingTags) != 1 {
			t.Errorf("Expected 1 remaining tag, got %d", len(remainingTags))
		}
	})

	t.Run("Templates", func(t *testing.T) {
		// Test creating template
		template := &models.IncidentTemplate{
			Name:                "Test Template",
			Description:         "Test template description",
			TitleTemplate:       "{{service}} Service Issue",
			DescriptionTemplate: "Service {{service}} is experiencing issues.\nImpact: {{impact}}",
			Severity:            models.SeverityMedium,
			DefaultTags: []models.TemplateTag{
				{Name: "template", Value: "test", Color: "#0000ff"},
			},
		}

		err := incidentService.CreateTemplate(template)
		if err != nil {
			t.Fatalf("Failed to create template: %v", err)
		}

		// Test listing templates
		templates, err := incidentService.ListTemplates()
		if err != nil {
			t.Fatalf("Failed to list templates: %v", err)
		}

		if len(templates) == 0 {
			t.Error("Expected at least 1 template")
		}

		// Test using template
		req := &models.CreateIncidentFromTemplateRequest{
			TemplateID: template.ID,
			Variables: map[string]string{
				"service": "user-api",
				"impact":  "high",
			},
		}

		newIncident, err := incidentService.UseTemplate(req, "test-user-1")
		if err != nil {
			t.Fatalf("Failed to use template: %v", err)
		}

		expectedTitle := "user-api Service Issue"
		if newIncident.Title != expectedTitle {
			t.Errorf("Expected title '%s', got '%s'", expectedTitle, newIncident.Title)
		}
	})

	t.Run("Search", func(t *testing.T) {
		// Test search functionality
		searchReq := &models.IncidentSearchRequest{
			Query:    "Test",
			Status:   []models.IncidentStatus{models.IncidentStatusOpen},
			Severity: []models.IncidentSeverity{models.SeverityHigh},
			Page:     1,
			Limit:    10,
		}

		searchResp, err := incidentService.SearchIncidents(searchReq)
		if err != nil {
			t.Fatalf("Failed to search incidents: %v", err)
		}

		if searchResp.Total == 0 {
			t.Error("Expected at least 1 search result")
		}
	})

	t.Run("BulkOperations", func(t *testing.T) {
		// Create another incident for bulk operations
		incident2, err := incidentService.CreateIncident(
			"Second Test Incident",
			"Second Test Description",
			models.SeverityMedium,
			[]string{},
		)
		if err != nil {
			t.Fatalf("Failed to create second incident: %v", err)
		}

		// Test bulk acknowledge
		response, err := incidentService.BulkAcknowledge(
			[]string{incident.ID, incident2.ID},
			"test-assignee",
			"test-user-1",
		)
		if err != nil {
			t.Fatalf("Failed to bulk acknowledge: %v", err)
		}

		if response.ProcessedCount != 2 {
			t.Errorf("Expected 2 processed incidents, got %d", response.ProcessedCount)
		}

		// Test bulk status update
		response, err = incidentService.BulkUpdateStatus(
			[]string{incident.ID, incident2.ID},
			models.IncidentStatusResolved,
			"test-user-1",
		)
		if err != nil {
			t.Fatalf("Failed to bulk update status: %v", err)
		}

		if response.ProcessedCount != 2 {
			t.Errorf("Expected 2 processed incidents, got %d", response.ProcessedCount)
		}
	})

	t.Run("Assignment", func(t *testing.T) {
		// Create another incident for assignment testing
		incident3, err := incidentService.CreateIncident(
			"Assignment Test Incident",
			"Assignment Test Description",
			models.SeverityLow,
			[]string{},
		)
		if err != nil {
			t.Fatalf("Failed to create incident for assignment test: %v", err)
		}

		// Test assignment
		err = incidentService.AssignIncident(incident3.ID, "test-assignee", "test-user-1")
		if err != nil {
			t.Fatalf("Failed to assign incident: %v", err)
		}

		// Verify assignment
		assignedIncident, err := store.GetIncident(incident3.ID)
		if err != nil {
			t.Fatalf("Failed to get assigned incident: %v", err)
		}

		if assignedIncident.AssigneeID != "test-assignee" {
			t.Errorf("Expected assignee 'test-assignee', got '%s'", assignedIncident.AssigneeID)
		}

		// Test reassignment
		err = incidentService.ReassignIncident(incident3.ID, "new-assignee", "test-user-1")
		if err != nil {
			t.Fatalf("Failed to reassign incident: %v", err)
		}

		// Verify reassignment
		reassignedIncident, err := store.GetIncident(incident3.ID)
		if err != nil {
			t.Fatalf("Failed to get reassigned incident: %v", err)
		}

		if reassignedIncident.AssigneeID != "new-assignee" {
			t.Errorf("Expected assignee 'new-assignee', got '%s'", reassignedIncident.AssigneeID)
		}
	})
}

func TestMemoryStoreEnhancedFeatures(t *testing.T) {
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}
	defer store.Close()

	// Create test incident
	incident := &models.Incident{
		ID:          "test-incident-1",
		Title:       "Test Incident",
		Description: "Test Description",
		Status:      models.IncidentStatusOpen,
		Severity:    models.SeverityHigh,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Labels:      make(map[string]string),
	}

	err = store.CreateIncident(incident)
	if err != nil {
		t.Fatalf("Failed to create test incident: %v", err)
	}

	t.Run("CommentStorage", func(t *testing.T) {
		comment := &models.IncidentComment{
			ID:          "comment-1",
			IncidentID:  incident.ID,
			Content:     "Test comment",
			CommentType: models.CommentTypeComment,
			CreatedAt:   time.Now(),
		}

		err := store.CreateIncidentComment(comment)
		if err != nil {
			t.Fatalf("Failed to create comment: %v", err)
		}

		comments, err := store.GetIncidentComments(incident.ID)
		if err != nil {
			t.Fatalf("Failed to get comments: %v", err)
		}

		if len(comments) != 1 {
			t.Errorf("Expected 1 comment, got %d", len(comments))
		}
	})

	t.Run("TagStorage", func(t *testing.T) {
		tag := &models.IncidentTag{
			ID:         "tag-1",
			IncidentID: incident.ID,
			TagName:    "test",
			TagValue:   strPtr("value"),
			Color:      "#ff0000",
			CreatedAt:  time.Now(),
		}

		err := store.CreateIncidentTag(tag)
		if err != nil {
			t.Fatalf("Failed to create tag: %v", err)
		}

		tags, err := store.GetIncidentTags(incident.ID)
		if err != nil {
			t.Fatalf("Failed to get tags: %v", err)
		}

		if len(tags) != 1 {
			t.Errorf("Expected 1 tag, got %d", len(tags))
		}

		err = store.DeleteIncidentTag(incident.ID, "test")
		if err != nil {
			t.Fatalf("Failed to delete tag: %v", err)
		}

		remainingTags, err := store.GetIncidentTags(incident.ID)
		if err != nil {
			t.Fatalf("Failed to get remaining tags: %v", err)
		}

		if len(remainingTags) != 0 {
			t.Errorf("Expected 0 remaining tags, got %d", len(remainingTags))
		}
	})

	t.Run("TemplateStorage", func(t *testing.T) {
		template := &models.IncidentTemplate{
			ID:                  "template-1",
			Name:                "Test Template",
			TitleTemplate:       "Test {{service}}",
			DescriptionTemplate: "Issue with {{service}}",
			Severity:            models.SeverityMedium,
			IsActive:            true,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		err := store.CreateIncidentTemplate(template)
		if err != nil {
			t.Fatalf("Failed to create template: %v", err)
		}

		retrievedTemplate, err := store.GetIncidentTemplate(template.ID)
		if err != nil {
			t.Fatalf("Failed to get template: %v", err)
		}

		if retrievedTemplate.Name != template.Name {
			t.Errorf("Expected name '%s', got '%s'", template.Name, retrievedTemplate.Name)
		}

		templates, err := store.ListIncidentTemplates(true)
		if err != nil {
			t.Fatalf("Failed to list templates: %v", err)
		}

		if len(templates) != 1 {
			t.Errorf("Expected 1 template, got %d", len(templates))
		}
	})
}

// Helper function to create string pointer
func strPtr(s string) *string {
	return &s
}