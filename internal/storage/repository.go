package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

// IncidentFilter defines filtering options for incident queries
type IncidentFilter struct {
	Status     *models.IncidentStatus
	Severity   *models.IncidentSeverity
	AssigneeID *string
	Limit      int
	Offset     int
	OrderBy    string // "created_at", "updated_at", etc.
}

// AlertFilter defines filtering options for alert queries (for legacy Store interface)
type AlertFilter struct {
	Status      *string
	IncidentID  *string
	Fingerprint *string
	Limit       int
	Offset      int
	OrderBy     string // "created_at", "starts_at", etc.
}

// IncidentRepository interface defines operations for incident management
type IncidentRepository interface {
	Create(ctx context.Context, incident *models.Incident) error
	GetByID(ctx context.Context, id string) (*models.Incident, error)
	List(ctx context.Context, filter IncidentFilter) ([]*models.Incident, error)
	Update(ctx context.Context, incident *models.Incident) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter IncidentFilter) (int, error)
}

// NewIncidentRepository creates a new incident repository
func NewIncidentRepository(store *PostgresStore) IncidentRepository {
	return &IncidentRepositoryImpl{store: store}
}

// IncidentRepositoryImpl provides incident repository implementation
type IncidentRepositoryImpl struct {
	store *PostgresStore
}

// Create implements IncidentRepository.Create
func (r *IncidentRepositoryImpl) Create(ctx context.Context, incident *models.Incident) error {
	labelsJSON, err := json.Marshal(incident.Labels)
	if err != nil {
		return fmt.Errorf("failed to marshal labels: %w", err)
	}
	
	query := `
		INSERT INTO incidents (id, title, description, status, severity, created_at, updated_at, assignee_id, labels)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err = r.store.db.ExecContext(ctx, query,
		incident.ID, incident.Title, incident.Description, incident.Status, incident.Severity,
		incident.CreatedAt, incident.UpdatedAt, incident.AssigneeID, labelsJSON,
	)
	
	return err
}

// GetByID implements IncidentRepository.GetByID
func (r *IncidentRepositoryImpl) GetByID(ctx context.Context, id string) (*models.Incident, error) {
	return r.store.GetIncidentByID(ctx, id)
}

// List implements IncidentRepository.List
func (r *IncidentRepositoryImpl) List(ctx context.Context, filter IncidentFilter) ([]*models.Incident, error) {
	return r.store.ListIncidentsWithFilter(ctx, filter)
}

// Update implements IncidentRepository.Update
func (r *IncidentRepositoryImpl) Update(ctx context.Context, incident *models.Incident) error {
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
	
	result, err := r.store.db.ExecContext(ctx, query,
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

// Delete implements IncidentRepository.Delete
func (r *IncidentRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM incidents WHERE id = $1`
	result, err := r.store.db.ExecContext(ctx, query, id)
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

// Count implements IncidentRepository.Count
func (r *IncidentRepositoryImpl) Count(ctx context.Context, filter IncidentFilter) (int, error) {
	return r.store.CountIncidents(ctx, filter)
}
