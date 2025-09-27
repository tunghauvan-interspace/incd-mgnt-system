-- Create custom enum types for incident status and severity
CREATE TYPE incident_status AS ENUM ('open', 'acknowledged', 'resolved');
CREATE TYPE incident_severity AS ENUM ('critical', 'high', 'medium', 'low');

-- Create incidents table
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

-- Create alerts table
CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fingerprint VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL,
    starts_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ends_at TIMESTAMP WITH TIME ZONE,
    labels JSONB DEFAULT '{}',
    annotations JSONB DEFAULT '{}',
    incident_id UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    -- Foreign key constraint
    CONSTRAINT fk_alerts_incident_id FOREIGN KEY (incident_id) REFERENCES incidents(id) ON DELETE SET NULL,
    
    -- Constraints
    CONSTRAINT alerts_timestamps_check CHECK (
        starts_at <= COALESCE(ends_at, starts_at) AND
        created_at >= starts_at
    )
);

-- Create comprehensive indexes for performance

-- Incidents indexes
CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_incidents_severity ON incidents(severity);
CREATE INDEX idx_incidents_created_at ON incidents(created_at DESC);
CREATE INDEX idx_incidents_updated_at ON incidents(updated_at DESC);
CREATE INDEX idx_incidents_assignee_id ON incidents(assignee_id) WHERE assignee_id IS NOT NULL;
CREATE INDEX idx_incidents_status_created ON incidents(status, created_at DESC);
CREATE INDEX idx_incidents_severity_created ON incidents(severity, created_at DESC);
CREATE INDEX idx_incidents_labels ON incidents USING gin(labels);

-- Covering index for common queries
CREATE INDEX idx_incidents_status_created_covering 
ON incidents(status, created_at DESC) 
INCLUDE (id, title, severity, assignee_id);

-- Alerts indexes
CREATE INDEX idx_alerts_fingerprint ON alerts(fingerprint);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_starts_at ON alerts(starts_at DESC);
CREATE INDEX idx_alerts_created_at ON alerts(created_at DESC);
CREATE INDEX idx_alerts_incident_id ON alerts(incident_id) WHERE incident_id IS NOT NULL;
CREATE INDEX idx_alerts_labels ON alerts USING gin(labels);
CREATE INDEX idx_alerts_annotations ON alerts USING gin(annotations);

-- Composite indexes for common query patterns
CREATE INDEX idx_alerts_status_starts ON alerts(status, starts_at DESC);
CREATE INDEX idx_alerts_incident_starts ON alerts(incident_id, starts_at DESC) WHERE incident_id IS NOT NULL;

-- Partial indexes for active alerts (performance optimization)
CREATE INDEX idx_alerts_active ON alerts(starts_at DESC) WHERE ends_at IS NULL;
CREATE INDEX idx_alerts_resolved ON alerts(ends_at DESC) WHERE ends_at IS NOT NULL;