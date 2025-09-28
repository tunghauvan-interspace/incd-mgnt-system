package storage

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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

	// User Management
	GetUser(id string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	ListUsers() ([]*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
	UpdateLastLogin(userID string, timestamp time.Time) error

	// Role Management
	GetRole(id string) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	ListRoles() ([]*models.Role, error)
	CreateRole(role *models.Role) error
	UpdateRole(role *models.Role) error
	DeleteRole(id string) error

	// User-Role Associations
	AssignRoleToUser(userID, roleID string) error
	RemoveRoleFromUser(userID, roleID string) error
	GetUserRoles(userID string) ([]*models.Role, error)
	
	// Permissions
	GetPermission(id string) (*models.Permission, error)
	ListPermissions() ([]*models.Permission, error)
	GetRolePermissions(roleID string) ([]*models.Permission, error)

	// User Activity Logging
	LogUserActivity(activity *models.UserActivity) error
	GetUserActivities(userID string, limit int) ([]*models.UserActivity, error)

	// Enhanced Incident Features - Comments
	CreateIncidentComment(comment *models.IncidentComment) error
	GetIncidentComments(incidentID string) ([]*models.IncidentComment, error)
	GetIncidentTimeline(incidentID string) ([]*models.IncidentComment, error)

	// Enhanced Incident Features - Tags  
	CreateIncidentTag(tag *models.IncidentTag) error
	GetIncidentTags(incidentID string) ([]*models.IncidentTag, error)
	DeleteIncidentTag(incidentID, tagName string) error

	// Enhanced Incident Features - Templates
	CreateIncidentTemplate(template *models.IncidentTemplate) error
	GetIncidentTemplate(id string) (*models.IncidentTemplate, error)
	ListIncidentTemplates(activeOnly bool) ([]*models.IncidentTemplate, error)
	UpdateIncidentTemplate(template *models.IncidentTemplate) error
	DeleteIncidentTemplate(id string) error

	// Enhanced Incident Features - Attachments
	CreateIncidentAttachment(attachment *models.IncidentAttachment) error
	GetIncidentAttachments(incidentID string) ([]*models.IncidentAttachment, error)
	DeleteIncidentAttachment(id string) error

	// Enhanced Incident Features - Search
	SearchIncidents(req *models.IncidentSearchRequest) ([]*models.Incident, int, error)

	// Close closes the store connection
	Close() error
}

// MemoryStore is an in-memory implementation of Store
type MemoryStore struct {
	incidents            map[string]*models.Incident
	alerts               map[string]*models.Alert
	notificationChannels map[string]*models.NotificationChannel
	escalationPolicies   map[string]*models.EscalationPolicy
	onCallSchedules      map[string]*models.OnCallSchedule
	users                map[string]*models.User
	usersByUsername      map[string]*models.User
	usersByEmail         map[string]*models.User
	roles                map[string]*models.Role
	rolesByName          map[string]*models.Role
	permissions          map[string]*models.Permission
	userRoles            map[string][]string // userID -> roleIDs
	rolePermissions      map[string][]string // roleID -> permissionIDs
	userActivities       map[string][]*models.UserActivity // userID -> activities
	// Enhanced incident features
	incidentComments     map[string][]*models.IncidentComment // incidentID -> comments
	incidentTags         map[string][]*models.IncidentTag     // incidentID -> tags
	incidentTemplates    map[string]*models.IncidentTemplate  // templateID -> template
	incidentAttachments  map[string][]*models.IncidentAttachment // incidentID -> attachments
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
		users:                make(map[string]*models.User),
		usersByUsername:      make(map[string]*models.User),
		usersByEmail:         make(map[string]*models.User),
		roles:                make(map[string]*models.Role),
		rolesByName:          make(map[string]*models.Role),
		permissions:          make(map[string]*models.Permission),
		userRoles:            make(map[string][]string),
		rolePermissions:      make(map[string][]string),
		userActivities:       make(map[string][]*models.UserActivity),
		// Enhanced incident features
		incidentComments:     make(map[string][]*models.IncidentComment),
		incidentTags:         make(map[string][]*models.IncidentTag),
		incidentTemplates:    make(map[string]*models.IncidentTemplate),
		incidentAttachments:  make(map[string][]*models.IncidentAttachment),
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

// Close closes the memory store (no-op for memory store)
func (s *MemoryStore) Close() error {
	return nil
}

// Enhanced Incident Features - Comments Implementation

func (s *MemoryStore) CreateIncidentComment(comment *models.IncidentComment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a copy to avoid external modifications
	commentCopy := *comment
	
	comments := s.incidentComments[comment.IncidentID]
	s.incidentComments[comment.IncidentID] = append(comments, &commentCopy)
	
	return nil
}

func (s *MemoryStore) GetIncidentComments(incidentID string) ([]*models.IncidentComment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	comments := s.incidentComments[incidentID]
	if len(comments) == 0 {
		return []*models.IncidentComment{}, nil
	}

	// Return copies to avoid external modifications
	result := make([]*models.IncidentComment, len(comments))
	for i, comment := range comments {
		commentCopy := *comment
		result[i] = &commentCopy
	}
	
	return result, nil
}

func (s *MemoryStore) GetIncidentTimeline(incidentID string) ([]*models.IncidentComment, error) {
	// For memory store, timeline is the same as comments (sorted by creation time)
	return s.GetIncidentComments(incidentID)
}

// Enhanced Incident Features - Tags Implementation

func (s *MemoryStore) CreateIncidentTag(tag *models.IncidentTag) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a copy to avoid external modifications
	tagCopy := *tag
	
	tags := s.incidentTags[tag.IncidentID]
	s.incidentTags[tag.IncidentID] = append(tags, &tagCopy)
	
	return nil
}

func (s *MemoryStore) GetIncidentTags(incidentID string) ([]*models.IncidentTag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tags := s.incidentTags[incidentID]
	if len(tags) == 0 {
		return []*models.IncidentTag{}, nil
	}

	// Return copies to avoid external modifications
	result := make([]*models.IncidentTag, len(tags))
	for i, tag := range tags {
		tagCopy := *tag
		result[i] = &tagCopy
	}
	
	return result, nil
}

func (s *MemoryStore) DeleteIncidentTag(incidentID, tagName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tags := s.incidentTags[incidentID]
	for i, tag := range tags {
		if tag.TagName == tagName {
			// Remove tag from slice
			s.incidentTags[incidentID] = append(tags[:i], tags[i+1:]...)
			return nil
		}
	}
	
	return ErrNotFound
}

// Enhanced Incident Features - Templates Implementation

func (s *MemoryStore) CreateIncidentTemplate(template *models.IncidentTemplate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a copy to avoid external modifications
	templateCopy := *template
	s.incidentTemplates[template.ID] = &templateCopy
	
	return nil
}

func (s *MemoryStore) GetIncidentTemplate(id string) (*models.IncidentTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	template, exists := s.incidentTemplates[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Return a copy to avoid external modifications
	templateCopy := *template
	return &templateCopy, nil
}

func (s *MemoryStore) ListIncidentTemplates(activeOnly bool) ([]*models.IncidentTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var templates []*models.IncidentTemplate
	for _, template := range s.incidentTemplates {
		if activeOnly && !template.IsActive {
			continue
		}
		
		// Create a copy to avoid external modifications
		templateCopy := *template
		templates = append(templates, &templateCopy)
	}
	
	return templates, nil
}

func (s *MemoryStore) UpdateIncidentTemplate(template *models.IncidentTemplate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.incidentTemplates[template.ID]; !exists {
		return ErrNotFound
	}

	// Create a copy to avoid external modifications
	templateCopy := *template
	s.incidentTemplates[template.ID] = &templateCopy
	
	return nil
}

func (s *MemoryStore) DeleteIncidentTemplate(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.incidentTemplates[id]; !exists {
		return ErrNotFound
	}

	delete(s.incidentTemplates, id)
	return nil
}

// Enhanced Incident Features - Attachments Implementation

func (s *MemoryStore) CreateIncidentAttachment(attachment *models.IncidentAttachment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a copy to avoid external modifications
	attachmentCopy := *attachment
	
	attachments := s.incidentAttachments[attachment.IncidentID]
	s.incidentAttachments[attachment.IncidentID] = append(attachments, &attachmentCopy)
	
	return nil
}

func (s *MemoryStore) GetIncidentAttachments(incidentID string) ([]*models.IncidentAttachment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	attachments := s.incidentAttachments[incidentID]
	if len(attachments) == 0 {
		return []*models.IncidentAttachment{}, nil
	}

	// Return copies to avoid external modifications
	result := make([]*models.IncidentAttachment, len(attachments))
	for i, attachment := range attachments {
		attachmentCopy := *attachment
		result[i] = &attachmentCopy
	}
	
	return result, nil
}

func (s *MemoryStore) DeleteIncidentAttachment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find and remove attachment by ID from all incidents
	for incidentID, attachments := range s.incidentAttachments {
		for i, attachment := range attachments {
			if attachment.ID == id {
				s.incidentAttachments[incidentID] = append(attachments[:i], attachments[i+1:]...)
				return nil
			}
		}
	}
	
	return ErrNotFound
}

// Enhanced Incident Features - Search Implementation

func (s *MemoryStore) SearchIncidents(req *models.IncidentSearchRequest) ([]*models.Incident, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matchingIncidents []*models.Incident

	// Filter incidents based on search criteria
	for _, incident := range s.incidents {
		if s.matchesSearchCriteria(incident, req) {
			// Create a copy to avoid external modifications
			incidentCopy := *incident
			matchingIncidents = append(matchingIncidents, &incidentCopy)
		}
	}

	// Sort incidents based on OrderBy and OrderDir
	s.sortIncidents(matchingIncidents, req.OrderBy, req.OrderDir)

	total := len(matchingIncidents)

	// Apply pagination
	start := (req.Page - 1) * req.Limit
	if start >= total {
		return []*models.Incident{}, total, nil
	}

	end := start + req.Limit
	if end > total {
		end = total
	}

	return matchingIncidents[start:end], total, nil
}

func (s *MemoryStore) matchesSearchCriteria(incident *models.Incident, req *models.IncidentSearchRequest) bool {
	// Text search in title and description
	if req.Query != "" {
		query := strings.ToLower(req.Query)
		if !strings.Contains(strings.ToLower(incident.Title), query) &&
		   !strings.Contains(strings.ToLower(incident.Description), query) {
			return false
		}
	}

	// Status filter
	if len(req.Status) > 0 {
		found := false
		for _, status := range req.Status {
			if incident.Status == status {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Severity filter
	if len(req.Severity) > 0 {
		found := false
		for _, severity := range req.Severity {
			if incident.Severity == severity {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Assignee filter
	if req.AssigneeID != nil && incident.AssigneeID != *req.AssigneeID {
		return false
	}

	// Date range filters
	if req.CreatedAfter != nil && incident.CreatedAt.Before(*req.CreatedAfter) {
		return false
	}
	if req.CreatedBefore != nil && incident.CreatedAt.After(*req.CreatedBefore) {
		return false
	}

	// Tag filter (simplified - would need to check incident tags in real implementation)
	if len(req.Tags) > 0 {
		incidentTags := s.incidentTags[incident.ID]
		for _, requiredTag := range req.Tags {
			found := false
			for _, tag := range incidentTags {
				if tag.TagName == requiredTag {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}

func (s *MemoryStore) sortIncidents(incidents []*models.Incident, orderBy, orderDir string) {
	// Simple sorting implementation - could be enhanced with more sophisticated sorting
	if orderBy == "" {
		orderBy = "created_at"
	}
	if orderDir == "" {
		orderDir = "desc"
	}

	// For simplicity, we'll implement basic sorting by created_at
	// A full implementation would handle all order fields
	if orderBy == "created_at" {
		if orderDir == "desc" {
			// Sort by creation time descending (newest first)
			for i := 0; i < len(incidents)-1; i++ {
				for j := i + 1; j < len(incidents); j++ {
					if incidents[i].CreatedAt.Before(incidents[j].CreatedAt) {
						incidents[i], incidents[j] = incidents[j], incidents[i]
					}
				}
			}
		} else {
			// Sort by creation time ascending (oldest first)  
			for i := 0; i < len(incidents)-1; i++ {
				for j := i + 1; j < len(incidents); j++ {
					if incidents[i].CreatedAt.After(incidents[j].CreatedAt) {
						incidents[i], incidents[j] = incidents[j], incidents[i]
					}
				}
			}
		}
	}
}

// User Management Methods

func (s *MemoryStore) GetUser(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, ErrNotFound
	}
	
	// Load user roles
	userCopy := *user
	roles, _ := s.getUserRolesInternal(id)
	userCopy.Roles = roles
	
	return &userCopy, nil
}

func (s *MemoryStore) GetUserByUsername(username string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.usersByUsername[username]
	if !exists {
		return nil, ErrNotFound
	}
	
	// Load user roles
	userCopy := *user
	roles, _ := s.getUserRolesInternal(user.ID)
	userCopy.Roles = roles
	
	return &userCopy, nil
}

func (s *MemoryStore) GetUserByEmail(email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.usersByEmail[email]
	if !exists {
		return nil, ErrNotFound
	}
	
	// Load user roles
	userCopy := *user
	roles, _ := s.getUserRolesInternal(user.ID)
	userCopy.Roles = roles
	
	return &userCopy, nil
}

func (s *MemoryStore) ListUsers() ([]*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		userCopy := *user
		roles, _ := s.getUserRolesInternal(user.ID)
		userCopy.Roles = roles
		users = append(users, &userCopy)
	}
	return users, nil
}

func (s *MemoryStore) CreateUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Check for existing username and email
	if _, exists := s.usersByUsername[user.Username]; exists {
		return errors.New("username already exists")
	}
	if _, exists := s.usersByEmail[user.Email]; exists {
		return errors.New("email already exists")
	}

	s.users[user.ID] = user
	s.usersByUsername[user.Username] = user
	s.usersByEmail[user.Email] = user
	return nil
}

func (s *MemoryStore) UpdateUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingUser, exists := s.users[user.ID]
	if !exists {
		return ErrNotFound
	}

	// Remove old username/email mappings
	delete(s.usersByUsername, existingUser.Username)
	delete(s.usersByEmail, existingUser.Email)

	// Check for conflicts with new username/email
	if user.Username != existingUser.Username {
		if _, exists := s.usersByUsername[user.Username]; exists {
			return errors.New("username already exists")
		}
	}
	if user.Email != existingUser.Email {
		if _, exists := s.usersByEmail[user.Email]; exists {
			return errors.New("email already exists")
		}
	}

	// Update mappings
	s.users[user.ID] = user
	s.usersByUsername[user.Username] = user
	s.usersByEmail[user.Email] = user
	return nil
}

func (s *MemoryStore) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return ErrNotFound
	}

	delete(s.users, id)
	delete(s.usersByUsername, user.Username)
	delete(s.usersByEmail, user.Email)
	delete(s.userRoles, id)
	delete(s.userActivities, id)
	return nil
}

func (s *MemoryStore) UpdateLastLogin(userID string, timestamp time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[userID]
	if !exists {
		return ErrNotFound
	}

	user.LastLogin = &timestamp
	return nil
}

// Role Management Methods

func (s *MemoryStore) GetRole(id string) (*models.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	role, exists := s.roles[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Load role permissions
	roleCopy := *role
	permissions, _ := s.getRolePermissionsInternal(id)
	roleCopy.Permissions = permissions

	return &roleCopy, nil
}

func (s *MemoryStore) GetRoleByName(name string) (*models.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	role, exists := s.rolesByName[name]
	if !exists {
		return nil, ErrNotFound
	}

	// Load role permissions
	roleCopy := *role
	permissions, _ := s.getRolePermissionsInternal(role.ID)
	roleCopy.Permissions = permissions

	return &roleCopy, nil
}

func (s *MemoryStore) ListRoles() ([]*models.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	roles := make([]*models.Role, 0, len(s.roles))
	for _, role := range s.roles {
		roleCopy := *role
		permissions, _ := s.getRolePermissionsInternal(role.ID)
		roleCopy.Permissions = permissions
		roles = append(roles, &roleCopy)
	}
	return roles, nil
}

func (s *MemoryStore) CreateRole(role *models.Role) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if role.ID == "" {
		role.ID = uuid.New().String()
	}

	if _, exists := s.rolesByName[role.Name]; exists {
		return errors.New("role name already exists")
	}

	s.roles[role.ID] = role
	s.rolesByName[role.Name] = role
	return nil
}

func (s *MemoryStore) UpdateRole(role *models.Role) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingRole, exists := s.roles[role.ID]
	if !exists {
		return ErrNotFound
	}

	// Remove old name mapping
	delete(s.rolesByName, existingRole.Name)

	// Check for name conflict
	if role.Name != existingRole.Name {
		if _, exists := s.rolesByName[role.Name]; exists {
			return errors.New("role name already exists")
		}
	}

	s.roles[role.ID] = role
	s.rolesByName[role.Name] = role
	return nil
}

func (s *MemoryStore) DeleteRole(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	role, exists := s.roles[id]
	if !exists {
		return ErrNotFound
	}

	delete(s.roles, id)
	delete(s.rolesByName, role.Name)
	delete(s.rolePermissions, id)

	// Remove role from all users
	for userID, roleIDs := range s.userRoles {
		newRoleIDs := make([]string, 0)
		for _, roleID := range roleIDs {
			if roleID != id {
				newRoleIDs = append(newRoleIDs, roleID)
			}
		}
		s.userRoles[userID] = newRoleIDs
	}

	return nil
}

// User-Role Association Methods

func (s *MemoryStore) AssignRoleToUser(userID, roleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify user and role exist
	if _, exists := s.users[userID]; !exists {
		return ErrNotFound
	}
	if _, exists := s.roles[roleID]; !exists {
		return ErrNotFound
	}

	// Check if already assigned
	roleIDs := s.userRoles[userID]
	for _, id := range roleIDs {
		if id == roleID {
			return nil // Already assigned
		}
	}

	s.userRoles[userID] = append(roleIDs, roleID)
	return nil
}

func (s *MemoryStore) RemoveRoleFromUser(userID, roleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	roleIDs := s.userRoles[userID]
	newRoleIDs := make([]string, 0)
	found := false

	for _, id := range roleIDs {
		if id != roleID {
			newRoleIDs = append(newRoleIDs, id)
		} else {
			found = true
		}
	}

	if !found {
		return ErrNotFound
	}

	s.userRoles[userID] = newRoleIDs
	return nil
}

func (s *MemoryStore) GetUserRoles(userID string) ([]*models.Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getUserRolesInternal(userID)
}

func (s *MemoryStore) getUserRolesInternal(userID string) ([]*models.Role, error) {
	roleIDs := s.userRoles[userID]
	roles := make([]*models.Role, 0, len(roleIDs))

	for _, roleID := range roleIDs {
		if role, exists := s.roles[roleID]; exists {
			roleCopy := *role
			permissions, _ := s.getRolePermissionsInternal(roleID)
			roleCopy.Permissions = permissions
			roles = append(roles, &roleCopy)
		}
	}

	return roles, nil
}

// Permission Methods

func (s *MemoryStore) GetPermission(id string) (*models.Permission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	permission, exists := s.permissions[id]
	if !exists {
		return nil, ErrNotFound
	}
	return permission, nil
}

func (s *MemoryStore) ListPermissions() ([]*models.Permission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	permissions := make([]*models.Permission, 0, len(s.permissions))
	for _, permission := range s.permissions {
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (s *MemoryStore) GetRolePermissions(roleID string) ([]*models.Permission, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getRolePermissionsInternal(roleID)
}

func (s *MemoryStore) getRolePermissionsInternal(roleID string) ([]*models.Permission, error) {
	permissionIDs := s.rolePermissions[roleID]
	permissions := make([]*models.Permission, 0, len(permissionIDs))

	for _, permID := range permissionIDs {
		if permission, exists := s.permissions[permID]; exists {
			permissions = append(permissions, permission)
		}
	}

	return permissions, nil
}

// User Activity Methods

func (s *MemoryStore) LogUserActivity(activity *models.UserActivity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if activity.ID == "" {
		activity.ID = uuid.New().String()
	}

	s.userActivities[activity.UserID] = append(s.userActivities[activity.UserID], activity)
	return nil
}

func (s *MemoryStore) GetUserActivities(userID string, limit int) ([]*models.UserActivity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	activities := s.userActivities[userID]
	if limit <= 0 || limit > len(activities) {
		limit = len(activities)
	}

	// Return most recent activities first
	result := make([]*models.UserActivity, 0, limit)
	for i := len(activities) - 1; i >= 0 && len(result) < limit; i-- {
		result = append(result, activities[i])
	}

	return result, nil
}
