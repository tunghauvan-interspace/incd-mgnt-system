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

// User represents a user in the system
type User struct {
	ID        string     `json:"id" db:"id"`
	Username  string     `json:"username" db:"username"`
	Email     string     `json:"email" db:"email"`
	FullName  string     `json:"full_name" db:"full_name"`
	Password  string     `json:"-" db:"password_hash"` // never expose password in JSON
	Roles     []*Role    `json:"roles,omitempty" db:"-"`  // populated by joins
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty" db:"last_login"`
}

// Role represents a role that can be assigned to users
type Role struct {
	ID          string        `json:"id" db:"id"`
	Name        string        `json:"name" db:"name"`
	DisplayName string        `json:"display_name" db:"display_name"`
	Description string        `json:"description" db:"description"`
	Permissions []*Permission `json:"permissions,omitempty" db:"-"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

// Permission represents a specific permission
type Permission struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Resource    string `json:"resource" db:"resource"`
	Action      string `json:"action" db:"action"`
	Description string `json:"description" db:"description"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID    string    `json:"user_id" db:"user_id"`
	RoleID    string    `json:"role_id" db:"role_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	RoleID       string `json:"role_id" db:"role_id"`
	PermissionID string `json:"permission_id" db:"permission_id"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,min=8"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"required,min=1,max=100"`
	Email    string `json:"email" validate:"required,email"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetConfirm represents a password reset confirmation
type PasswordResetConfirm struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// AuthResponse represents an authentication response with JWT token
type AuthResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	User         User      `json:"user"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// UserActivity represents user activity for audit trails
type UserActivity struct {
	ID          string                 `json:"id" db:"id"`
	UserID      string                 `json:"user_id" db:"user_id"`
	Action      string                 `json:"action" db:"action"`
	Resource    string                 `json:"resource" db:"resource"`
	ResourceID  string                 `json:"resource_id,omitempty" db:"resource_id"`
	IPAddress   string                 `json:"ip_address" db:"ip_address"`
	UserAgent   string                 `json:"user_agent" db:"user_agent"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
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

// IncidentComment represents a comment or timeline event on an incident
type IncidentComment struct {
	ID          string                 `json:"id" db:"id"`
	IncidentID  string                 `json:"incident_id" db:"incident_id"`
	UserID      *string                `json:"user_id,omitempty" db:"user_id"`
	User        *User                  `json:"user,omitempty" db:"-"` // populated by joins
	Content     string                 `json:"content" db:"content"`
	CommentType IncidentCommentType    `json:"comment_type" db:"comment_type"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}

// IncidentCommentType represents the type of timeline event
type IncidentCommentType string

const (
	CommentTypeComment         IncidentCommentType = "comment"
	CommentTypeStatusChange    IncidentCommentType = "status_change"
	CommentTypeAssignment      IncidentCommentType = "assignment"
	CommentTypeSeverityChange  IncidentCommentType = "severity_change"
	CommentTypeTagAdded        IncidentCommentType = "tag_added"
	CommentTypeTagRemoved      IncidentCommentType = "tag_removed"
	CommentTypeAttachmentAdded IncidentCommentType = "attachment_added"
)

// IncidentTag represents a tag applied to an incident
type IncidentTag struct {
	ID         string     `json:"id" db:"id"`
	IncidentID string     `json:"incident_id" db:"incident_id"`
	TagName    string     `json:"tag_name" db:"tag_name"`
	TagValue   *string    `json:"tag_value,omitempty" db:"tag_value"`
	Color      string     `json:"color" db:"color"`
	CreatedBy  *string    `json:"created_by,omitempty" db:"created_by"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	User       *User      `json:"created_by_user,omitempty" db:"-"` // populated by joins
}

// IncidentTemplate represents a template for creating incidents
type IncidentTemplate struct {
	ID                  string            `json:"id" db:"id"`
	Name                string            `json:"name" db:"name"`
	Description         string            `json:"description" db:"description"`
	TitleTemplate       string            `json:"title_template" db:"title_template"`
	DescriptionTemplate string            `json:"description_template" db:"description_template"`
	Severity            IncidentSeverity  `json:"severity" db:"severity"`
	DefaultTags         []TemplateTag     `json:"default_tags" db:"default_tags"`
	IsActive            bool              `json:"is_active" db:"is_active"`
	CreatedBy           *string           `json:"created_by,omitempty" db:"created_by"`
	CreatedAt           time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at" db:"updated_at"`
	User                *User             `json:"created_by_user,omitempty" db:"-"` // populated by joins
}

// TemplateTag represents a default tag in a template
type TemplateTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Color string `json:"color"`
}

// IncidentAttachment represents a file attached to an incident
type IncidentAttachment struct {
	ID             string              `json:"id" db:"id"`
	IncidentID     string              `json:"incident_id" db:"incident_id"`
	FileName       string              `json:"file_name" db:"file_name"`
	OriginalName   string              `json:"original_name" db:"original_name"`
	FileSize       int64               `json:"file_size" db:"file_size"`
	MimeType       string              `json:"mime_type" db:"mime_type"`
	FilePath       string              `json:"file_path" db:"file_path"`
	AttachmentType AttachmentType      `json:"attachment_type" db:"attachment_type"`
	UploadedBy     *string             `json:"uploaded_by,omitempty" db:"uploaded_by"`
	CreatedAt      time.Time           `json:"created_at" db:"created_at"`
	User           *User               `json:"uploaded_by_user,omitempty" db:"-"` // populated by joins
	DownloadURL    string              `json:"download_url,omitempty" db:"-"` // generated field
}

// AttachmentType represents the type of attachment
type AttachmentType string

const (
	AttachmentTypeRunbook   AttachmentType = "runbook"
	AttachmentTypeScreenshot AttachmentType = "screenshot"
	AttachmentTypeLog       AttachmentType = "log"
	AttachmentTypeDocument  AttachmentType = "document"
	AttachmentTypeGeneral   AttachmentType = "general"
)

// IncidentSearchRequest represents a search request for incidents
type IncidentSearchRequest struct {
	Query      string              `json:"query"`
	Status     []IncidentStatus    `json:"status"`
	Severity   []IncidentSeverity  `json:"severity"`
	AssigneeID *string             `json:"assignee_id"`
	Tags       []string            `json:"tags"`
	CreatedAfter  *time.Time       `json:"created_after"`
	CreatedBefore *time.Time       `json:"created_before"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	OrderBy    string              `json:"order_by"` // created_at, updated_at, severity
	OrderDir   string              `json:"order_dir"` // asc, desc
}

// IncidentSearchResponse represents a search response
type IncidentSearchResponse struct {
	Incidents    []*Incident `json:"incidents"`
	Total        int         `json:"total"`
	Page         int         `json:"page"`
	Limit        int         `json:"limit"`
	TotalPages   int         `json:"total_pages"`
}

// BulkOperationRequest represents a bulk operation request
type BulkOperationRequest struct {
	IncidentIDs []string           `json:"incident_ids"`
	Operation   BulkOperationType  `json:"operation"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// BulkOperationType represents the type of bulk operation
type BulkOperationType string

const (
	BulkOperationAcknowledge  BulkOperationType = "acknowledge"
	BulkOperationResolve      BulkOperationType = "resolve"
	BulkOperationUpdateStatus BulkOperationType = "update_status"
	BulkOperationAssign       BulkOperationType = "assign"
	BulkOperationAddTags      BulkOperationType = "add_tags"
	BulkOperationRemoveTags   BulkOperationType = "remove_tags"
)

// BulkOperationResponse represents the result of a bulk operation
type BulkOperationResponse struct {
	ProcessedCount int                      `json:"processed_count"`
	FailedCount    int                      `json:"failed_count"`
	Failures       []BulkOperationFailure   `json:"failures,omitempty"`
}

// BulkOperationFailure represents a failure in bulk operation
type BulkOperationFailure struct {
	IncidentID string `json:"incident_id"`
	Error      string `json:"error"`
}

// CreateIncidentFromTemplateRequest represents a request to create incident from template
type CreateIncidentFromTemplateRequest struct {
	TemplateID  string            `json:"template_id"`
	Variables   map[string]string `json:"variables"`
	AssigneeID  *string           `json:"assignee_id"`
	AdditionalTags []TemplateTag  `json:"additional_tags"`
}