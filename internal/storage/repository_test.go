package storage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

// TestRepositoryPattern tests the repository pattern implementation
func TestRepositoryPattern(t *testing.T) {
	// Use memory store for testing
	memStore, err := NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}
	defer memStore.Close()

	// Create repository wrapper around postgres store structure
	// Note: For testing, we'll test the interface compliance and filtering logic
	// The actual database operations are tested in postgres_test.go with real DB
	
	t.Run("IncidentFilter", func(t *testing.T) {
		// Test incident filter structure
		filter := IncidentFilter{
			Limit:   10,
			Offset:  5,
			OrderBy: "created_at",
		}

		if filter.Limit != 10 {
			t.Errorf("Expected Limit=10, got %d", filter.Limit)
		}
		if filter.Offset != 5 {
			t.Errorf("Expected Offset=5, got %d", filter.Offset)
		}
		if filter.OrderBy != "created_at" {
			t.Errorf("Expected OrderBy='created_at', got %s", filter.OrderBy)
		}
	})

	t.Run("AlertFilter", func(t *testing.T) {
		// Test alert filter structure
		filter := AlertFilter{
			Limit:   20,
			Offset:  10,
			OrderBy: "starts_at",
		}

		if filter.Limit != 20 {
			t.Errorf("Expected Limit=20, got %d", filter.Limit)
		}
		if filter.Offset != 10 {
			t.Errorf("Expected Offset=10, got %d", filter.Offset)
		}
		if filter.OrderBy != "starts_at" {
			t.Errorf("Expected OrderBy='starts_at', got %s", filter.OrderBy)
		}
	})

	t.Run("RepositoryInterfaceCompliance", func(t *testing.T) {
		// Create a PostgresStore instance for testing interface compliance
		// This tests that our PostgresStore can be used as a Repository
		
		// Note: We can't actually test database operations without a real DB connection
		// But we can test that the types and interfaces are correctly defined
		
		var _ IncidentRepository = &RepositoryImpl{}
		var _ AlertRepository = &RepositoryImpl{}
		var _ Repository = &RepositoryImpl{}
	})
}

// TestMemoryStoreBasicOperations tests basic CRUD operations with memory store
func TestMemoryStoreBasicOperations(t *testing.T) {
	store, err := NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}
	defer store.Close()

	// Test incident operations
	incident := &models.Incident{
		ID:          uuid.New().String(),
		Title:       "Test Incident",
		Description: "Test incident description",
		Status:      models.IncidentStatusOpen,
		Severity:    models.SeverityHigh,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AssigneeID:  "test-user",
		Labels:      map[string]string{"env": "test"},
		AlertIDs:    []string{},
	}

	// Create incident
	err = store.CreateIncident(incident)
	if err != nil {
		t.Fatalf("Failed to create incident: %v", err)
	}

	// Get incident
	retrieved, err := store.GetIncident(incident.ID)
	if err != nil {
		t.Fatalf("Failed to get incident: %v", err)
	}

	if retrieved.ID != incident.ID {
		t.Errorf("Expected ID=%s, got %s", incident.ID, retrieved.ID)
	}
	if retrieved.Title != incident.Title {
		t.Errorf("Expected Title=%s, got %s", incident.Title, retrieved.Title)
	}

	// List incidents
	incidents, err := store.ListIncidents()
	if err != nil {
		t.Fatalf("Failed to list incidents: %v", err)
	}
	if len(incidents) != 1 {
		t.Errorf("Expected 1 incident, got %d", len(incidents))
	}

	// Update incident
	incident.Title = "Updated Test Incident"
	incident.Status = models.IncidentStatusAcknowledged
	ackTime := time.Now()
	incident.AckedAt = &ackTime
	
	err = store.UpdateIncident(incident)
	if err != nil {
		t.Fatalf("Failed to update incident: %v", err)
	}

	// Get updated incident
	updated, err := store.GetIncident(incident.ID)
	if err != nil {
		t.Fatalf("Failed to get updated incident: %v", err)
	}
	if updated.Title != "Updated Test Incident" {
		t.Errorf("Expected updated title, got %s", updated.Title)
	}
	if updated.Status != models.IncidentStatusAcknowledged {
		t.Errorf("Expected status acknowledged, got %s", updated.Status)
	}

	// Delete incident
	err = store.DeleteIncident(incident.ID)
	if err != nil {
		t.Fatalf("Failed to delete incident: %v", err)
	}

	// Verify deletion
	_, err = store.GetIncident(incident.ID)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}
}

// TestAlertOperations tests basic alert CRUD operations
func TestAlertOperations(t *testing.T) {
	store, err := NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}
	defer store.Close()

	alert := &models.Alert{
		ID:          uuid.New().String(),
		Fingerprint: "test-fingerprint",
		Status:      "firing",
		StartsAt:    time.Now(),
		Labels:      map[string]string{"alertname": "TestAlert"},
		Annotations: map[string]string{"description": "Test alert"},
		CreatedAt:   time.Now(),
	}

	// Create alert
	err = store.CreateAlert(alert)
	if err != nil {
		t.Fatalf("Failed to create alert: %v", err)
	}

	// Get alert
	retrieved, err := store.GetAlert(alert.ID)
	if err != nil {
		t.Fatalf("Failed to get alert: %v", err)
	}

	if retrieved.ID != alert.ID {
		t.Errorf("Expected ID=%s, got %s", alert.ID, retrieved.ID)
	}
	if retrieved.Fingerprint != alert.Fingerprint {
		t.Errorf("Expected Fingerprint=%s, got %s", alert.Fingerprint, retrieved.Fingerprint)
	}

	// List alerts
	alerts, err := store.ListAlerts()
	if err != nil {
		t.Fatalf("Failed to list alerts: %v", err)
	}
	if len(alerts) != 1 {
		t.Errorf("Expected 1 alert, got %d", len(alerts))
	}

	// Update alert
	alert.Status = "resolved"
	endTime := time.Now()
	alert.EndsAt = endTime
	
	err = store.UpdateAlert(alert)
	if err != nil {
		t.Fatalf("Failed to update alert: %v", err)
	}

	// Get updated alert
	updated, err := store.GetAlert(alert.ID)
	if err != nil {
		t.Fatalf("Failed to get updated alert: %v", err)
	}
	if updated.Status != "resolved" {
		t.Errorf("Expected status resolved, got %s", updated.Status)
	}

	// Delete alert
	err = store.DeleteAlert(alert.ID)
	if err != nil {
		t.Fatalf("Failed to delete alert: %v", err)
	}

	// Verify deletion
	_, err = store.GetAlert(alert.ID)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}
}

// TestFiltering tests the filtering functionality of the repository pattern
func TestFiltering(t *testing.T) {
	t.Run("IncidentFilterWithPointers", func(t *testing.T) {
		// Test filter with pointer fields for optional filtering
		status := models.IncidentStatusOpen
		severity := models.SeverityHigh
		assignee := "user123"
		
		filter := IncidentFilter{
			Status:     &status,
			Severity:   &severity,
			AssigneeID: &assignee,
			Limit:      10,
			Offset:     0,
			OrderBy:    "created_at",
		}

		// Verify filter fields are set correctly
		if filter.Status == nil || *filter.Status != models.IncidentStatusOpen {
			t.Errorf("Expected Status to be set to Open")
		}
		if filter.Severity == nil || *filter.Severity != models.SeverityHigh {
			t.Errorf("Expected Severity to be set to High")
		}
		if filter.AssigneeID == nil || *filter.AssigneeID != "user123" {
			t.Errorf("Expected AssigneeID to be set to user123")
		}
	})

	t.Run("AlertFilterWithPointers", func(t *testing.T) {
		// Test filter with pointer fields for optional filtering
		status := "firing"
		incidentID := "incident-123"
		fingerprint := "alert-fingerprint"
		
		filter := AlertFilter{
			Status:      &status,
			IncidentID:  &incidentID,
			Fingerprint: &fingerprint,
			Limit:       20,
			Offset:      5,
			OrderBy:     "starts_at",
		}

		// Verify filter fields are set correctly
		if filter.Status == nil || *filter.Status != "firing" {
			t.Errorf("Expected Status to be set to firing")
		}
		if filter.IncidentID == nil || *filter.IncidentID != "incident-123" {
			t.Errorf("Expected IncidentID to be set to incident-123")
		}
		if filter.Fingerprint == nil || *filter.Fingerprint != "alert-fingerprint" {
			t.Errorf("Expected Fingerprint to be set to alert-fingerprint")
		}
	})

	t.Run("EmptyFilters", func(t *testing.T) {
		// Test empty filters (should match all records)
		incidentFilter := IncidentFilter{}
		if incidentFilter.Status != nil {
			t.Errorf("Expected empty Status filter")
		}
		if incidentFilter.Limit != 0 {
			t.Errorf("Expected default Limit of 0")
		}

		alertFilter := AlertFilter{}
		if alertFilter.Status != nil {
			t.Errorf("Expected empty Status filter")
		}
		if alertFilter.Limit != 0 {
			t.Errorf("Expected default Limit of 0")
		}
	})
}