package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

// AlertService handles alert operations and grouping
type AlertService struct {
	store           storage.Store
	incidentService *IncidentService
}

// NewAlertService creates a new alert service
func NewAlertService(store storage.Store, incidentService *IncidentService) *AlertService {
	return &AlertService{
		store:           store,
		incidentService: incidentService,
	}
}

// AlertmanagerAlert represents an alert from Alertmanager
type AlertmanagerAlert struct {
	Fingerprint string            `json:"fingerprint"`
	Status      string            `json:"status"`
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// AlertmanagerWebhook represents the webhook payload from Alertmanager
type AlertmanagerWebhook struct {
	Version           string              `json:"version"`
	GroupKey          string              `json:"groupKey"`
	Status            string              `json:"status"`
	Receiver          string              `json:"receiver"`
	GroupLabels       map[string]string   `json:"groupLabels"`
	CommonLabels      map[string]string   `json:"commonLabels"`
	CommonAnnotations map[string]string   `json:"commonAnnotations"`
	ExternalURL       string              `json:"externalURL"`
	Alerts            []AlertmanagerAlert `json:"alerts"`
}

// ProcessAlertmanagerWebhook processes alerts from Alertmanager
func (s *AlertService) ProcessAlertmanagerWebhook(webhook *AlertmanagerWebhook) error {
	for _, amAlert := range webhook.Alerts {
		alert := &models.Alert{
			ID:          uuid.New().String(),
			Fingerprint: amAlert.Fingerprint,
			Status:      amAlert.Status,
			StartsAt:    amAlert.StartsAt,
			EndsAt:      amAlert.EndsAt,
			Labels:      amAlert.Labels,
			Annotations: amAlert.Annotations,
			CreatedAt:   time.Now(),
		}

		// Check if we already have this alert
		existingAlert, err := s.findAlertByFingerprint(amAlert.Fingerprint)
		if err != nil && err != storage.ErrNotFound {
			return fmt.Errorf("failed to check existing alert: %w", err)
		}

		if existingAlert != nil {
			// Update existing alert
			existingAlert.Status = alert.Status
			existingAlert.EndsAt = alert.EndsAt
			if err := s.store.UpdateAlert(existingAlert); err != nil {
				return fmt.Errorf("failed to update alert: %w", err)
			}
			alert = existingAlert
		} else {
			// Create new alert
			if err := s.store.CreateAlert(alert); err != nil {
				return fmt.Errorf("failed to create alert: %w", err)
			}
		}

		// Group alert into incident if it's firing
		if alert.Status == "firing" && alert.IncidentID == "" {
			if err := s.groupAlertIntoIncident(alert); err != nil {
				return fmt.Errorf("failed to group alert into incident: %w", err)
			}
		}
	}

	return nil
}

// findAlertByFingerprint finds an alert by its fingerprint
func (s *AlertService) findAlertByFingerprint(fingerprint string) (*models.Alert, error) {
	alerts, err := s.store.ListAlerts()
	if err != nil {
		return nil, err
	}

	for _, alert := range alerts {
		if alert.Fingerprint == fingerprint {
			return alert, nil
		}
	}

	return nil, storage.ErrNotFound
}

// groupAlertIntoIncident groups an alert into an appropriate incident
func (s *AlertService) groupAlertIntoIncident(alert *models.Alert) error {
	// Find existing incidents that this alert could be grouped into
	incidents, err := s.store.ListIncidents()
	if err != nil {
		return err
	}

	// Look for an open incident with similar labels
	for _, incident := range incidents {
		if incident.Status == models.IncidentStatusResolved {
			continue
		}

		if s.shouldGroupAlertWithIncident(alert, incident) {
			// Add alert to existing incident
			incident.AlertIDs = append(incident.AlertIDs, alert.ID)
			alert.IncidentID = incident.ID

			if err := s.store.UpdateIncident(incident); err != nil {
				return err
			}
			return s.store.UpdateAlert(alert)
		}
	}

	// Create new incident for this alert
	severity := s.determineSeverity(alert)
	title := s.generateIncidentTitle(alert)
	description := s.generateIncidentDescription(alert)

	incident, err := s.incidentService.CreateIncident(title, description, severity, []string{alert.ID})
	if err != nil {
		return err
	}

	alert.IncidentID = incident.ID
	return s.store.UpdateAlert(alert)
}

// shouldGroupAlertWithIncident determines if an alert should be grouped with an incident
func (s *AlertService) shouldGroupAlertWithIncident(alert *models.Alert, incident *models.Incident) bool {
	// Get first alert of the incident to compare
	if len(incident.AlertIDs) == 0 {
		return false
	}

	firstAlert, err := s.store.GetAlert(incident.AlertIDs[0])
	if err != nil {
		return false
	}

	// Group by service label
	if alert.Labels["service"] != "" && firstAlert.Labels["service"] != "" {
		return alert.Labels["service"] == firstAlert.Labels["service"]
	}

	// Group by instance label
	if alert.Labels["instance"] != "" && firstAlert.Labels["instance"] != "" {
		return alert.Labels["instance"] == firstAlert.Labels["instance"]
	}

	// Group by alertname
	if alert.Labels["alertname"] != "" && firstAlert.Labels["alertname"] != "" {
		return alert.Labels["alertname"] == firstAlert.Labels["alertname"]
	}

	return false
}

// determineSeverity determines the severity of an incident based on alert
func (s *AlertService) determineSeverity(alert *models.Alert) models.IncidentSeverity {
	severity, exists := alert.Labels["severity"]
	if !exists {
		severity = alert.Labels["priority"] // fallback
	}

	switch strings.ToLower(severity) {
	case "critical", "p0":
		return models.SeverityCritical
	case "high", "p1":
		return models.SeverityHigh
	case "medium", "p2":
		return models.SeverityMedium
	case "low", "p3":
		return models.SeverityLow
	default:
		return models.SeverityMedium
	}
}

// generateIncidentTitle generates a title for an incident from an alert
func (s *AlertService) generateIncidentTitle(alert *models.Alert) string {
	if summary := alert.Annotations["summary"]; summary != "" {
		return summary
	}

	if alertname := alert.Labels["alertname"]; alertname != "" {
		if instance := alert.Labels["instance"]; instance != "" {
			return fmt.Sprintf("%s on %s", alertname, instance)
		}
		return alertname
	}

	return "Alert Incident"
}

// generateIncidentDescription generates a description for an incident from an alert
func (s *AlertService) generateIncidentDescription(alert *models.Alert) string {
	description := ""

	if desc := alert.Annotations["description"]; desc != "" {
		description = desc
	} else if summary := alert.Annotations["summary"]; summary != "" {
		description = summary
	}

	// Add alert details
	if description != "" {
		description += "\n\n"
	}
	description += "Alert Details:\n"
	
	for key, value := range alert.Labels {
		description += fmt.Sprintf("- %s: %s\n", key, value)
	}

	return description
}

// GetAlert retrieves an alert by ID
func (s *AlertService) GetAlert(id string) (*models.Alert, error) {
	return s.store.GetAlert(id)
}

// ListAlerts retrieves all alerts
func (s *AlertService) ListAlerts() ([]*models.Alert, error) {
	return s.store.ListAlerts()
}