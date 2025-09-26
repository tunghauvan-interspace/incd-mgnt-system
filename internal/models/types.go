package models

import (
	"time"
)

// IncidentStatus represents the status of an incident
type IncidentStatus string

const (
	IncidentStatusOpen         IncidentStatus = "open"
	IncidentStatusAcknowledged IncidentStatus = "acknowledged"
	IncidentStatusResolved     IncidentStatus = "resolved"
)

// IncidentSeverity represents the severity level of an incident
type IncidentSeverity string

const (
	SeverityCritical IncidentSeverity = "critical"
	SeverityHigh     IncidentSeverity = "high"
	SeverityMedium   IncidentSeverity = "medium"
	SeverityLow      IncidentSeverity = "low"
)

// Incident represents an incident in the system
type Incident struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      IncidentStatus    `json:"status"`
	Severity    IncidentSeverity  `json:"severity"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	AckedAt     *time.Time        `json:"acked_at,omitempty"`
	ResolvedAt  *time.Time        `json:"resolved_at,omitempty"`
	AssigneeID  string            `json:"assignee_id,omitempty"`
	AlertIDs    []string          `json:"alert_ids"`
	Labels      map[string]string `json:"labels"`
}

// Alert represents an alert from Prometheus/Alertmanager
type Alert struct {
	ID          string            `json:"id"`
	Fingerprint string            `json:"fingerprint"`
	Status      string            `json:"status"`
	StartsAt    time.Time         `json:"starts_at"`
	EndsAt      time.Time         `json:"ends_at"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	IncidentID  string            `json:"incident_id,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

// NotificationChannel represents a notification destination
type NotificationChannel struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"` // slack, email, telegram
	Config   map[string]string `json:"config"`
	Enabled  bool              `json:"enabled"`
	Name     string            `json:"name"`
}

// EscalationPolicy defines how incidents should be escalated
type EscalationPolicy struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Rules []EscalationRule `json:"rules"`
}

// EscalationRule defines a single escalation step
type EscalationRule struct {
	DelayMinutes int      `json:"delay_minutes"`
	Targets      []string `json:"targets"` // user IDs or notification channel IDs
}

// OnCallSchedule represents an on-call schedule
type OnCallSchedule struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Timezone string        `json:"timezone"`
	Layers   []ScheduleLayer `json:"layers"`
}

// ScheduleLayer represents a layer in the on-call schedule
type ScheduleLayer struct {
	Name         string        `json:"name"`
	Users        []string      `json:"users"` // user IDs
	Rotation     RotationType  `json:"rotation"`
	Start        time.Time     `json:"start"`
	Restrictions []Restriction `json:"restrictions"`
}

// RotationType defines how the rotation works
type RotationType struct {
	Type     string `json:"type"` // daily, weekly, monthly
	Length   int    `json:"length"`
	Handoff  string `json:"handoff"` // time of handoff
}

// Restriction defines time restrictions for on-call
type Restriction struct {
	Type       string `json:"type"` // daily, weekly
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	StartDay   int    `json:"start_day"`   // for weekly restrictions
	EndDay     int    `json:"end_day"`     // for weekly restrictions
}

// Metrics represents incident metrics for dashboard
type Metrics struct {
	TotalIncidents     int           `json:"total_incidents"`
	OpenIncidents      int           `json:"open_incidents"`
	ResolvedIncidents  int           `json:"resolved_incidents"`
	MTTA               time.Duration `json:"mtta"` // Mean Time To Acknowledge
	MTTR               time.Duration `json:"mttr"` // Mean Time To Resolve
	IncidentsByStatus  map[string]int `json:"incidents_by_status"`
	IncidentsBySeverity map[string]int `json:"incidents_by_severity"`
}