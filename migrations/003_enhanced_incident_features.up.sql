-- Create incident_comments table for timeline tracking
CREATE TABLE incident_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    comment_type VARCHAR(50) NOT NULL DEFAULT 'comment', -- comment, status_change, assignment, etc.
    metadata JSONB DEFAULT '{}', -- additional context like old/new values for status changes
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT incident_comments_type_check CHECK (
        comment_type IN ('comment', 'status_change', 'assignment', 'severity_change', 'tag_added', 'tag_removed', 'attachment_added')
    )
);

-- Create incident_tags table for flexible tagging
CREATE TABLE incident_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    tag_name VARCHAR(100) NOT NULL,
    tag_value VARCHAR(500), -- optional value for key-value tags
    color VARCHAR(7) DEFAULT '#007bff', -- hex color code
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT incident_tags_name_check CHECK (tag_name ~ '^[a-zA-Z0-9_-]+$'),
    CONSTRAINT incident_tags_color_check CHECK (color ~ '^#[0-9A-Fa-f]{6}$'),
    UNIQUE(incident_id, tag_name, tag_value)
);

-- Create incident_templates table for reusable templates
CREATE TABLE incident_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    title_template VARCHAR(500) NOT NULL,
    description_template TEXT,
    severity incident_severity NOT NULL DEFAULT 'medium',
    default_tags JSONB DEFAULT '[]', -- array of {name, value, color} objects
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT incident_templates_name_unique UNIQUE (name)
);

-- Create incident_attachments table for runbook support
CREATE TABLE incident_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_path VARCHAR(500) NOT NULL, -- path to file in storage
    attachment_type VARCHAR(50) NOT NULL DEFAULT 'general', -- runbook, screenshot, log, general
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT incident_attachments_size_check CHECK (file_size > 0 AND file_size <= 52428800), -- 50MB max
    CONSTRAINT incident_attachments_type_check CHECK (
        attachment_type IN ('runbook', 'screenshot', 'log', 'document', 'general')
    )
);

-- Add indexes for performance
CREATE INDEX idx_incident_comments_incident_id ON incident_comments(incident_id);
CREATE INDEX idx_incident_comments_user_id ON incident_comments(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_incident_comments_created_at ON incident_comments(created_at DESC);
CREATE INDEX idx_incident_comments_type ON incident_comments(comment_type);

CREATE INDEX idx_incident_tags_incident_id ON incident_tags(incident_id);
CREATE INDEX idx_incident_tags_name ON incident_tags(tag_name);
CREATE INDEX idx_incident_tags_created_by ON incident_tags(created_by) WHERE created_by IS NOT NULL;

CREATE INDEX idx_incident_templates_name ON incident_templates(name);
CREATE INDEX idx_incident_templates_severity ON incident_templates(severity);
CREATE INDEX idx_incident_templates_is_active ON incident_templates(is_active) WHERE is_active = true;
CREATE INDEX idx_incident_templates_created_by ON incident_templates(created_by) WHERE created_by IS NOT NULL;

CREATE INDEX idx_incident_attachments_incident_id ON incident_attachments(incident_id);
CREATE INDEX idx_incident_attachments_type ON incident_attachments(attachment_type);
CREATE INDEX idx_incident_attachments_uploaded_by ON incident_attachments(uploaded_by) WHERE uploaded_by IS NOT NULL;

-- Add full-text search support for incidents
ALTER TABLE incidents ADD COLUMN search_vector tsvector;
CREATE INDEX idx_incidents_search_vector ON incidents USING gin(search_vector);

-- Create function to update search vector
CREATE OR REPLACE FUNCTION update_incident_search_vector() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := 
        setweight(to_tsvector('english', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(NEW.assignee_id::text, '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update search vector
CREATE TRIGGER update_incident_search_vector_trigger
    BEFORE INSERT OR UPDATE ON incidents
    FOR EACH ROW EXECUTE FUNCTION update_incident_search_vector();

-- Update existing incidents' search vectors
UPDATE incidents SET search_vector = 
    setweight(to_tsvector('english', COALESCE(title, '')), 'A') ||
    setweight(to_tsvector('english', COALESCE(description, '')), 'B') ||
    setweight(to_tsvector('english', COALESCE(assignee_id::text, '')), 'C');

-- Insert default incident templates
INSERT INTO incident_templates (name, description, title_template, description_template, severity, default_tags) VALUES
('Critical Service Outage', 'Template for critical service outage incidents', 
 'Service Outage: {{service_name}}', 
 'Service {{service_name}} is experiencing a complete outage.\n\nImpact: {{impact}}\nAffected users: {{affected_users}}\nStarted at: {{start_time}}\n\n## Investigation Steps\n1. Check service status\n2. Review recent deployments\n3. Check infrastructure health\n4. Review error logs\n\n## Communication\n- [ ] Notify status page\n- [ ] Update stakeholders\n- [ ] Post to incident channel',
 'critical',
 '[{"name": "outage", "value": "", "color": "#dc3545"}, {"name": "critical", "value": "", "color": "#dc3545"}]'),

('High CPU Alert', 'Template for high CPU utilization alerts',
 'High CPU Usage: {{server_name}}',
 'High CPU usage detected on {{server_name}}.\n\nCurrent CPU: {{cpu_percentage}}%\nThreshold: {{threshold}}%\nDuration: {{duration}}\n\n## Investigation Steps\n1. Check current processes\n2. Review resource utilization\n3. Check for runaway processes\n4. Scale resources if needed\n\n## Next Steps\n- [ ] Investigate root cause\n- [ ] Apply mitigation\n- [ ] Monitor metrics',
 'high',
 '[{"name": "performance", "value": "", "color": "#ffc107"}, {"name": "cpu", "value": "", "color": "#17a2b8"}]'),

('Security Incident', 'Template for security-related incidents',
 'Security Alert: {{incident_type}}',
 'Security incident detected: {{incident_type}}\n\nSeverity: {{severity}}\nAffected systems: {{systems}}\nDetected at: {{detection_time}}\n\n## Immediate Actions\n- [ ] Contain the incident\n- [ ] Preserve evidence\n- [ ] Notify security team\n- [ ] Document findings\n\n## Investigation\n- [ ] Analyze logs\n- [ ] Check for IOCs\n- [ ] Assess impact\n- [ ] Implement fixes',
 'high',
 '[{"name": "security", "value": "", "color": "#dc3545"}, {"name": "investigation", "value": "", "color": "#6f42c1"}]');