package storage

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/config"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

// setupTestDB creates a test database connection for integration tests
func setupTestDB(t *testing.T) (*PostgresStore, func()) {
	t.Helper()
	
	// Use a test database URL or skip if not provided
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping PostgreSQL integration tests")
	}

	cfg := &config.Config{
		DatabaseURL:       dbURL,
		DBMaxOpenConns:    10,
		DBMaxIdleConns:    2,
		DBConnMaxLifetime: 5 * time.Minute,
	}

	store, err := NewPostgresStore(cfg)
	if err != nil {
		t.Fatalf("Failed to create test PostgreSQL store: %v", err)
	}

	// Cleanup function to clean test data and close connection
	cleanup := func() {
		// Clean up test data
		store.db.Exec("DELETE FROM alerts")
		store.db.Exec("DELETE FROM incidents") 
		store.Close()
	}

	return store, cleanup
}

// TestPostgresStore_IncidentCRUD tests complete CRUD operations for incidents
func TestPostgresStore_IncidentCRUD(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create test incident
	incident := &models.Incident{
		ID:          uuid.New().String(),
		Title:       "Test Incident",
		Description: "Test incident description",
		Status:      models.IncidentStatusOpen,
		Severity:    models.SeverityHigh,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AssigneeID:  "test-user",
		Labels:      map[string]string{"env": "test", "service": "api"},
		AlertIDs:    []string{},
	}

	// Test Create
	err := store.CreateIncident(incident)
	if err != nil {
		t.Fatalf("Failed to create incident: %v", err)
	}

	// Test GetIncident
	retrieved, err := store.GetIncident(incident.ID)
	if err != nil {
		t.Fatalf("Failed to get incident: %v", err)
	}

	// Verify retrieved incident matches created incident
	if retrieved.ID != incident.ID {
		t.Errorf("Expected ID %s, got %s", incident.ID, retrieved.ID)
	}
	if retrieved.Title != incident.Title {
		t.Errorf("Expected title %s, got %s", incident.Title, retrieved.Title)
	}
	if retrieved.Status != incident.Status {
		t.Errorf("Expected status %s, got %s", incident.Status, retrieved.Status)
	}
	if retrieved.Severity != incident.Severity {
		t.Errorf("Expected severity %s, got %s", incident.Severity, retrieved.Severity)
	}
	if len(retrieved.Labels) != 2 || retrieved.Labels["env"] != "test" {
		t.Errorf("Labels not properly stored/retrieved: %v", retrieved.Labels)
	}

	// Test Update
	retrieved.Status = models.IncidentStatusAcknowledged
	retrieved.UpdatedAt = time.Now()
	ackedTime := time.Now()
	retrieved.AckedAt = &ackedTime

	err = store.UpdateIncident(retrieved)
	if err != nil {
		t.Fatalf("Failed to update incident: %v", err)
	}

	// Verify update
	updated, err := store.GetIncident(incident.ID)
	if err != nil {
		t.Fatalf("Failed to get updated incident: %v", err)
	}
	if updated.Status != models.IncidentStatusAcknowledged {
		t.Errorf("Expected status %s, got %s", models.IncidentStatusAcknowledged, updated.Status)
	}
	if updated.AckedAt == nil {
		t.Error("Expected AckedAt to be set")
	}

	// Test ListIncidents
	incidents, err := store.ListIncidents()
	if err != nil {
		t.Fatalf("Failed to list incidents: %v", err)
	}
	if len(incidents) != 1 {
		t.Errorf("Expected 1 incident, got %d", len(incidents))
	}

	// Test Delete
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

// TestPostgresStore_AlertCRUD tests complete CRUD operations for alerts
func TestPostgresStore_AlertCRUD(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create test incident first (for foreign key)
	incident := &models.Incident{
		ID:          uuid.New().String(),
		Title:       "Test Incident for Alert",
		Description: "Test incident for alert testing",
		Status:      models.IncidentStatusOpen,
		Severity:    models.SeverityMedium,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Labels:      make(map[string]string),
		AlertIDs:    []string{},
	}
	err := store.CreateIncident(incident)
	if err != nil {
		t.Fatalf("Failed to create test incident: %v", err)
	}

	// Create test alert
	alert := &models.Alert{
		ID:          uuid.New().String(),
		Fingerprint: "test-fingerprint-" + uuid.New().String(),
		Status:      "firing",
		StartsAt:    time.Now().Add(-1 * time.Hour),
		EndsAt:      time.Now(),
		Labels:      map[string]string{"alertname": "test", "severity": "warning"},
		Annotations: map[string]string{"summary": "Test alert", "description": "Test alert description"},
		IncidentID:  incident.ID,
		CreatedAt:   time.Now(),
	}

	// Test Create
	err = store.CreateAlert(alert)
	if err != nil {
		t.Fatalf("Failed to create alert: %v", err)
	}

	// Test GetAlert
	retrieved, err := store.GetAlert(alert.ID)
	if err != nil {
		t.Fatalf("Failed to get alert: %v", err)
	}

	// Verify retrieved alert matches created alert
	if retrieved.ID != alert.ID {
		t.Errorf("Expected ID %s, got %s", alert.ID, retrieved.ID)
	}
	if retrieved.Fingerprint != alert.Fingerprint {
		t.Errorf("Expected fingerprint %s, got %s", alert.Fingerprint, retrieved.Fingerprint)
	}
	if retrieved.Status != alert.Status {
		t.Errorf("Expected status %s, got %s", alert.Status, retrieved.Status)
	}
	if retrieved.IncidentID != alert.IncidentID {
		t.Errorf("Expected incident ID %s, got %s", alert.IncidentID, retrieved.IncidentID)
	}
	if len(retrieved.Labels) != 2 || retrieved.Labels["alertname"] != "test" {
		t.Errorf("Labels not properly stored/retrieved: %v", retrieved.Labels)
	}
	if len(retrieved.Annotations) != 2 || retrieved.Annotations["summary"] != "Test alert" {
		t.Errorf("Annotations not properly stored/retrieved: %v", retrieved.Annotations)
	}

	// Test Update
	retrieved.Status = "resolved"
	err = store.UpdateAlert(retrieved)
	if err != nil {
		t.Fatalf("Failed to update alert: %v", err)
	}

	// Verify update
	updated, err := store.GetAlert(alert.ID)
	if err != nil {
		t.Fatalf("Failed to get updated alert: %v", err)
	}
	if updated.Status != "resolved" {
		t.Errorf("Expected status resolved, got %s", updated.Status)
	}

	// Test ListAlerts
	alerts, err := store.ListAlerts()
	if err != nil {
		t.Fatalf("Failed to list alerts: %v", err)
	}
	if len(alerts) != 1 {
		t.Errorf("Expected 1 alert, got %d", len(alerts))
	}

	// Test Delete
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

// TestPostgresStore_IncidentAlertRelationship tests the foreign key relationship
func TestPostgresStore_IncidentAlertRelationship(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create incident
	incident := &models.Incident{
		ID:          uuid.New().String(),
		Title:       "Parent Incident",
		Description: "Incident with alerts",
		Status:      models.IncidentStatusOpen,
		Severity:    models.SeverityCritical,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Labels:      make(map[string]string),
		AlertIDs:    []string{},
	}
	err := store.CreateIncident(incident)
	if err != nil {
		t.Fatalf("Failed to create incident: %v", err)
	}

	// Create multiple alerts for the incident
	alert1 := &models.Alert{
		ID:          uuid.New().String(),
		Fingerprint: "alert1-" + uuid.New().String(),
		Status:      "firing",
		StartsAt:    time.Now(),
		EndsAt:      time.Now().Add(time.Hour),
		Labels:      map[string]string{"instance": "server1"},
		Annotations: map[string]string{"summary": "Alert 1"},
		IncidentID:  incident.ID,
		CreatedAt:   time.Now(),
	}

	alert2 := &models.Alert{
		ID:          uuid.New().String(),
		Fingerprint: "alert2-" + uuid.New().String(),
		Status:      "firing",
		StartsAt:    time.Now(),
		EndsAt:      time.Now().Add(time.Hour),
		Labels:      map[string]string{"instance": "server2"},
		Annotations: map[string]string{"summary": "Alert 2"},
		IncidentID:  incident.ID,
		CreatedAt:   time.Now(),
	}

	err = store.CreateAlert(alert1)
	if err != nil {
		t.Fatalf("Failed to create alert1: %v", err)
	}

	err = store.CreateAlert(alert2)
	if err != nil {
		t.Fatalf("Failed to create alert2: %v", err)
	}

	// Retrieve incident and verify alert IDs are populated
	retrievedIncident, err := store.GetIncident(incident.ID)
	if err != nil {
		t.Fatalf("Failed to get incident: %v", err)
	}

	if len(retrievedIncident.AlertIDs) != 2 {
		t.Errorf("Expected 2 alert IDs, got %d", len(retrievedIncident.AlertIDs))
	}

	// Verify the alert IDs contain our created alerts
	alertIDMap := make(map[string]bool)
	for _, id := range retrievedIncident.AlertIDs {
		alertIDMap[id] = true
	}

	if !alertIDMap[alert1.ID] {
		t.Errorf("Expected alert1 ID %s in incident alert IDs", alert1.ID)
	}
	if !alertIDMap[alert2.ID] {
		t.Errorf("Expected alert2 ID %s in incident alert IDs", alert2.ID)
	}
}

// TestPostgresStore_DatabaseConstraints tests database constraints
func TestPostgresStore_DatabaseConstraints(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Test unique fingerprint constraint on alerts
	alert1 := &models.Alert{
		ID:          uuid.New().String(),
		Fingerprint: "duplicate-fingerprint",
		Status:      "firing",
		StartsAt:    time.Now(),
		EndsAt:      time.Now().Add(time.Hour),
		Labels:      make(map[string]string),
		Annotations: make(map[string]string),
		CreatedAt:   time.Now(),
	}

	alert2 := &models.Alert{
		ID:          uuid.New().String(),
		Fingerprint: "duplicate-fingerprint", // Same fingerprint
		Status:      "firing",
		StartsAt:    time.Now(),
		EndsAt:      time.Now().Add(time.Hour),
		Labels:      make(map[string]string),
		Annotations: make(map[string]string),
		CreatedAt:   time.Now(),
	}

	// First alert should succeed
	err := store.CreateAlert(alert1)
	if err != nil {
		t.Fatalf("Failed to create first alert: %v", err)
	}

	// Second alert should fail due to unique constraint
	err = store.CreateAlert(alert2)
	if err == nil {
		t.Error("Expected error due to unique fingerprint constraint, but got none")
	}
}

// TestPostgresStore_Migration tests migration functionality
func TestPostgresStore_Migration(t *testing.T) {
	// Use a separate database URL for migration testing
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping migration tests")
	}

	cfg := &config.Config{
		DatabaseURL:       dbURL,
		DBMaxOpenConns:    10,
		DBMaxIdleConns:    2,
		DBConnMaxLifetime: 5 * time.Minute,
	}

	// Open raw connection to test migration
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Create store which should run migrations
	store, err := NewPostgresStore(cfg)
	if err != nil {
		t.Fatalf("Failed to create PostgreSQL store: %v", err)
	}
	defer store.Close()

	// Verify tables exist
	tables := []string{"incidents", "alerts"}
	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT count(*) FROM information_schema.tables WHERE table_name = '%s'", table)
		err = db.QueryRow(query).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check table %s: %v", table, err)
		}
		if count != 1 {
			t.Errorf("Expected table %s to exist", table)
		}
	}

	// Verify custom types exist
	types := []string{"incident_status", "incident_severity"}
	for _, typeName := range types {
		var count int
		query := fmt.Sprintf("SELECT count(*) FROM pg_type WHERE typname = '%s'", typeName)
		err = db.QueryRow(query).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check type %s: %v", typeName, err)
		}
		if count != 1 {
			t.Errorf("Expected type %s to exist", typeName)
		}
	}
}