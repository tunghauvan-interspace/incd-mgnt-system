-- Drop indexes first
DROP INDEX IF EXISTS idx_alerts_resolved;
DROP INDEX IF EXISTS idx_alerts_active;
DROP INDEX IF EXISTS idx_alerts_incident_starts;
DROP INDEX IF EXISTS idx_alerts_status_starts;
DROP INDEX IF EXISTS idx_alerts_annotations;
DROP INDEX IF EXISTS idx_alerts_labels;
DROP INDEX IF EXISTS idx_alerts_incident_id;
DROP INDEX IF EXISTS idx_alerts_created_at;
DROP INDEX IF EXISTS idx_alerts_starts_at;
DROP INDEX IF EXISTS idx_alerts_status;
DROP INDEX IF EXISTS idx_alerts_fingerprint;

DROP INDEX IF EXISTS idx_incidents_status_created_covering;
DROP INDEX IF EXISTS idx_incidents_labels;
DROP INDEX IF EXISTS idx_incidents_severity_created;
DROP INDEX IF EXISTS idx_incidents_status_created;
DROP INDEX IF EXISTS idx_incidents_assignee_id;
DROP INDEX IF EXISTS idx_incidents_updated_at;
DROP INDEX IF EXISTS idx_incidents_created_at;
DROP INDEX IF EXISTS idx_incidents_severity;
DROP INDEX IF EXISTS idx_incidents_status;

-- Drop tables in order due to foreign key constraints
DROP TABLE IF EXISTS alerts;
DROP TABLE IF EXISTS incidents;

-- Drop custom enum types
DROP TYPE IF EXISTS incident_severity;
DROP TYPE IF EXISTS incident_status;