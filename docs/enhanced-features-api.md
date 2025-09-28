# Enhanced Incident Features API Documentation

This document demonstrates the enhanced incident management features implemented in Phase 1.

## API Endpoints

### 1. Comments and Timeline

#### Add a comment to an incident
```bash
POST /api/incidents/{incident_id}/comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "Investigating the root cause of the outage",
  "comment_type": "comment",
  "user_id": "user123"
}
```

#### Get all comments for an incident
```bash
GET /api/incidents/{incident_id}/comments
Authorization: Bearer <token>
```

#### Get incident timeline (comments + system events)
```bash
GET /api/incidents/{incident_id}/timeline
Authorization: Bearer <token>
```

### 2. Tags

#### Add tags to an incident
```bash
POST /api/incidents/{incident_id}/tags
Authorization: Bearer <token>
Content-Type: application/json

{
  "tags": [
    {
      "name": "environment",
      "value": "production",
      "color": "#ff0000"
    },
    {
      "name": "service",
      "value": "user-api",
      "color": "#00ff00"
    }
  ],
  "user_id": "user123"
}
```

#### Get tags for an incident
```bash
GET /api/incidents/{incident_id}/tags
Authorization: Bearer <token>
```

#### Remove tags from an incident
```bash
DELETE /api/incidents/{incident_id}/tags
Authorization: Bearer <token>
Content-Type: application/json

{
  "tag_names": ["environment"],
  "user_id": "user123"
}
```

### 3. Templates

#### Create an incident template
```bash
POST /api/templates
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Service Outage Template",
  "description": "Template for service outage incidents",
  "title_template": "Service Outage: {{service_name}}",
  "description_template": "Service {{service_name}} is experiencing a complete outage.\n\nImpact: {{impact}}\nAffected users: {{affected_users}}\nStarted at: {{start_time}}\n\n## Investigation Steps\n1. Check service status\n2. Review recent deployments\n3. Check infrastructure health\n4. Review error logs",
  "severity": "critical",
  "default_tags": [
    {
      "name": "outage",
      "value": "",
      "color": "#dc3545"
    }
  ]
}
```

#### List all templates
```bash
GET /api/templates
Authorization: Bearer <token>
```

#### Create incident from template
```bash
POST /api/incidents/from-template
Authorization: Bearer <token>
Content-Type: application/json

{
  "template_id": "template-123",
  "variables": {
    "service_name": "user-authentication-service",
    "impact": "Users cannot log in or access the platform",
    "affected_users": "~50,000 active users",
    "start_time": "2024-01-15 14:30 UTC"
  },
  "assignee_id": "oncall-engineer-123",
  "additional_tags": [
    {
      "name": "priority",
      "value": "p0",
      "color": "#ff0000"
    }
  ]
}
```

### 4. Advanced Search

#### Search incidents
```bash
POST /api/incidents/search
Authorization: Bearer <token>
Content-Type: application/json

{
  "query": "database connection timeout",
  "status": ["open", "acknowledged"],
  "severity": ["critical", "high"],
  "assignee_id": "user123",
  "tags": ["database", "timeout"],
  "created_after": "2024-01-01T00:00:00Z",
  "created_before": "2024-01-31T23:59:59Z",
  "page": 1,
  "limit": 20,
  "order_by": "created_at",
  "order_dir": "desc"
}
```

Response:
```json
{
  "incidents": [...],
  "total": 42,
  "page": 1,
  "limit": 20,
  "total_pages": 3
}
```

### 5. Bulk Operations

#### Bulk acknowledge incidents
```bash
POST /api/incidents/bulk
Authorization: Bearer <token>
Content-Type: application/json

{
  "incident_ids": ["incident-1", "incident-2", "incident-3"],
  "operation": "acknowledge",
  "parameters": {
    "assignee_id": "oncall-engineer"
  }
}
```

#### Bulk status update
```bash
POST /api/incidents/bulk
Authorization: Bearer <token>
Content-Type: application/json

{
  "incident_ids": ["incident-1", "incident-2"],
  "operation": "update_status",
  "parameters": {
    "status": "resolved"
  }
}
```

Response:
```json
{
  "processed_count": 2,
  "failed_count": 0,
  "failures": []
}
```

### 6. Assignment

#### Assign incident to user
```bash
POST /api/incidents/{incident_id}/assign
Authorization: Bearer <token>
Content-Type: application/json

{
  "assignee_id": "engineer-456",
  "user_id": "manager-123"
}
```

## Example Workflow

### 1. Create incident from template
```bash
# Use the "Critical Service Outage" template
curl -X POST https://api.example.com/api/incidents/from-template \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": "critical-outage-template",
    "variables": {
      "service_name": "payment-processor",
      "impact": "Payment processing is completely down",
      "affected_users": "All users attempting payments",
      "start_time": "2024-01-15 15:45 UTC"
    }
  }'
```

### 2. Add investigation comment
```bash
curl -X POST https://api.example.com/api/incidents/inc-123/comments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Initial investigation shows database connection pool exhaustion. Scaling up connection limits.",
    "comment_type": "comment"
  }'
```

### 3. Add tags for tracking
```bash
curl -X POST https://api.example.com/api/incidents/inc-123/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "tags": [
      {"name": "component", "value": "database", "color": "#007bff"},
      {"name": "root_cause", "value": "connection_pool", "color": "#ffc107"},
      {"name": "severity_reason", "value": "payment_down", "color": "#dc3545"}
    ]
  }'
```

### 4. Assign to engineer
```bash
curl -X POST https://api.example.com/api/incidents/inc-123/assign \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "assignee_id": "database-expert-456"
  }'
```

### 5. Search for similar incidents
```bash
curl -X POST https://api.example.com/api/incidents/search \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "database connection pool",
    "tags": ["database"],
    "severity": ["critical", "high"],
    "limit": 10
  }'
```

### 6. Bulk resolve related incidents
```bash
curl -X POST https://api.example.com/api/incidents/bulk \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "incident_ids": ["inc-123", "inc-124", "inc-125"],
    "operation": "update_status",
    "parameters": {
      "status": "resolved"
    }
  }'
```

## Features Demonstrated

1. **Template-based Incident Creation**: Rapidly create consistent incidents using predefined templates with variable substitution
2. **Rich Timeline**: Track all activities including comments, status changes, assignments, and tag modifications
3. **Flexible Tagging**: Organize incidents with colored, searchable tags that support key-value pairs
4. **Advanced Search**: Find incidents using full-text search combined with multiple filters
5. **Bulk Operations**: Efficiently manage multiple incidents simultaneously
6. **Assignment Workflow**: Clear assignment and reassignment with audit trail

These features transform the basic incident management system into a comprehensive, enterprise-ready incident response platform.