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

// AlertFilter defines filtering options for alert queries
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
	CreateIncident(ctx context.Context, incident *models.Incident) error
	GetIncidentByID(ctx context.Context, id string) (*models.Incident, error)
	ListIncidents(ctx context.Context, filter IncidentFilter) ([]*models.Incident, error)
	UpdateIncident(ctx context.Context, incident *models.Incident) error
	DeleteIncident(ctx context.Context, id string) error
	CountIncidents(ctx context.Context, filter IncidentFilter) (int, error)
}

// AlertRepository interface defines operations for alert management
type AlertRepository interface {
	CreateAlert(ctx context.Context, alert *models.Alert) error
	GetAlertByID(ctx context.Context, id string) (*models.Alert, error)
	ListAlerts(ctx context.Context, filter AlertFilter) ([]*models.Alert, error)
	UpdateAlert(ctx context.Context, alert *models.Alert) error
	DeleteAlert(ctx context.Context, id string) error
	CountAlerts(ctx context.Context, filter AlertFilter) (int, error)
}

// Repository combines all repository interfaces
type Repository interface {
	IncidentRepository
	AlertRepository
	
	// Transaction support
	WithTransaction(ctx context.Context, fn func(Repository) error) error
	
	// Health and stats
	HealthCheck() error
	Close() error
}