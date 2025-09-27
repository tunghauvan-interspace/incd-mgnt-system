# Repository Pattern Implementation

This document describes the repository pattern implementation for the incident management system.

## Overview

The repository pattern provides a clean abstraction layer over the database operations with support for:
- Context-aware operations
- Query filtering and pagination
- Transaction support
- Proper error handling
- Type-safe interfaces

## Quick Start

```go
package main

import (
    "context"
    "time"
    
    "github.com/tunghauvan-interspace/incd-mgnt-system/internal/config"
    "github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
    "github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

func main() {
    // Initialize PostgreSQL store
    cfg := &config.Config{
        DatabaseURL: "postgres://user:password@localhost/incident_db?sslmode=disable",
    }
    
    store, err := storage.NewPostgresStore(cfg)
    if err != nil {
        panic(err)
    }
    defer store.Close()
    
    // Create repository
    repo := storage.NewRepository(store)
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Use repository methods...
    incidents, err := repo.ListIncidents(ctx, storage.IncidentFilter{
        Limit: 10,
        OrderBy: "created_at",
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Found %d incidents\n", len(incidents))
}
```

## Interfaces

### IncidentRepository

```go
type IncidentRepository interface {
    CreateIncident(ctx context.Context, incident *models.Incident) error
    GetIncidentByID(ctx context.Context, id string) (*models.Incident, error)
    ListIncidents(ctx context.Context, filter IncidentFilter) ([]*models.Incident, error)
    UpdateIncident(ctx context.Context, incident *models.Incident) error
    DeleteIncident(ctx context.Context, id string) error
    CountIncidents(ctx context.Context, filter IncidentFilter) (int, error)
}
```

### AlertRepository

```go
type AlertRepository interface {
    CreateAlert(ctx context.Context, alert *models.Alert) error
    GetAlertByID(ctx context.Context, id string) (*models.Alert, error)
    ListAlerts(ctx context.Context, filter AlertFilter) ([]*models.Alert, error)
    UpdateAlert(ctx context.Context, alert *models.Alert) error
    DeleteAlert(ctx context.Context, id string) error
    CountAlerts(ctx context.Context, filter AlertFilter) (int, error)
}
```

## Filtering

### IncidentFilter

```go
type IncidentFilter struct {
    Status     *models.IncidentStatus   // Optional: filter by status
    Severity   *models.IncidentSeverity // Optional: filter by severity
    AssigneeID *string                  // Optional: filter by assignee
    Limit      int                      // Pagination: max results
    Offset     int                      // Pagination: skip results
    OrderBy    string                   // Sorting: "created_at", "updated_at", etc.
}
```

Example usage:

```go
// Get open high-severity incidents assigned to specific user
openStatus := models.IncidentStatusOpen
highSeverity := models.SeverityHigh
assigneeID := "user-123"

filter := storage.IncidentFilter{
    Status:     &openStatus,
    Severity:   &highSeverity,
    AssigneeID: &assigneeID,
    Limit:      25,
    OrderBy:    "created_at",
}

incidents, err := repo.ListIncidents(ctx, filter)
count, err := repo.CountIncidents(ctx, filter)
```

### AlertFilter

```go
type AlertFilter struct {
    Status      *string  // Optional: filter by status
    IncidentID  *string  // Optional: filter by incident ID
    Fingerprint *string  // Optional: filter by fingerprint
    Limit       int      // Pagination: max results
    Offset      int      // Pagination: skip results
    OrderBy     string   // Sorting: "created_at", "starts_at", etc.
}
```

## Transaction Support

Use transactions for operations that need to be atomic:

```go
err := repo.WithTransaction(ctx, func(txRepo storage.Repository) error {
    // Create incident
    incident := &models.Incident{
        ID:    "incident-123",
        Title: "Database Connection Lost",
        // ... other fields
    }
    
    if err := txRepo.CreateIncident(ctx, incident); err != nil {
        return err // This will rollback the transaction
    }
    
    // Create related alerts
    alert := &models.Alert{
        ID:         "alert-456", 
        IncidentID: incident.ID,
        // ... other fields
    }
    
    if err := txRepo.CreateAlert(ctx, alert); err != nil {
        return err // This will rollback the transaction
    }
    
    return nil // This will commit the transaction
})

if err != nil {
    // Handle transaction error
    log.Printf("Transaction failed: %v", err)
}
```

## Error Handling

The repository returns standard errors:

```go
incident, err := repo.GetIncidentByID(ctx, "some-id")
if err != nil {
    if errors.Is(err, storage.ErrNotFound) {
        // Handle not found case
        return nil, fmt.Errorf("incident not found: %s", "some-id")
    }
    // Handle other errors
    return nil, fmt.Errorf("failed to get incident: %w", err)
}
```

## Best Practices

1. **Always use context**: Pass context to all repository methods for proper cancellation and timeout handling.

2. **Use reasonable limits**: Always set `Limit` in filters to prevent accidentally loading large datasets.

3. **Handle ErrNotFound**: Check for `storage.ErrNotFound` when getting records by ID.

4. **Use transactions for related operations**: When creating/updating multiple related records, use `WithTransaction`.

5. **Use pointer fields in filters**: This allows distinguishing between "not filtered" (nil) and "filtered to specific value".

6. **Order results consistently**: Always specify `OrderBy` for predictable pagination.

## Migration from Old Interface

The repository pattern works alongside the existing `Store` interface:

```go
// Old way (still supported)
incidents, err := store.ListIncidents()

// New way with repository pattern
incidents, err := repo.ListIncidents(ctx, storage.IncidentFilter{
    Limit: 100,
    OrderBy: "created_at",
})
```

Both approaches work, but the repository pattern provides more flexibility and better practices.