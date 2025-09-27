# Low-Level Architecture

## Overview

This document provides a detailed technical architecture of the Incident Management System, including component interactions, data flows, API specifications, and implementation details.

## System Components

### 1. Application Layer

#### Main Entry Point (`cmd/server/main.go`)
- **Purpose**: Application bootstrap and dependency injection
- **Responsibilities**:
  - Load configuration from environment variables
  - Initialize storage layer (currently in-memory)
  - Create service instances with dependencies
  - Set up HTTP server and routing
  - Handle graceful shutdown
- **Key Dependencies**: config, handlers, services, storage packages

#### Configuration (`internal/config/config.go`)
- **Purpose**: Centralized configuration management
- **Structure**:
  ```go
  type Config struct {
      Port, LogLevel, DatabaseURL string
      SlackToken, SlackChannel string
      EmailSMTPHost, EmailSMTPPort, EmailUsername, EmailPassword string
      TelegramBotToken, TelegramChatID string
      AlertmanagerURL, AlertmanagerTimeout string
      MetricsEnabled bool
      MetricsPort string
  }
  ```
- **Loading Logic**: Environment variable parsing with defaults
- **Validation**: Basic type checking and required field validation

### 2. Service Layer

#### IncidentService (`internal/services/incident.go`)
- **Core Methods**:
  - `CreateIncident()`: Creates new incident with UUID generation
  - `GetIncident()`, `ListIncidents()`: CRUD operations
  - `AcknowledgeIncident()`, `ResolveIncident()`: Status transitions with timestamps
  - `CalculateMetrics()`: MTTA/MTTR calculation logic
- **Business Logic**:
  - Status validation (open → acknowledged → resolved)
  - Timestamp tracking (CreatedAt, AckedAt, ResolvedAt)
  - Assignee management
- **Dependencies**: Storage interface

#### AlertService (`internal/services/alert.go`)
- **Webhook Processing**:
  - `ProcessAlertmanagerWebhook()`: Main entry point for alert ingestion
  - Fingerprint-based deduplication
  - Alert grouping logic based on labels (service, instance, alertname)
- **Grouping Rules**:
  ```go
  // Priority: service > instance > alertname
  if alert.Labels["service"] == existingIncident.FirstAlert.Labels["service"] {
      // Group into incident
  }
  ```
- **Severity Mapping**:
  ```go
  switch strings.ToLower(severity) {
  case "critical", "p0": return SeverityCritical
  case "high", "p1": return SeverityHigh
  // ...
  }
  ```

#### NotificationService (`internal/services/notification.go`)
- **Channel Support**: Slack, Email (SMTP), Telegram
- **Message Templates**:
  - Incident creation: "🚨 New Incident Created"
  - Acknowledgment: "✅ Incident Acknowledged"
  - Resolution: "🎉 Incident Resolved"
- **Error Handling**: Continues with other channels if one fails
- **External API Calls**: HTTP requests with timeouts and error logging

### 3. Handler Layer (`internal/handlers/handlers.go`)

#### Route Structure
```
GET  /api/incidents           → handleListIncidents
GET  /api/incidents/{id}      → handleGetIncident
POST /api/incidents/{id}/acknowledge → handleAcknowledgeIncident
POST /api/incidents/{id}/resolve     → handleResolveIncident
GET  /api/alerts             → handleListAlerts
POST /api/webhooks/alertmanager     → handleAlertmanagerWebhook
GET  /api/metrics            → handleGetMetrics
GET  /health                 → handleHealth
GET  /                       → handleDashboard (serve template)
GET  /incidents              → handleIncidentsPage
GET  /alerts                 → handleAlertsPage
```

#### Request/Response Patterns
- **JSON API**: All `/api/*` endpoints return JSON
- **HTTP Status Codes**: 200 (success), 400 (bad request), 404 (not found), 500 (server error)
- **Error Handling**: Consistent error responses with descriptive messages

### 4. Storage Layer (`internal/storage/memory.go`)

#### Interface Definition
```go
type Store interface {
    // Incidents
    CreateIncident(*Incident) error
    GetIncident(id string) (*Incident, error)
    ListIncidents() ([]*Incident, error)
    UpdateIncident(*Incident) error
    DeleteIncident(id string) error

    // Alerts (similar pattern)
    // NotificationChannels, EscalationPolicies, OnCallSchedules
}
```

#### Current Implementation: MemoryStore
- **Concurrency**: `sync.RWMutex` for thread-safe operations
- **Data Structures**: `map[string]*Model` for each entity type
- **Limitations**: No persistence, data lost on restart
- **Performance**: O(1) for ID-based lookups, O(n) for list operations

### 5. Data Models (`internal/models/types.go`)

#### Core Entities

**Incident**:
```go
type Incident struct {
    ID, Title, Description string
    Status IncidentStatus  // open | acknowledged | resolved
    Severity IncidentSeverity // critical | high | medium | low
    CreatedAt, UpdatedAt time.Time
    AckedAt, ResolvedAt *time.Time
    AssigneeID string
    AlertIDs []string
    Labels map[string]string
}
```

**Alert**:
```go
type Alert struct {
    ID, Fingerprint string
    Status string  // firing | resolved
    StartsAt, EndsAt time.Time
    Labels, Annotations map[string]string
    IncidentID string
    CreatedAt time.Time
}
```

**Metrics**:
```go
type Metrics struct {
    TotalIncidents, OpenIncidents, ResolvedIncidents int
    MTTA, MTTR time.Duration
    IncidentsByStatus, IncidentsBySeverity map[string]int
}
```

## Data Flow Architecture

### Alert Ingestion Flow

```
Alertmanager Webhook
        ↓
handleAlertmanagerWebhook()
        ↓
AlertService.ProcessAlertmanagerWebhook()
        ↓
For each alert:
    - Check fingerprint deduplication
    - Create/Update alert in storage
    - If firing: groupAlertIntoIncident()
        ↓
        Find existing incident or create new
        ↓
        Update incident with alert ID
        ↓
        Send creation notification
```

### Incident Management Flow

```
User Action (acknowledge/resolve)
        ↓
HTTP Handler (handleAcknowledgeIncident)
        ↓
IncidentService.AcknowledgeIncident()
        ↓
Validate incident exists
        ↓
Update status, timestamps, assignee
        ↓
Persist to storage
        ↓
Send acknowledgment notification
```

### Dashboard Data Flow

```
Browser Request → HTTP Handler → Service Layer → Storage → JSON Response → Chart.js Rendering
```

## API Specifications

### REST API Endpoints

#### Incidents API
```
GET /api/incidents
- Query Params: status, severity, limit, offset
- Response: {"incidents": [Incident...], "total": number}

GET /api/incidents/{id}
- Response: Incident object

POST /api/incidents/{id}/acknowledge
- Body: {"assignee_id": "string"}
- Response: Updated Incident

POST /api/incidents/{id}/resolve
- Response: Updated Incident
```

#### Alerts API
```
GET /api/alerts
- Response: {"alerts": [Alert...]}

POST /api/webhooks/alertmanager
- Body: AlertmanagerWebhook payload
- Response: {"status": "ok"}
```

#### Metrics API
```
GET /api/metrics
- Response: Metrics object with MTTA/MTTR calculations
```

### Alertmanager Integration

#### Webhook Payload Structure
```json
{
  "version": "4",
  "groupKey": "...",
  "status": "firing",
  "receiver": "incident-management",
  "groupLabels": {...},
  "commonLabels": {...},
  "commonAnnotations": {...},
  "externalURL": "...",
  "alerts": [
    {
      "fingerprint": "...",
      "status": "firing",
      "startsAt": "2023-...",
      "endsAt": "0001-...",
      "labels": {"alertname": "...", "severity": "..."},
      "annotations": {"summary": "...", "description": "..."}
    }
  ]
}
```

#### Expected Response
- Status: 200 OK
- Body: `{"status": "ok"}`
- Processing: Asynchronous, no immediate response expected

## PostgreSQL Persistence Implementation Plan

### Migration Strategy

#### Phase 1: Database Setup and Connection
1. **Add PostgreSQL Dependencies**
   ```go
   // go.mod additions
   require (
       github.com/lib/pq v1.10.9
       github.com/golang-migrate/migrate/v4 v4.16.2
   )
   ```

2. **Database Connection Configuration**
   ```go
   type Config struct {
       // ... existing fields ...
       DatabaseURL string
       DBMaxOpenConns int    // default: 25
       DBMaxIdleConns int    // default: 5
       DBConnMaxLifetime time.Duration // default: 5m
   }
   ```

3. **Connection Pool Setup**
   ```go
   func NewPostgresStore(cfg *config.Config) (*PostgresStore, error) {
       db, err := sql.Open("postgres", cfg.DatabaseURL)
       if err != nil {
           return nil, err
       }

       db.SetMaxOpenConns(cfg.DBMaxOpenConns)
       db.SetMaxIdleConns(cfg.DBMaxIdleConns)
       db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

       if err := db.Ping(); err != nil {
           return nil, err
       }

       return &PostgresStore{db: db}, nil
   }
   ```

#### Phase 2: Schema Design and Migrations

##### Database Schema Design

**incidents table**
```sql
CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status incident_status NOT NULL DEFAULT 'open',
    severity incident_severity NOT NULL DEFAULT 'medium',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    acked_at TIMESTAMP WITH TIME ZONE,
    resolved_at TIMESTAMP WITH TIME ZONE,
    assignee_id VARCHAR(255),
    labels JSONB DEFAULT '{}',

    -- Constraints
    CONSTRAINT incidents_status_check CHECK (status IN ('open', 'acknowledged', 'resolved')),
    CONSTRAINT incidents_severity_check CHECK (severity IN ('critical', 'high', 'medium', 'low')),
    CONSTRAINT incidents_ack_sequence CHECK (
        (status = 'open' AND acked_at IS NULL AND resolved_at IS NULL) OR
        (status = 'acknowledged' AND acked_at IS NOT NULL AND resolved_at IS NULL) OR
        (status = 'resolved' AND acked_at IS NOT NULL AND resolved_at IS NOT NULL)
    ),
    CONSTRAINT incidents_timestamps_check CHECK (
        created_at <= updated_at AND
        (acked_at IS NULL OR acked_at >= created_at) AND
        (resolved_at IS NULL OR resolved_at >= acked_at)
    )
);

-- Custom enum types
CREATE TYPE incident_status AS ENUM ('open', 'acknowledged', 'resolved');
CREATE TYPE incident_severity AS ENUM ('critical', 'high', 'medium', 'low');
```

**alerts table**
```sql
CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fingerprint VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL,
    starts_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ends_at TIMESTAMP WITH TIME ZONE,
    labels JSONB DEFAULT '{}',
    annotations JSONB DEFAULT '{}',
    incident_id UUID REFERENCES incidents(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT alerts_fingerprint_unique UNIQUE (fingerprint),
    CONSTRAINT alerts_timestamps_check CHECK (
        starts_at <= ends_at OR ends_at IS NULL
    ),
    CONSTRAINT alerts_status_check CHECK (status IN ('firing', 'resolved'))
);
```

**notification_channels table (Future)**
```sql
CREATE TABLE notification_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- slack, email, telegram
    config JSONB NOT NULL DEFAULT '{}',
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT notification_channels_type_check CHECK (type IN ('slack', 'email', 'telegram'))
);
```

##### Indexes for Performance
```sql
-- Incidents indexes
CREATE INDEX idx_incidents_status_created ON incidents(status, created_at DESC);
CREATE INDEX idx_incidents_severity ON incidents(severity);
CREATE INDEX idx_incidents_assignee ON incidents(assignee_id) WHERE assignee_id IS NOT NULL;
CREATE INDEX idx_incidents_labels ON incidents USING GIN (labels);
CREATE INDEX idx_incidents_updated ON incidents(updated_at DESC);

-- Alerts indexes
CREATE INDEX idx_alerts_fingerprint ON alerts(fingerprint);
CREATE INDEX idx_alerts_incident_id ON alerts(incident_id);
CREATE INDEX idx_alerts_status_starts ON alerts(status, starts_at DESC);
CREATE INDEX idx_alerts_labels ON alerts USING GIN (labels);
CREATE INDEX idx_alerts_created ON alerts(created_at DESC);

-- Composite indexes for common queries
CREATE INDEX idx_incidents_status_severity ON incidents(status, severity);
CREATE INDEX idx_alerts_incident_status ON alerts(incident_id, status);
```

#### Phase 3: Repository Implementation

##### Repository Pattern Structure
```go
type PostgresStore struct {
    db *sql.DB
}

type IncidentRepository interface {
    Create(ctx context.Context, incident *models.Incident) error
    GetByID(ctx context.Context, id string) (*models.Incident, error)
    List(ctx context.Context, filter IncidentFilter) ([]*models.Incident, error)
    Update(ctx context.Context, incident *models.Incident) error
    Delete(ctx context.Context, id string) error
    Count(ctx context.Context, filter IncidentFilter) (int, error)
}

type IncidentFilter struct {
    Status    *models.IncidentStatus
    Severity  *models.IncidentSeverity
    AssigneeID *string
    Limit     int
    Offset    int
    OrderBy   string // "created_at", "updated_at", etc.
}
```

##### Sample Repository Methods
```go
func (r *PostgresStore) CreateIncident(ctx context.Context, incident *models.Incident) error {
    query := `
        INSERT INTO incidents (id, title, description, status, severity, assignee_id, labels)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING created_at, updated_at`

    labelsJSON, err := json.Marshal(incident.Labels)
    if err != nil {
        return fmt.Errorf("failed to marshal labels: %w", err)
    }

    err = r.db.QueryRowContext(ctx, query,
        incident.ID, incident.Title, incident.Description,
        incident.Status, incident.Severity, incident.AssigneeID, labelsJSON,
    ).Scan(&incident.CreatedAt, &incident.UpdatedAt)

    return err
}

func (r *PostgresStore) ListIncidents(ctx context.Context, filter IncidentFilter) ([]*models.Incident, error) {
    query := `
        SELECT id, title, description, status, severity, created_at, updated_at,
               acked_at, resolved_at, assignee_id, labels
        FROM incidents
        WHERE ($1::incident_status IS NULL OR status = $1)
          AND ($2::incident_severity IS NULL OR severity = $2)
          AND ($3::text IS NULL OR assignee_id = $3)
        ORDER BY ` + filter.OrderBy + ` DESC
        LIMIT $4 OFFSET $5`

    rows, err := r.db.QueryContext(ctx, query,
        filter.Status, filter.Severity, filter.AssigneeID,
        filter.Limit, filter.Offset)
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
            &incident.Status, &incident.Severity, &incident.CreatedAt,
            &incident.UpdatedAt, &incident.AckedAt, &incident.ResolvedAt,
            &incident.AssigneeID, &labelsJSON,
        )
        if err != nil {
            return nil, err
        }

        if err := json.Unmarshal(labelsJSON, &incident.Labels); err != nil {
            return nil, fmt.Errorf("failed to unmarshal labels: %w", err)
        }

        incidents = append(incidents, &incident)
    }

    return incidents, rows.Err()
}
```

#### Phase 4: Migration and Data Seeding

##### Database Migrations
```go
// migrations/001_initial_schema.up.sql
-- Create custom types
CREATE TYPE incident_status AS ENUM ('open', 'acknowledged', 'resolved');
CREATE TYPE incident_severity AS ENUM ('critical', 'high', 'medium', 'low');

-- Create incidents table
CREATE TABLE incidents (...);

-- Create alerts table
CREATE TABLE alerts (...);

-- Create indexes
CREATE INDEX ...;

// migrations/001_initial_schema.down.sql
DROP TABLE IF EXISTS alerts;
DROP TABLE IF EXISTS incidents;
DROP TYPE IF EXISTS incident_severity;
DROP TYPE IF EXISTS incident_status;
```

##### Migration Tool Integration
```go
import (
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(db *sql.DB) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return err
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "postgres", driver)
    if err != nil {
        return err
    }

    return m.Up()
}
```

#### Phase 5: Configuration and Environment Setup

##### Docker Compose Updates
```yaml
version: '3.8'
services:
  incident-management:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:password@postgres:5432/incidentdb?sslmode=disable
      - DB_MAX_OPEN_CONNS=25
      - DB_MAX_IDLE_CONNS=5
      - DB_CONN_MAX_LIFETIME=5m
    depends_on:
      - postgres
    # ... rest of config

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=incidentdb
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d incidentdb"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

##### Environment Variables
```bash
# Database Configuration
DATABASE_URL=postgres://user:password@localhost:5432/incidentdb?sslmode=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Migration Configuration
DB_MIGRATION_PATH=./migrations
```

#### Phase 6: Performance Optimization

##### Connection Pool Tuning
- **Max Open Connections**: 25 (based on expected concurrent users)
- **Max Idle Connections**: 5 (to maintain some warm connections)
- **Connection Max Lifetime**: 5 minutes (to prevent stale connections)

##### Query Optimization
```sql
-- Use EXPLAIN ANALYZE to identify slow queries
EXPLAIN ANALYZE SELECT * FROM incidents WHERE status = 'open' ORDER BY created_at DESC LIMIT 50;

-- Optimize with covering indexes
CREATE INDEX CONCURRENTLY idx_incidents_status_created_covering
ON incidents(status, created_at DESC) INCLUDE (id, title, severity);
```

##### Caching Strategy
```go
// Redis for metrics caching
type MetricsCache struct {
    redis *redis.Client
    ttl   time.Duration
}

func (c *MetricsCache) GetMetrics(ctx context.Context) (*models.Metrics, error) {
    key := "metrics:global"
    val, err := c.redis.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil // Cache miss
    }
    // Deserialize and return
}

func (c *MetricsCache) SetMetrics(ctx context.Context, metrics *models.Metrics) error {
    key := "metrics:global"
    data, _ := json.Marshal(metrics)
    return c.redis.Set(ctx, key, data, c.ttl).Err()
}
```

#### Phase 7: Backup and Recovery

##### Backup Strategy
```bash
# Daily backup script
#!/bin/bash
BACKUP_DIR="/backups"
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -h localhost -U user -d incidentdb > $BACKUP_DIR/backup_$DATE.sql

# Retention policy: keep last 30 days
find $BACKUP_DIR -name "backup_*.sql" -mtime +30 -delete
```

##### Point-in-Time Recovery
```sql
-- Create base backup
SELECT pg_start_backup('backup_label');

-- Copy data directory
-- Then stop backup
SELECT pg_stop_backup();
```

##### High Availability Setup (Future)
- PostgreSQL streaming replication
- Automatic failover with repmgr or Patroni
- Read replicas for scaling reads

#### Phase 8: Monitoring and Maintenance

##### Database Metrics to Monitor
```sql
-- Active connections
SELECT count(*) FROM pg_stat_activity WHERE state = 'active';

-- Table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename))
FROM pg_tables WHERE schemaname = 'public' ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Index usage
SELECT schemaname, tablename, indexname, idx_scan, idx_tup_read, idx_tup_fetch
FROM pg_stat_user_indexes ORDER BY idx_scan DESC;
```

##### Maintenance Tasks
```sql
-- Analyze tables for query planner
ANALYZE incidents;
ANALYZE alerts;

-- Vacuum for space reclamation
VACUUM (VERBOSE, ANALYZE) incidents;

-- Reindex if needed
REINDEX INDEX CONCURRENTLY idx_incidents_status_created;
```

#### Phase 9: Testing Strategy

##### Unit Tests for Repository
```go
func TestIncidentRepository_Create(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := NewPostgresStore(db)
    incident := &models.Incident{
        ID: "test-id",
        Title: "Test Incident",
        // ...
    }

    err := repo.CreateIncident(context.Background(), incident)
    assert.NoError(t, err)
    assert.NotZero(t, incident.CreatedAt)
}
```

##### Integration Tests
```go
func TestIncidentLifecycle(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    store := NewPostgresStore(db)
    service := services.NewIncidentService(store)

    // Create incident
    incident, err := service.CreateIncident("Test", "Description", models.SeverityHigh, nil)
    assert.NoError(t, err)

    // Acknowledge incident
    err = service.AcknowledgeIncident(incident.ID, "user123")
    assert.NoError(t, err)

    // Verify state
    updated, err := service.GetIncident(incident.ID)
    assert.Equal(t, models.IncidentStatusAcknowledged, updated.Status)
}
```

##### Database Migration Tests
```go
func TestMigrations(t *testing.T) {
    db := setupEmptyTestDB(t)
    defer db.Close()

    // Run migrations
    err := runMigrations(db)
    assert.NoError(t, err)

    // Verify schema
    var count int
    err = db.QueryRow("SELECT count(*) FROM information_schema.tables WHERE table_name = 'incidents'").Scan(&count)
    assert.NoError(t, err)
    assert.Equal(t, 1, count)
}
```

### Implementation Timeline

1. **Week 1**: Database setup, connection pooling, basic schema
2. **Week 2**: Repository implementation, migrations, basic CRUD
3. **Week 3**: Service layer updates, integration testing
4. **Week 4**: Performance optimization, monitoring, production deployment

### Success Criteria

- [ ] All existing functionality works with PostgreSQL
- [ ] Query performance < 100ms for common operations
- [ ] Data integrity maintained during migration
- [ ] Comprehensive test coverage for database operations
- [ ] Monitoring and alerting for database health
- [ ] Backup and recovery procedures documented and tested

## Security Architecture

### Authentication (Future)
- JWT-based authentication
- Role-based access control (RBAC)
- Session management with secure cookies

### Authorization
- API endpoint protection
- Incident ownership validation
- Admin-only operations

### Input Validation
- JSON schema validation for webhooks
- SQL injection prevention (prepared statements)
- XSS protection in templates

### Secrets Management
- Environment variables for sensitive data
- No hardcoded credentials
- Future: Integration with HashiCorp Vault

## Performance Considerations

### Current Bottlenecks
- In-memory storage: O(n) list operations
- No caching layer
- Synchronous notification sending
- Single-threaded processing

### Optimization Strategies
- Database indexing for common queries
- Redis caching for metrics
- Async notification processing
- Connection pooling
- Query optimization and pagination

### Scalability Patterns
- Horizontal scaling with load balancer
- Database read replicas
- Message queue for async processing
- CDN for static assets

## Deployment Architecture

### Docker Compose Setup
```yaml
services:
  incident-management:
    build: .
    ports: ["8080:8080"]
    environment: {...}
    depends_on: [prometheus, alertmanager]

  prometheus:
    image: prom/prometheus:latest
    volumes: ["./deployments/prometheus.yml:/etc/prometheus/prometheus.yml"]
    ports: ["9090:9090"]

  alertmanager:
    image: prom/alertmanager:latest
    volumes: ["./deployments/alertmanager.yml:/etc/alertmanager/alertmanager.yml"]
    ports: ["9093:9093"]
```

### Kubernetes Deployment (Future)
- Deployment manifests for each service
- ConfigMaps for configuration
- Secrets for sensitive data
- Ingress for external access
- PersistentVolumeClaims for database

## Monitoring & Observability

### Metrics to Collect
- HTTP request duration and status codes
- Webhook processing success/failure rates
- Notification delivery success rates
- Database query performance
- Memory and CPU usage

### Logging
- Structured logging with levels (debug, info, warn, error)
- Request ID tracing
- Error context and stack traces
- Audit logging for sensitive operations

### Health Checks
- `/health` endpoint for basic health
- `/ready` endpoint for readiness probes
- Database connectivity checks
- External service availability

## Error Handling Patterns

### Service Layer Errors
- Custom error types with context
- Error wrapping with `fmt.Errorf`
- Consistent error propagation

### HTTP Layer Errors
- HTTP status code mapping
- JSON error responses
- User-friendly error messages

### Recovery Strategies
- Graceful degradation (continue with partial failures)
- Retry logic with exponential backoff
- Circuit breaker pattern for external services

## Testing Architecture

### Unit Tests
- Service layer testing with mocks
- Handler testing with httptest
- Model validation testing

### Integration Tests
- API endpoint testing
- Database integration testing
- External service mocking

### End-to-End Tests
- Full webhook processing flow
- UI interaction testing
- Performance testing

## Future Architecture Extensions

### Microservices Decomposition
- Alert processing service
- Notification service
- Metrics aggregation service
- User management service

### Event-Driven Architecture
- Event sourcing for incident changes
- Message queues for async processing
- Event streaming for real-time updates

### Multi-Tenant Architecture
- Database schema per tenant
- Tenant context in requests
- Resource isolation and quotas

This low-level architecture document provides the technical foundation for implementing, maintaining, and extending the Incident Management System. All components are designed with modularity, testability, and scalability in mind.</content>
<filePath>c:\Users\tung4\incd-mgnt-system\Low-Level-Architecture.md