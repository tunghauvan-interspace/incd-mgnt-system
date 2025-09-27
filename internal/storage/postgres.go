package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/config"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

// PostgresStore implements the Store interface using PostgreSQL
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgreSQL store with connection pooling
func NewPostgresStore(cfg *config.Config) (*PostgresStore, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &PostgresStore{db: db}

	// Run migrations
	if err := store.runMigrations(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return store, nil
}

// Close closes the database connection
func (s *PostgresStore) Close() error {
	return s.db.Close()
}

// runMigrations runs database migrations
func (s *PostgresStore) runMigrations() error {
	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Try different migration paths to handle different execution contexts
	migrationPaths := []string{
		"file://migrations",
		"file://../../migrations",
		"file://" + getMigrationsPath(),
	}

	var m *migrate.Migrate
	for _, path := range migrationPaths {
		m, err = migrate.NewWithDatabaseInstance(path, "postgres", driver)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// getMigrationsPath returns the absolute path to migrations directory
func getMigrationsPath() string {
	// Try to find migrations directory relative to current working directory
	paths := []string{
		"migrations",
		"../migrations",
		"../../migrations",
		"../../../migrations",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			abs, _ := filepath.Abs(path)
			return abs
		}
	}

	// Fallback to relative path
	return "migrations"
}

// Incident methods

// GetIncidentByID implements IncidentRepository.GetIncidentByID
func (s *PostgresStore) GetIncidentByID(ctx context.Context, id string) (*models.Incident, error) {
	query := `
		SELECT id, title, description, status, severity, created_at, updated_at,
		       acked_at, resolved_at, assignee_id, labels
		FROM incidents
		WHERE id = $1
	`

	var incident models.Incident
	var labelsJSON []byte

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&incident.ID, &incident.Title, &incident.Description,
		&incident.Status, &incident.Severity, &incident.CreatedAt, &incident.UpdatedAt,
		&incident.AckedAt, &incident.ResolvedAt, &incident.AssigneeID, &labelsJSON,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Parse labels JSON
	if len(labelsJSON) > 0 {
		if err := json.Unmarshal(labelsJSON, &incident.Labels); err != nil {
			return nil, fmt.Errorf("failed to parse labels: %w", err)
		}
	} else {
		incident.Labels = make(map[string]string)
	}

	// Get associated alert IDs
	alertQuery := `SELECT id FROM alerts WHERE incident_id = $1`
	rows, err := s.db.QueryContext(ctx, alertQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alertIDs []string
	for rows.Next() {
		var alertID string
		if err := rows.Scan(&alertID); err != nil {
			return nil, err
		}
		alertIDs = append(alertIDs, alertID)
	}
	incident.AlertIDs = alertIDs

	return &incident, nil
}

// GetIncident provides backward compatibility
func (s *PostgresStore) GetIncident(id string) (*models.Incident, error) {
	return s.GetIncidentByID(context.Background(), id)
}

// ListIncidents implements IncidentRepository.ListIncidents with filtering and pagination
func (s *PostgresStore) ListIncidentsWithFilter(ctx context.Context, filter IncidentFilter) ([]*models.Incident, error) {
	// Build query with filtering
	query := `
		SELECT id, title, description, status, severity, created_at, updated_at,
		       acked_at, resolved_at, assignee_id, labels
		FROM incidents
		WHERE ($1::incident_status IS NULL OR status = $1)
		  AND ($2::incident_severity IS NULL OR severity = $2)
		  AND ($3::text IS NULL OR assignee_id = $3)
	`

	// Handle ordering with SQL injection protection
	orderBy := "created_at"
	if filter.OrderBy != "" {
		switch filter.OrderBy {
		case "created_at", "updated_at", "title", "status", "severity":
			orderBy = filter.OrderBy
		default:
			orderBy = "created_at"
		}
	}
	query += " ORDER BY " + orderBy + " DESC"

	// Add pagination
	if filter.Limit > 0 {
		query += " LIMIT $4"
		if filter.Offset > 0 {
			query += " OFFSET $5"
		}
	}

	// Execute query with appropriate parameters
	var rows *sql.Rows
	var err error

	if filter.Limit > 0 {
		if filter.Offset > 0 {
			rows, err = s.db.QueryContext(ctx, query, filter.Status, filter.Severity, filter.AssigneeID, filter.Limit, filter.Offset)
		} else {
			rows, err = s.db.QueryContext(ctx, query, filter.Status, filter.Severity, filter.AssigneeID, filter.Limit)
		}
	} else {
		// Remove LIMIT clause if no limit specified
		query = `
			SELECT id, title, description, status, severity, created_at, updated_at,
			       acked_at, resolved_at, assignee_id, labels
			FROM incidents
			WHERE ($1::incident_status IS NULL OR status = $1)
			  AND ($2::incident_severity IS NULL OR severity = $2)
			  AND ($3::text IS NULL OR assignee_id = $3)
			ORDER BY ` + orderBy + ` DESC`
		rows, err = s.db.QueryContext(ctx, query, filter.Status, filter.Severity, filter.AssigneeID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []*models.Incident
	for rows.Next() {
		var incident models.Incident
		var labelsJSON []byte

		err := rows.Scan(
			&incident.ID, &incident.Title, &incident.Description,
			&incident.Status, &incident.Severity, &incident.CreatedAt, &incident.UpdatedAt,
			&incident.AckedAt, &incident.ResolvedAt, &incident.AssigneeID, &labelsJSON,
		)
		if err != nil {
			return nil, err
		}

		// Parse labels JSON
		if len(labelsJSON) > 0 {
			if err := json.Unmarshal(labelsJSON, &incident.Labels); err != nil {
				return nil, fmt.Errorf("failed to parse labels: %w", err)
			}
		} else {
			incident.Labels = make(map[string]string)
		}

		// Get associated alert IDs for each incident
		alertQuery := `SELECT id FROM alerts WHERE incident_id = $1`
		alertRows, err := s.db.QueryContext(ctx, alertQuery, incident.ID)
		if err != nil {
			return nil, err
		}

		var alertIDs []string
		for alertRows.Next() {
			var alertID string
			if err := alertRows.Scan(&alertID); err != nil {
				alertRows.Close()
				return nil, err
			}
			alertIDs = append(alertIDs, alertID)
		}
		alertRows.Close()
		incident.AlertIDs = alertIDs

		incidents = append(incidents, &incident)
	}

	return incidents, nil
}

// ListIncidents provides backward compatibility for the old Store interface
func (s *PostgresStore) ListIncidents() ([]*models.Incident, error) {
	ctx := context.Background()
	query := `
		SELECT id, title, description, status, severity, created_at, updated_at,
		       acked_at, resolved_at, assignee_id, labels
		FROM incidents
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []*models.Incident
	for rows.Next() {
		var incident models.Incident
		var labelsJSON []byte

		err := rows.Scan(
			&incident.ID, &incident.Title, &incident.Description,
			&incident.Status, &incident.Severity, &incident.CreatedAt, &incident.UpdatedAt,
			&incident.AckedAt, &incident.ResolvedAt, &incident.AssigneeID, &labelsJSON,
		)
		if err != nil {
			return nil, err
		}

		// Parse labels JSON
		if len(labelsJSON) > 0 {
			if err := json.Unmarshal(labelsJSON, &incident.Labels); err != nil {
				return nil, fmt.Errorf("failed to parse labels: %w", err)
			}
		} else {
			incident.Labels = make(map[string]string)
		}

		// Get associated alert IDs for each incident
		alertQuery := `SELECT id FROM alerts WHERE incident_id = $1`
		alertRows, err := s.db.QueryContext(ctx, alertQuery, incident.ID)
		if err != nil {
			return nil, err
		}

		var alertIDs []string
		for alertRows.Next() {
			var alertID string
			if err := alertRows.Scan(&alertID); err != nil {
				alertRows.Close()
				return nil, err
			}
			alertIDs = append(alertIDs, alertID)
		}
		alertRows.Close()
		incident.AlertIDs = alertIDs

		incidents = append(incidents, &incident)
	}

	return incidents, nil
}

// CreateIncident implements IncidentRepository.CreateIncident
func (s *PostgresStore) CreateIncidentWithContext(ctx context.Context, incident *models.Incident) error {
	labelsJSON, err := json.Marshal(incident.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	query := `
		INSERT INTO incidents (id, title, description, status, severity, created_at, updated_at, assignee_id, labels)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = s.db.ExecContext(ctx, query,
		incident.ID, incident.Title, incident.Description, incident.Status, incident.Severity,
		incident.CreatedAt, incident.UpdatedAt, incident.AssigneeID, labelsJSON,
	)

	return err
}

// Old interface backward compatibility methods

// CreateIncident (old interface signature)
// CreateIncident provides backward compatibility for the old Store interface
func (s *PostgresStore) CreateIncident(incident *models.Incident) error {
	return s.CreateIncidentWithContext(context.Background(), incident)
}

// UpdateIncident implements IncidentRepository.UpdateIncident
func (s *PostgresStore) UpdateIncidentWithContext(ctx context.Context, incident *models.Incident) error {
	labelsJSON, err := json.Marshal(incident.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	query := `
		UPDATE incidents 
		SET title = $2, description = $3, status = $4, severity = $5,
		    updated_at = $6, acked_at = $7, resolved_at = $8, assignee_id = $9, labels = $10
		WHERE id = $1
	`

	result, err := s.db.ExecContext(ctx, query,
		incident.ID, incident.Title, incident.Description, incident.Status, incident.Severity,
		incident.UpdatedAt, incident.AckedAt, incident.ResolvedAt, incident.AssigneeID, labelsJSON,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateIncident (old interface signature - backward compatibility)
// UpdateIncident provides backward compatibility for the old Store interface
func (s *PostgresStore) UpdateIncident(incident *models.Incident) error {
	return s.UpdateIncidentWithContext(context.Background(), incident)
}

// DeleteIncident implements IncidentRepository.DeleteIncident
func (s *PostgresStore) DeleteIncidentWithContext(ctx context.Context, id string) error {
	query := `DELETE FROM incidents WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// DeleteIncident (old interface signature - backward compatibility)
// DeleteIncident provides backward compatibility for the old Store interface
func (s *PostgresStore) DeleteIncident(id string) error {
	return s.DeleteIncidentWithContext(context.Background(), id)
}

// CountIncidents implements IncidentRepository.CountIncidents
func (s *PostgresStore) CountIncidents(ctx context.Context, filter IncidentFilter) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM incidents
		WHERE ($1::incident_status IS NULL OR status = $1)
		  AND ($2::incident_severity IS NULL OR severity = $2)
		  AND ($3::text IS NULL OR assignee_id = $3)
	`

	var count int
	err := s.db.QueryRowContext(ctx, query, filter.Status, filter.Severity, filter.AssigneeID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Alert methods

// GetByID implements AlertRepository.GetByID for alerts
func (s *PostgresStore) GetAlertByID(ctx context.Context, id string) (*models.Alert, error) {
	query := `
		SELECT id, fingerprint, status, starts_at, ends_at, labels, annotations, incident_id, created_at
		FROM alerts
		WHERE id = $1
	`

	var alert models.Alert
	var labelsJSON, annotationsJSON []byte
	var incidentID sql.NullString

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&alert.ID, &alert.Fingerprint, &alert.Status, &alert.StartsAt, &alert.EndsAt,
		&labelsJSON, &annotationsJSON, &incidentID, &alert.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Handle nullable incident_id
	if incidentID.Valid {
		alert.IncidentID = incidentID.String
	}

	// Parse JSON fields
	if len(labelsJSON) > 0 {
		if err := json.Unmarshal(labelsJSON, &alert.Labels); err != nil {
			return nil, fmt.Errorf("failed to parse labels: %w", err)
		}
	} else {
		alert.Labels = make(map[string]string)
	}

	if len(annotationsJSON) > 0 {
		if err := json.Unmarshal(annotationsJSON, &alert.Annotations); err != nil {
			return nil, fmt.Errorf("failed to parse annotations: %w", err)
		}
	} else {
		alert.Annotations = make(map[string]string)
	}

	return &alert, nil
}

// GetAlert provides backward compatibility
func (s *PostgresStore) GetAlert(id string) (*models.Alert, error) {
	return s.GetAlertByID(context.Background(), id)
}

// ListAlerts implements AlertRepository.ListAlerts with filtering and pagination
func (s *PostgresStore) ListAlertsWithFilter(ctx context.Context, filter AlertFilter) ([]*models.Alert, error) {
	// Build query with filtering
	query := `
		SELECT id, fingerprint, status, starts_at, ends_at, labels, annotations, incident_id, created_at
		FROM alerts
		WHERE ($1::text IS NULL OR status = $1)
		  AND ($2::uuid IS NULL OR incident_id = $2::uuid)
		  AND ($3::text IS NULL OR fingerprint = $3)
	`

	// Handle ordering with SQL injection protection
	orderBy := "created_at"
	if filter.OrderBy != "" {
		switch filter.OrderBy {
		case "created_at", "starts_at", "status", "fingerprint":
			orderBy = filter.OrderBy
		default:
			orderBy = "created_at"
		}
	}
	query += " ORDER BY " + orderBy + " DESC"

	// Add pagination
	if filter.Limit > 0 {
		query += " LIMIT $4"
		if filter.Offset > 0 {
			query += " OFFSET $5"
		}
	}

	// Execute query with appropriate parameters
	var rows *sql.Rows
	var err error

	if filter.Limit > 0 {
		if filter.Offset > 0 {
			rows, err = s.db.QueryContext(ctx, query, filter.Status, filter.IncidentID, filter.Fingerprint, filter.Limit, filter.Offset)
		} else {
			rows, err = s.db.QueryContext(ctx, query, filter.Status, filter.IncidentID, filter.Fingerprint, filter.Limit)
		}
	} else {
		// Remove LIMIT clause if no limit specified
		query = `
			SELECT id, fingerprint, status, starts_at, ends_at, labels, annotations, incident_id, created_at
			FROM alerts
			WHERE ($1::text IS NULL OR status = $1)
			  AND ($2::uuid IS NULL OR incident_id = $2::uuid)
			  AND ($3::text IS NULL OR fingerprint = $3)
			ORDER BY ` + orderBy + ` DESC`
		rows, err = s.db.QueryContext(ctx, query, filter.Status, filter.IncidentID, filter.Fingerprint)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.Alert
	for rows.Next() {
		var alert models.Alert
		var labelsJSON, annotationsJSON []byte
		var incidentID sql.NullString

		err := rows.Scan(
			&alert.ID, &alert.Fingerprint, &alert.Status, &alert.StartsAt, &alert.EndsAt,
			&labelsJSON, &annotationsJSON, &incidentID, &alert.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable incident_id
		if incidentID.Valid {
			alert.IncidentID = incidentID.String
		}

		// Parse JSON fields
		if len(labelsJSON) > 0 {
			if err := json.Unmarshal(labelsJSON, &alert.Labels); err != nil {
				return nil, fmt.Errorf("failed to parse labels: %w", err)
			}
		} else {
			alert.Labels = make(map[string]string)
		}

		if len(annotationsJSON) > 0 {
			if err := json.Unmarshal(annotationsJSON, &alert.Annotations); err != nil {
				return nil, fmt.Errorf("failed to parse annotations: %w", err)
			}
		} else {
			alert.Annotations = make(map[string]string)
		}

		alerts = append(alerts, &alert)
	}

	return alerts, nil
}

// ListAlerts (old interface signature - backward compatibility)
func (s *PostgresStore) ListAlerts() ([]*models.Alert, error) {
	ctx := context.Background()
	query := `
		SELECT id, fingerprint, status, starts_at, ends_at, labels, annotations, incident_id, created_at
		FROM alerts
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.Alert
	for rows.Next() {
		var alert models.Alert
		var labelsJSON, annotationsJSON []byte
		var incidentID sql.NullString

		err := rows.Scan(
			&alert.ID, &alert.Fingerprint, &alert.Status, &alert.StartsAt, &alert.EndsAt,
			&labelsJSON, &annotationsJSON, &incidentID, &alert.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable incident_id
		if incidentID.Valid {
			alert.IncidentID = incidentID.String
		}

		// Parse JSON fields
		if len(labelsJSON) > 0 {
			if err := json.Unmarshal(labelsJSON, &alert.Labels); err != nil {
				return nil, fmt.Errorf("failed to parse labels: %w", err)
			}
		} else {
			alert.Labels = make(map[string]string)
		}

		if len(annotationsJSON) > 0 {
			if err := json.Unmarshal(annotationsJSON, &alert.Annotations); err != nil {
				return nil, fmt.Errorf("failed to parse annotations: %w", err)
			}
		} else {
			alert.Annotations = make(map[string]string)
		}

		alerts = append(alerts, &alert)
	}

	return alerts, nil
}

// CreateAlert implements AlertRepository.CreateAlert
func (s *PostgresStore) CreateAlertWithContext(ctx context.Context, alert *models.Alert) error {
	labelsJSON, err := json.Marshal(alert.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	annotationsJSON, err := json.Marshal(alert.Annotations)
	if err != nil {
		return fmt.Errorf("failed to marshal annotations: %w", err)
	}

	query := `
		INSERT INTO alerts (id, fingerprint, status, starts_at, ends_at, labels, annotations, incident_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	// Handle empty incident_id as NULL
	var incidentID interface{}
	if alert.IncidentID == "" {
		incidentID = nil
	} else {
		incidentID = alert.IncidentID
	}

	_, err = s.db.ExecContext(ctx, query,
		alert.ID, alert.Fingerprint, alert.Status, alert.StartsAt, alert.EndsAt,
		labelsJSON, annotationsJSON, incidentID, alert.CreatedAt,
	)

	return err
}

// UpdateAlert implements AlertRepository.UpdateAlert
func (s *PostgresStore) UpdateAlertWithContext(ctx context.Context, alert *models.Alert) error {
	labelsJSON, err := json.Marshal(alert.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	annotationsJSON, err := json.Marshal(alert.Annotations)
	if err != nil {
		return fmt.Errorf("failed to marshal annotations: %w", err)
	}

	query := `
		UPDATE alerts 
		SET fingerprint = $2, status = $3, starts_at = $4, ends_at = $5,
		    labels = $6, annotations = $7, incident_id = $8
		WHERE id = $1
	`

	// Handle empty incident_id as NULL
	var incidentID interface{}
	if alert.IncidentID == "" {
		incidentID = nil
	} else {
		incidentID = alert.IncidentID
	}

	result, err := s.db.ExecContext(ctx, query,
		alert.ID, alert.Fingerprint, alert.Status, alert.StartsAt, alert.EndsAt,
		labelsJSON, annotationsJSON, incidentID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// DeleteAlert implements AlertRepository.DeleteAlert
func (s *PostgresStore) DeleteAlertWithContext(ctx context.Context, id string) error {
	query := `DELETE FROM alerts WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// CountAlerts implements AlertRepository.CountAlerts
func (s *PostgresStore) CountAlerts(ctx context.Context, filter AlertFilter) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM alerts
		WHERE ($1::text IS NULL OR status = $1)
		  AND ($2::uuid IS NULL OR incident_id = $2::uuid)
		  AND ($3::text IS NULL OR fingerprint = $3)
	`

	var count int
	err := s.db.QueryRowContext(ctx, query, filter.Status, filter.IncidentID, filter.Fingerprint).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Old interface backward compatibility methods for alerts

func (s *PostgresStore) CreateAlert(alert *models.Alert) error {
	ctx := context.Background()

	labelsJSON, err := json.Marshal(alert.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	annotationsJSON, err := json.Marshal(alert.Annotations)
	if err != nil {
		return fmt.Errorf("failed to marshal annotations: %w", err)
	}

	query := `
		INSERT INTO alerts (id, fingerprint, status, starts_at, ends_at, labels, annotations, incident_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	// Handle empty incident_id as NULL
	var incidentID interface{}
	if alert.IncidentID == "" {
		incidentID = nil
	} else {
		incidentID = alert.IncidentID
	}

	_, err = s.db.ExecContext(ctx, query,
		alert.ID, alert.Fingerprint, alert.Status, alert.StartsAt, alert.EndsAt,
		labelsJSON, annotationsJSON, incidentID, alert.CreatedAt,
	)

	return err
}

func (s *PostgresStore) UpdateAlert(alert *models.Alert) error {
	ctx := context.Background()

	labelsJSON, err := json.Marshal(alert.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}

	annotationsJSON, err := json.Marshal(alert.Annotations)
	if err != nil {
		return fmt.Errorf("failed to marshal annotations: %w", err)
	}

	query := `
		UPDATE alerts 
		SET fingerprint = $2, status = $3, starts_at = $4, ends_at = $5,
		    labels = $6, annotations = $7, incident_id = $8
		WHERE id = $1
	`

	// Handle empty incident_id as NULL
	var incidentID interface{}
	if alert.IncidentID == "" {
		incidentID = nil
	} else {
		incidentID = alert.IncidentID
	}

	result, err := s.db.ExecContext(ctx, query,
		alert.ID, alert.Fingerprint, alert.Status, alert.StartsAt, alert.EndsAt,
		labelsJSON, annotationsJSON, incidentID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresStore) DeleteAlert(id string) error {
	ctx := context.Background()

	query := `DELETE FROM alerts WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Placeholder implementations for other entity types (to maintain interface compatibility)
// These will be implemented in subsequent phases

func (s *PostgresStore) GetNotificationChannel(id string) (*models.NotificationChannel, error) {
	return nil, fmt.Errorf("notification channels not yet implemented in postgres store")
}

func (s *PostgresStore) ListNotificationChannels() ([]*models.NotificationChannel, error) {
	return nil, fmt.Errorf("notification channels not yet implemented in postgres store")
}

func (s *PostgresStore) CreateNotificationChannel(channel *models.NotificationChannel) error {
	return fmt.Errorf("notification channels not yet implemented in postgres store")
}

func (s *PostgresStore) UpdateNotificationChannel(channel *models.NotificationChannel) error {
	return fmt.Errorf("notification channels not yet implemented in postgres store")
}

func (s *PostgresStore) DeleteNotificationChannel(id string) error {
	return fmt.Errorf("notification channels not yet implemented in postgres store")
}

func (s *PostgresStore) GetEscalationPolicy(id string) (*models.EscalationPolicy, error) {
	return nil, fmt.Errorf("escalation policies not yet implemented in postgres store")
}

func (s *PostgresStore) ListEscalationPolicies() ([]*models.EscalationPolicy, error) {
	return nil, fmt.Errorf("escalation policies not yet implemented in postgres store")
}

func (s *PostgresStore) CreateEscalationPolicy(policy *models.EscalationPolicy) error {
	return fmt.Errorf("escalation policies not yet implemented in postgres store")
}

func (s *PostgresStore) UpdateEscalationPolicy(policy *models.EscalationPolicy) error {
	return fmt.Errorf("escalation policies not yet implemented in postgres store")
}

func (s *PostgresStore) DeleteEscalationPolicy(id string) error {
	return fmt.Errorf("escalation policies not yet implemented in postgres store")
}

func (s *PostgresStore) GetOnCallSchedule(id string) (*models.OnCallSchedule, error) {
	return nil, fmt.Errorf("on-call schedules not yet implemented in postgres store")
}

func (s *PostgresStore) ListOnCallSchedules() ([]*models.OnCallSchedule, error) {
	return nil, fmt.Errorf("on-call schedules not yet implemented in postgres store")
}

func (s *PostgresStore) CreateOnCallSchedule(schedule *models.OnCallSchedule) error {
	return fmt.Errorf("on-call schedules not yet implemented in postgres store")
}

func (s *PostgresStore) UpdateOnCallSchedule(schedule *models.OnCallSchedule) error {
	return fmt.Errorf("on-call schedules not yet implemented in postgres store")
}

func (s *PostgresStore) DeleteOnCallSchedule(id string) error {
	return fmt.Errorf("on-call schedules not yet implemented in postgres store")
}

// HealthCheck tests the database connection
func (s *PostgresStore) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.db.PingContext(ctx)
}

// GetDBStats returns database connection statistics
func (s *PostgresStore) GetDBStats() sql.DBStats {
	return s.db.Stats()
}
