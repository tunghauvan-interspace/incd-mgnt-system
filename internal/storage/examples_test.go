package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

// TestRepositoryUsage demonstrates how to use the repository pattern
func TestRepositoryUsage(t *testing.T) {
	// This test demonstrates the expected usage of the repository pattern
	// It shows how filtering, pagination, and context would be used in practice

	t.Run("IncidentRepositoryUsage", func(t *testing.T) {
		// In a real application, you would create a PostgresStore connected to a database
		// Here we demonstrate the interface and filtering capabilities
		
		ctx := context.Background()
		
		// Example of creating filters for different use cases
		
		// Filter for open high-severity incidents
		openHighSeverityFilter := storage.IncidentFilter{
			Status:   &[]models.IncidentStatus{models.IncidentStatusOpen}[0],
			Severity: &[]models.IncidentSeverity{models.SeverityHigh}[0],
			Limit:    10,
			Offset:   0,
			OrderBy:  "created_at",
		}

		// Filter for incidents assigned to a specific user
		assigneeID := "user-123"
		userIncidentsFilter := storage.IncidentFilter{
			AssigneeID: &assigneeID,
			Limit:      20,
			OrderBy:    "updated_at",
		}

		// Pagination filter for incident listing
		paginationFilter := storage.IncidentFilter{
			Limit:   50,
			Offset:  100, // Third page with 50 items per page
			OrderBy: "created_at",
		}

		// Verify filters are structured correctly
		if openHighSeverityFilter.Status == nil || *openHighSeverityFilter.Status != models.IncidentStatusOpen {
			t.Error("Open filter not set correctly")
		}
		if userIncidentsFilter.AssigneeID == nil || *userIncidentsFilter.AssigneeID != "user-123" {
			t.Error("Assignee filter not set correctly")
		}
		if paginationFilter.Limit != 50 || paginationFilter.Offset != 100 {
			t.Error("Pagination filter not set correctly")
		}

		// In a real implementation, you would use these filters like:
		// incidents, err := repository.ListIncidents(ctx, openHighSeverityFilter)
		// count, err := repository.CountIncidents(ctx, openHighSeverityFilter)

		t.Logf("Filters created successfully for context: %v", ctx)
	})

	t.Run("AlertRepositoryUsage", func(t *testing.T) {
		ctx := context.Background()
		
		// Example alert filters
		
		// Filter for firing alerts
		status := "firing"
		firingAlertsFilter := storage.AlertFilter{
			Status:  &status,
			Limit:   25,
			OrderBy: "starts_at",
		}

		// Filter for alerts associated with an incident
		incidentID := "incident-456"
		incidentAlertsFilter := storage.AlertFilter{
			IncidentID: &incidentID,
			OrderBy:    "starts_at",
		}

		// Filter by fingerprint
		fingerprint := "specific-alert-fingerprint"
		specificAlertFilter := storage.AlertFilter{
			Fingerprint: &fingerprint,
		}

		// Verify filters
		if firingAlertsFilter.Status == nil || *firingAlertsFilter.Status != "firing" {
			t.Error("Firing status filter not set correctly")
		}
		if incidentAlertsFilter.IncidentID == nil || *incidentAlertsFilter.IncidentID != "incident-456" {
			t.Error("Incident ID filter not set correctly")
		}
		if specificAlertFilter.Fingerprint == nil || *specificAlertFilter.Fingerprint != "specific-alert-fingerprint" {
			t.Error("Fingerprint filter not set correctly")
		}

		t.Logf("Alert filters created successfully for context: %v", ctx)
	})

	t.Run("TransactionUsageExample", func(t *testing.T) {
		// This demonstrates how transactions would be used in practice
		ctx := context.Background()
		
		// Example of a complex operation that needs to be atomic
		// In practice, this would use a real repository implementation
		
		exampleTransactionLogic := func(repo storage.Repository) error {
			// This would be the logic inside a transaction
			
			// Create an incident
			incident := &models.Incident{
				ID:          uuid.New().String(),
				Title:       "Database Connection Lost",
				Description: "Unable to connect to primary database",
				Status:      models.IncidentStatusOpen,
				Severity:    models.SeverityCritical,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Labels:      map[string]string{"service": "database", "env": "production"},
			}
			
			if err := repo.CreateIncident(ctx, incident); err != nil {
				return err
			}
			
			// Create related alerts
			alert1 := &models.Alert{
				ID:          uuid.New().String(),
				Fingerprint: "db-connection-alert-1",
				Status:      "firing",
				StartsAt:    time.Now(),
				Labels:      map[string]string{"alertname": "DatabaseDown", "instance": "db-1"},
				Annotations: map[string]string{"description": "Database instance db-1 is down"},
				IncidentID:  incident.ID,
				CreatedAt:   time.Now(),
			}

			alert2 := &models.Alert{
				ID:          uuid.New().String(),
				Fingerprint: "db-connection-alert-2", 
				Status:      "firing",
				StartsAt:    time.Now(),
				Labels:      map[string]string{"alertname": "DatabaseDown", "instance": "db-2"},
				Annotations: map[string]string{"description": "Database instance db-2 is down"},
				IncidentID:  incident.ID,
				CreatedAt:   time.Now(),
			}

			if err := repo.CreateAlert(ctx, alert1); err != nil {
				return err
			}
			
			if err := repo.CreateAlert(ctx, alert2); err != nil {
				return err
			}
			
			// Update incident with alert references
			incident.AlertIDs = []string{alert1.ID, alert2.ID}
			return repo.UpdateIncident(ctx, incident)
		}

		// In a real application, this would be:
		// err := repository.WithTransaction(ctx, exampleTransactionLogic)
		
		// For this test, we just verify the function structure
		if exampleTransactionLogic == nil {
			t.Error("Transaction logic should be defined")
		}
		
		t.Log("Transaction usage example created successfully")
	})
}

// TestRepositoryPatternBestPractices demonstrates best practices for using the repository pattern
func TestRepositoryPatternBestPractices(t *testing.T) {
	t.Run("FilteringBestPractices", func(t *testing.T) {
		// Best practice: Use pointer fields for optional filters to distinguish
		// between "not filtered" (nil) and "filtered to empty/default value"
		
		// Good: Filter by specific status
		openStatus := models.IncidentStatusOpen
		filter1 := storage.IncidentFilter{
			Status: &openStatus,
		}
		
		// Good: No status filter (will match all statuses)
		filter2 := storage.IncidentFilter{
			Status: nil,
		}
		
		if filter1.Status == nil {
			t.Error("Filter1 should have status set")
		}
		if filter2.Status != nil {
			t.Error("Filter2 should have no status filter")
		}
		
		// Best practice: Always set reasonable limits for pagination
		productionFilter := storage.IncidentFilter{
			Limit:   100, // Prevent accidentally loading thousands of records
			OrderBy: "created_at",
		}
		
		if productionFilter.Limit == 0 {
			t.Error("Production filters should have reasonable limits")
		}
	})

	t.Run("ContextBestPractices", func(t *testing.T) {
		// Best practice: Use context with timeout for database operations
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		// Best practice: Pass context through all repository operations
		filter := storage.IncidentFilter{Limit: 10}
		
		// In real code:
		// incidents, err := repository.ListIncidents(ctx, filter)
		// count, err := repository.CountIncidents(ctx, filter)
		
		// Verify context is properly structured
		if ctx == nil {
			t.Error("Context should not be nil")
		}
		
		// Verify filter is set up correctly
		if filter.Limit != 10 {
			t.Error("Filter limit should be 10")
		}
		
		select {
		case <-ctx.Done():
			t.Error("Context should not be cancelled immediately")
		default:
			// Good - context is still active
		}
		
		t.Log("Context best practices demonstrated")
	})

	t.Run("ErrorHandlingBestPractices", func(t *testing.T) {
		// Best practice: Handle ErrNotFound appropriately
		// In real code, you would check for storage.ErrNotFound
		
		err := storage.ErrNotFound
		if err == nil {
			t.Error("ErrNotFound should be defined")
		}
		
		// Best practice: Use errors.Is() for error comparison
		// Example: if errors.Is(err, storage.ErrNotFound) { /* handle not found */ }
		
		t.Log("Error handling patterns documented")
	})
}