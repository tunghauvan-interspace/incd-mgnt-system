package storage

import (
	"errors"
	"sync"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
)

var (
	ErrNotFound = errors.New("not found")
)

// Store defines the storage interface
type Store interface {
	// Incidents
	GetIncident(id string) (*models.Incident, error)
	ListIncidents() ([]*models.Incident, error)
	CreateIncident(incident *models.Incident) error
	UpdateIncident(incident *models.Incident) error
	DeleteIncident(id string) error

	// Alerts
	GetAlert(id string) (*models.Alert, error)
	ListAlerts() ([]*models.Alert, error)
	CreateAlert(alert *models.Alert) error
	UpdateAlert(alert *models.Alert) error
	DeleteAlert(id string) error

	// Notification Channels
	GetNotificationChannel(id string) (*models.NotificationChannel, error)
	ListNotificationChannels() ([]*models.NotificationChannel, error)
	CreateNotificationChannel(channel *models.NotificationChannel) error
	UpdateNotificationChannel(channel *models.NotificationChannel) error
	DeleteNotificationChannel(id string) error

	// Escalation Policies
	GetEscalationPolicy(id string) (*models.EscalationPolicy, error)
	ListEscalationPolicies() ([]*models.EscalationPolicy, error)
	CreateEscalationPolicy(policy *models.EscalationPolicy) error
	UpdateEscalationPolicy(policy *models.EscalationPolicy) error
	DeleteEscalationPolicy(id string) error

	// On-Call Schedules
	GetOnCallSchedule(id string) (*models.OnCallSchedule, error)
	ListOnCallSchedules() ([]*models.OnCallSchedule, error)
	CreateOnCallSchedule(schedule *models.OnCallSchedule) error
	UpdateOnCallSchedule(schedule *models.OnCallSchedule) error
	DeleteOnCallSchedule(id string) error
}

// MemoryStore is an in-memory implementation of Store
type MemoryStore struct {
	incidents            map[string]*models.Incident
	alerts               map[string]*models.Alert
	notificationChannels map[string]*models.NotificationChannel
	escalationPolicies   map[string]*models.EscalationPolicy
	onCallSchedules      map[string]*models.OnCallSchedule
	mu                   sync.RWMutex
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() (*MemoryStore, error) {
	return &MemoryStore{
		incidents:            make(map[string]*models.Incident),
		alerts:               make(map[string]*models.Alert),
		notificationChannels: make(map[string]*models.NotificationChannel),
		escalationPolicies:   make(map[string]*models.EscalationPolicy),
		onCallSchedules:      make(map[string]*models.OnCallSchedule),
	}, nil
}

// Incident methods
func (s *MemoryStore) GetIncident(id string) (*models.Incident, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	incident, exists := s.incidents[id]
	if !exists {
		return nil, ErrNotFound
	}
	return incident, nil
}

func (s *MemoryStore) ListIncidents() ([]*models.Incident, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	incidents := make([]*models.Incident, 0, len(s.incidents))
	for _, incident := range s.incidents {
		incidents = append(incidents, incident)
	}
	return incidents, nil
}

func (s *MemoryStore) CreateIncident(incident *models.Incident) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.incidents[incident.ID] = incident
	return nil
}

func (s *MemoryStore) UpdateIncident(incident *models.Incident) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.incidents[incident.ID]; !exists {
		return ErrNotFound
	}
	s.incidents[incident.ID] = incident
	return nil
}

func (s *MemoryStore) DeleteIncident(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.incidents[id]; !exists {
		return ErrNotFound
	}
	delete(s.incidents, id)
	return nil
}

// Alert methods
func (s *MemoryStore) GetAlert(id string) (*models.Alert, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	alert, exists := s.alerts[id]
	if !exists {
		return nil, ErrNotFound
	}
	return alert, nil
}

func (s *MemoryStore) ListAlerts() ([]*models.Alert, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	alerts := make([]*models.Alert, 0, len(s.alerts))
	for _, alert := range s.alerts {
		alerts = append(alerts, alert)
	}
	return alerts, nil
}

func (s *MemoryStore) CreateAlert(alert *models.Alert) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.alerts[alert.ID] = alert
	return nil
}

func (s *MemoryStore) UpdateAlert(alert *models.Alert) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.alerts[alert.ID]; !exists {
		return ErrNotFound
	}
	s.alerts[alert.ID] = alert
	return nil
}

func (s *MemoryStore) DeleteAlert(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.alerts[id]; !exists {
		return ErrNotFound
	}
	delete(s.alerts, id)
	return nil
}

// NotificationChannel methods
func (s *MemoryStore) GetNotificationChannel(id string) (*models.NotificationChannel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	channel, exists := s.notificationChannels[id]
	if !exists {
		return nil, ErrNotFound
	}
	return channel, nil
}

func (s *MemoryStore) ListNotificationChannels() ([]*models.NotificationChannel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	channels := make([]*models.NotificationChannel, 0, len(s.notificationChannels))
	for _, channel := range s.notificationChannels {
		channels = append(channels, channel)
	}
	return channels, nil
}

func (s *MemoryStore) CreateNotificationChannel(channel *models.NotificationChannel) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notificationChannels[channel.ID] = channel
	return nil
}

func (s *MemoryStore) UpdateNotificationChannel(channel *models.NotificationChannel) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notificationChannels[channel.ID]; !exists {
		return ErrNotFound
	}
	s.notificationChannels[channel.ID] = channel
	return nil
}

func (s *MemoryStore) DeleteNotificationChannel(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notificationChannels[id]; !exists {
		return ErrNotFound
	}
	delete(s.notificationChannels, id)
	return nil
}

// EscalationPolicy methods
func (s *MemoryStore) GetEscalationPolicy(id string) (*models.EscalationPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	policy, exists := s.escalationPolicies[id]
	if !exists {
		return nil, ErrNotFound
	}
	return policy, nil
}

func (s *MemoryStore) ListEscalationPolicies() ([]*models.EscalationPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	policies := make([]*models.EscalationPolicy, 0, len(s.escalationPolicies))
	for _, policy := range s.escalationPolicies {
		policies = append(policies, policy)
	}
	return policies, nil
}

func (s *MemoryStore) CreateEscalationPolicy(policy *models.EscalationPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.escalationPolicies[policy.ID] = policy
	return nil
}

func (s *MemoryStore) UpdateEscalationPolicy(policy *models.EscalationPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.escalationPolicies[policy.ID]; !exists {
		return ErrNotFound
	}
	s.escalationPolicies[policy.ID] = policy
	return nil
}

func (s *MemoryStore) DeleteEscalationPolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.escalationPolicies[id]; !exists {
		return ErrNotFound
	}
	delete(s.escalationPolicies, id)
	return nil
}

// OnCallSchedule methods
func (s *MemoryStore) GetOnCallSchedule(id string) (*models.OnCallSchedule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	schedule, exists := s.onCallSchedules[id]
	if !exists {
		return nil, ErrNotFound
	}
	return schedule, nil
}

func (s *MemoryStore) ListOnCallSchedules() ([]*models.OnCallSchedule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	schedules := make([]*models.OnCallSchedule, 0, len(s.onCallSchedules))
	for _, schedule := range s.onCallSchedules {
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

func (s *MemoryStore) CreateOnCallSchedule(schedule *models.OnCallSchedule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.onCallSchedules[schedule.ID] = schedule
	return nil
}

func (s *MemoryStore) UpdateOnCallSchedule(schedule *models.OnCallSchedule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.onCallSchedules[schedule.ID]; !exists {
		return ErrNotFound
	}
	s.onCallSchedules[schedule.ID] = schedule
	return nil
}

func (s *MemoryStore) DeleteOnCallSchedule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.onCallSchedules[id]; !exists {
		return ErrNotFound
	}
	delete(s.onCallSchedules, id)
	return nil
}