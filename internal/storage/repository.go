package storage

import (
	"context"

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
	return r.store.CreateIncidentWithContext(ctx, incident)
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
	return r.store.UpdateIncidentWithContext(ctx, incident)
}

// Delete implements IncidentRepository.Delete
func (r *IncidentRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.store.DeleteIncidentWithContext(ctx, id)
}

// Count implements IncidentRepository.Count
func (r *IncidentRepositoryImpl) Count(ctx context.Context, filter IncidentFilter) (int, error) {
	return r.store.CountIncidents(ctx, filter)
}
