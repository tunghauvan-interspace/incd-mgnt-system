package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

// IncidentService handles incident operations
type IncidentService struct {
	store storage.Store
}

// NewIncidentService creates a new incident service
func NewIncidentService(store storage.Store) *IncidentService {
	return &IncidentService{
		store: store,
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

	err := s.store.CreateIncident(incident)
	if err != nil {
		return nil, fmt.Errorf("failed to create incident: %w", err)
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