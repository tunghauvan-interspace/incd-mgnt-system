-- Create notification_channels table for storing notification channel configurations
CREATE TABLE notification_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- slack, email, telegram
    config JSONB NOT NULL DEFAULT '{}',
    enabled BOOLEAN NOT NULL DEFAULT true,
    preferences JSONB, -- batching preferences, quiet hours, etc.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT notification_channels_type_check CHECK (type IN ('slack', 'email', 'telegram')),
    CONSTRAINT notification_channels_name_unique UNIQUE (name)
);

-- Create notification_templates table for customizable notification templates
CREATE TABLE notification_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- incident_created, incident_acked, incident_resolved, test
    channel VARCHAR(50) NOT NULL, -- slack, email, telegram
    subject VARCHAR(500), -- for email templates
    body TEXT NOT NULL, -- template content with variables
    variables JSONB DEFAULT '{}', -- available variables with descriptions
    is_default BOOLEAN NOT NULL DEFAULT false,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- user-specific template
    org_id VARCHAR(255), -- organization-specific template
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT notification_templates_type_check CHECK (type IN ('incident_created', 'incident_acked', 'incident_resolved', 'test')),
    CONSTRAINT notification_templates_channel_check CHECK (channel IN ('slack', 'email', 'telegram')),
    CONSTRAINT notification_templates_name_unique UNIQUE (name, type, channel)
);

-- Create notification_history table for tracking delivery attempts
CREATE TABLE notification_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID REFERENCES incidents(id) ON DELETE CASCADE,
    channel_id UUID REFERENCES notification_channels(id) ON DELETE SET NULL,
    template_id UUID REFERENCES notification_templates(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL, -- incident_created, incident_acked, incident_resolved, test
    channel VARCHAR(50) NOT NULL, -- slack, email, telegram
    recipient VARCHAR(500) NOT NULL, -- email address, slack user, telegram chat id
    subject VARCHAR(500), -- email subject
    content TEXT NOT NULL, -- final rendered content
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, sent, delivered, failed, retrying
    error_msg TEXT, -- error message if failed
    retry_count INTEGER NOT NULL DEFAULT 0,
    scheduled_at TIMESTAMP WITH TIME ZONE, -- for scheduled notifications
    sent_at TIMESTAMP WITH TIME ZONE, -- when notification was sent
    delivered_at TIMESTAMP WITH TIME ZONE, -- when delivery was confirmed (if supported)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT notification_history_type_check CHECK (type IN ('incident_created', 'incident_acked', 'incident_resolved', 'test')),
    CONSTRAINT notification_history_channel_check CHECK (channel IN ('slack', 'email', 'telegram')),
    CONSTRAINT notification_history_status_check CHECK (status IN ('pending', 'sent', 'delivered', 'failed', 'retrying')),
    CONSTRAINT notification_history_retry_count_check CHECK (retry_count >= 0)
);

-- Create notification_batches table for batch processing
CREATE TABLE notification_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id UUID REFERENCES notification_channels(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL, -- incident_created, incident_acked, incident_resolved
    count INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, sent, delivered, failed, retrying
    notifications JSONB NOT NULL DEFAULT '[]', -- array of notification history IDs
    scheduled_at TIMESTAMP WITH TIME ZONE, -- when to process the batch
    processed_at TIMESTAMP WITH TIME ZONE, -- when batch was processed
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT notification_batches_type_check CHECK (type IN ('incident_created', 'incident_acked', 'incident_resolved')),
    CONSTRAINT notification_batches_status_check CHECK (status IN ('pending', 'sent', 'delivered', 'failed', 'retrying')),
    CONSTRAINT notification_batches_count_check CHECK (count >= 0)
);

-- Create indexes for performance
CREATE INDEX idx_notification_channels_type ON notification_channels(type);
CREATE INDEX idx_notification_channels_enabled ON notification_channels(enabled) WHERE enabled = true;

CREATE INDEX idx_notification_templates_type_channel ON notification_templates(type, channel);
CREATE INDEX idx_notification_templates_is_default ON notification_templates(is_default) WHERE is_default = true;
CREATE INDEX idx_notification_templates_user_id ON notification_templates(user_id) WHERE user_id IS NOT NULL;

CREATE INDEX idx_notification_history_incident_id ON notification_history(incident_id) WHERE incident_id IS NOT NULL;
CREATE INDEX idx_notification_history_channel_id ON notification_history(channel_id) WHERE channel_id IS NOT NULL;
CREATE INDEX idx_notification_history_status ON notification_history(status);
CREATE INDEX idx_notification_history_created_at ON notification_history(created_at DESC);
CREATE INDEX idx_notification_history_type_channel ON notification_history(type, channel);
CREATE INDEX idx_notification_history_scheduled_at ON notification_history(scheduled_at) WHERE scheduled_at IS NOT NULL;

CREATE INDEX idx_notification_batches_channel_id ON notification_batches(channel_id) WHERE channel_id IS NOT NULL;
CREATE INDEX idx_notification_batches_status ON notification_batches(status);
CREATE INDEX idx_notification_batches_scheduled_at ON notification_batches(scheduled_at) WHERE scheduled_at IS NOT NULL;

-- Insert default notification templates
INSERT INTO notification_templates (name, type, channel, subject, body, variables, is_default) VALUES
-- Slack templates
('Default Incident Created (Slack)', 'incident_created', 'slack', NULL, 
 '🚨 *New Incident Created*\n\n*Title:* {{.Title}}\n*Severity:* {{.Severity}}\n*Status:* {{.Status}}\n*Created:* {{.CreatedAt}}\n\n{{if .Description}}*Description:*\n{{.Description}}\n\n{{end}}*Incident ID:* {{.ID}}',
 '{"Title": "Incident title", "Severity": "Incident severity", "Status": "Incident status", "CreatedAt": "Creation timestamp", "Description": "Incident description", "ID": "Incident unique ID"}', true),

('Default Incident Acknowledged (Slack)', 'incident_acked', 'slack', NULL,
 '✅ *Incident Acknowledged*\n\n*Title:* {{.Title}}\n*Severity:* {{.Severity}}\n*Status:* {{.Status}}\n*Acknowledged:* {{.AckedAt}}\n\n*Incident ID:* {{.ID}}',
 '{"Title": "Incident title", "Severity": "Incident severity", "Status": "Incident status", "AckedAt": "Acknowledgment timestamp", "ID": "Incident unique ID"}', true),

('Default Incident Resolved (Slack)', 'incident_resolved', 'slack', NULL,
 '🎉 *Incident Resolved*\n\n*Title:* {{.Title}}\n*Severity:* {{.Severity}}\n*Status:* {{.Status}}\n*Resolved:* {{.ResolvedAt}}\n\n*Incident ID:* {{.ID}}',
 '{"Title": "Incident title", "Severity": "Incident severity", "Status": "Incident status", "ResolvedAt": "Resolution timestamp", "ID": "Incident unique ID"}', true),

-- Email templates
('Default Incident Created (Email)', 'incident_created', 'email', 
 '🚨 New Incident: {{.Title}}',
 '<h2>🚨 New Incident Created</h2>\n<p><strong>Title:</strong> {{.Title}}</p>\n<p><strong>Severity:</strong> {{.Severity}}</p>\n<p><strong>Status:</strong> {{.Status}}</p>\n<p><strong>Created:</strong> {{.CreatedAt}}</p>\n{{if .Description}}<p><strong>Description:</strong><br>{{.Description}}</p>{{end}}\n<p><strong>Incident ID:</strong> {{.ID}}</p>',
 '{"Title": "Incident title", "Severity": "Incident severity", "Status": "Incident status", "CreatedAt": "Creation timestamp", "Description": "Incident description", "ID": "Incident unique ID"}', true),

('Default Incident Acknowledged (Email)', 'incident_acked', 'email',
 '✅ Incident Acknowledged: {{.Title}}',
 '<h2>✅ Incident Acknowledged</h2>\n<p><strong>Title:</strong> {{.Title}}</p>\n<p><strong>Severity:</strong> {{.Severity}}</p>\n<p><strong>Status:</strong> {{.Status}}</p>\n<p><strong>Acknowledged:</strong> {{.AckedAt}}</p>\n<p><strong>Incident ID:</strong> {{.ID}}</p>',
 '{"Title": "Incident title", "Severity": "Incident severity", "Status": "Incident status", "AckedAt": "Acknowledgment timestamp", "ID": "Incident unique ID"}', true),

('Default Incident Resolved (Email)', 'incident_resolved', 'email',
 '🎉 Incident Resolved: {{.Title}}',
 '<h2>🎉 Incident Resolved</h2>\n<p><strong>Title:</strong> {{.Title}}</p>\n<p><strong>Severity:</strong> {{.Severity}}</p>\n<p><strong>Status:</strong> {{.Status}}</p>\n<p><strong>Resolved:</strong> {{.ResolvedAt}}</p>\n<p><strong>Incident ID:</strong> {{.ID}}</p>',
 '{"Title": "Incident title", "Severity": "Incident severity", "Status": "Incident status", "ResolvedAt": "Resolution timestamp", "ID": "Incident unique ID"}', true),

-- Telegram templates
('Default Incident Created (Telegram)', 'incident_created', 'telegram', NULL,
 '🚨 <b>New Incident Created</b>\n\n<b>Title:</b> {{.Title}}\n<b>Severity:</b> {{.Severity}}\n<b>Status:</b> {{.Status}}\n<b>Created:</b> {{.CreatedAt}}\n\n{{if .Description}}<b>Description:</b>\n{{.Description}}\n\n{{end}}<b>Incident ID:</b> {{.ID}}',
 '{"Title": "Incident title", "Severity": "Incident severity", "Status": "Incident status", "CreatedAt": "Creation timestamp", "Description": "Incident description", "ID": "Incident unique ID"}', true),

('Default Incident Acknowledged (Telegram)', 'incident_acked', 'telegram', NULL,
 '✅ <b>Incident Acknowledged</b>\n\n<b>Title:</b> {{.Title}}\n<b>Severity:</b> {{.Severity}}\n<b>Status:</b> {{.Status}}\n<b>Acknowledged:</b> {{.AckedAt}}\n\n<b>Incident ID:</b> {{.ID}}',
 '{"Title": "Incident title", "Severity": "Incident severity", "Status": "Incident status", "AckedAt": "Acknowledgment timestamp", "ID": "Incident unique ID"}', true),

('Default Incident Resolved (Telegram)', 'incident_resolved', 'telegram', NULL,
 '🎉 <b>Incident Resolved</b>\n\n<b>Title:</b> {{.Title}}\n<b>Severity:</b> {{.Severity}}\n<b>Status:</b> {{.Status}}\n<b>Resolved:</b> {{.ResolvedAt}}\n\n<b>Incident ID:</b> {{.ID}}',
 '{"Title": "Incident title", "Severity": "Incident severity", "Status": "Incident status", "ResolvedAt": "Resolution timestamp", "ID": "Incident unique ID"}', true);