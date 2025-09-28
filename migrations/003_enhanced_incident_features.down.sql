-- Drop triggers and functions
DROP TRIGGER IF EXISTS update_incident_search_vector_trigger ON incidents;
DROP FUNCTION IF EXISTS update_incident_search_vector();

-- Drop indexes
DROP INDEX IF EXISTS idx_incidents_search_vector;
DROP INDEX IF EXISTS idx_incident_attachments_uploaded_by;
DROP INDEX IF EXISTS idx_incident_attachments_type;
DROP INDEX IF EXISTS idx_incident_attachments_incident_id;
DROP INDEX IF EXISTS idx_incident_templates_created_by;
DROP INDEX IF EXISTS idx_incident_templates_is_active;
DROP INDEX IF EXISTS idx_incident_templates_severity;
DROP INDEX IF EXISTS idx_incident_templates_name;
DROP INDEX IF EXISTS idx_incident_tags_created_by;
DROP INDEX IF EXISTS idx_incident_tags_name;
DROP INDEX IF EXISTS idx_incident_tags_incident_id;
DROP INDEX IF EXISTS idx_incident_comments_type;
DROP INDEX IF EXISTS idx_incident_comments_created_at;
DROP INDEX IF EXISTS idx_incident_comments_user_id;
DROP INDEX IF EXISTS idx_incident_comments_incident_id;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS incident_attachments;
DROP TABLE IF EXISTS incident_templates;
DROP TABLE IF EXISTS incident_tags;
DROP TABLE IF EXISTS incident_comments;

-- Remove search_vector column
ALTER TABLE incidents DROP COLUMN IF EXISTS search_vector;