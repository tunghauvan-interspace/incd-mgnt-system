-- Create permissions table
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT permissions_name_unique UNIQUE (name),
    CONSTRAINT permissions_resource_action_unique UNIQUE (resource, action)
);

-- Create roles table
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT roles_name_check CHECK (name ~ '^[a-z_]+$')
);

-- Create role_permissions junction table
CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    PRIMARY KEY (role_id, permission_id)
);

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_login TIMESTAMP WITH TIME ZONE,

    CONSTRAINT users_username_check CHECK (username ~ '^[a-zA-Z0-9_-]+$'),
    CONSTRAINT users_email_check CHECK (email ~ '^[^@]+@[^@]+\.[^@]+$'),
    CONSTRAINT users_password_check CHECK (length(password_hash) >= 60) -- bcrypt hash length
);

-- Create user_roles junction table
CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, role_id)
);

-- Create user_activities table for audit trails
CREATE TABLE user_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    resource_id VARCHAR(255),
    ip_address INET,
    user_agent TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active) WHERE is_active = true;
CREATE INDEX idx_users_created_at ON users(created_at DESC);
CREATE INDEX idx_users_last_login ON users(last_login DESC) WHERE last_login IS NOT NULL;

CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_created_at ON roles(created_at DESC);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

CREATE INDEX idx_user_activities_user_id ON user_activities(user_id);
CREATE INDEX idx_user_activities_action ON user_activities(action);
CREATE INDEX idx_user_activities_resource ON user_activities(resource);
CREATE INDEX idx_user_activities_created_at ON user_activities(created_at DESC);
CREATE INDEX idx_user_activities_resource_id ON user_activities(resource_id) WHERE resource_id IS NOT NULL;

-- Insert default permissions
INSERT INTO permissions (name, resource, action, description) VALUES
-- Incident permissions
('incidents.read', 'incidents', 'read', 'View incidents'),
('incidents.create', 'incidents', 'create', 'Create new incidents'),
('incidents.update', 'incidents', 'update', 'Update incident details'),
('incidents.delete', 'incidents', 'delete', 'Delete incidents'),
('incidents.acknowledge', 'incidents', 'acknowledge', 'Acknowledge incidents'),
('incidents.resolve', 'incidents', 'resolve', 'Resolve incidents'),
('incidents.assign', 'incidents', 'assign', 'Assign incidents to users'),

-- Alert permissions
('alerts.read', 'alerts', 'read', 'View alerts'),
('alerts.update', 'alerts', 'update', 'Update alert details'),
('alerts.delete', 'alerts', 'delete', 'Delete alerts'),

-- User management permissions
('users.read', 'users', 'read', 'View users'),
('users.create', 'users', 'create', 'Create new users'),
('users.update', 'users', 'update', 'Update user details'),
('users.delete', 'users', 'delete', 'Delete users'),
('users.manage_roles', 'users', 'manage_roles', 'Assign roles to users'),

-- Role management permissions
('roles.read', 'roles', 'read', 'View roles'),
('roles.create', 'roles', 'create', 'Create new roles'),
('roles.update', 'roles', 'update', 'Update role details'),
('roles.delete', 'roles', 'delete', 'Delete roles'),
('roles.manage_permissions', 'roles', 'manage_permissions', 'Assign permissions to roles'),

-- System permissions
('metrics.read', 'metrics', 'read', 'View system metrics'),
('system.health', 'system', 'health', 'View system health'),
('audit.read', 'audit', 'read', 'View audit logs');

-- Insert default roles
INSERT INTO roles (name, display_name, description) VALUES
('admin', 'Administrator', 'Full system access with all permissions'),
('responder', 'Incident Responder', 'Can manage incidents and alerts'),
('viewer', 'Viewer', 'Read-only access to incidents and alerts');

-- Get role IDs for permission assignments
DO $$
DECLARE
    admin_role_id UUID;
    responder_role_id UUID;
    viewer_role_id UUID;
    perm_record RECORD;
BEGIN
    -- Get role IDs
    SELECT id INTO admin_role_id FROM roles WHERE name = 'admin';
    SELECT id INTO responder_role_id FROM roles WHERE name = 'responder';
    SELECT id INTO viewer_role_id FROM roles WHERE name = 'viewer';

    -- Admin gets all permissions
    FOR perm_record IN SELECT id FROM permissions LOOP
        INSERT INTO role_permissions (role_id, permission_id) 
        VALUES (admin_role_id, perm_record.id);
    END LOOP;

    -- Responder permissions
    INSERT INTO role_permissions (role_id, permission_id)
    SELECT responder_role_id, id FROM permissions WHERE name IN (
        'incidents.read', 'incidents.create', 'incidents.update', 'incidents.acknowledge', 
        'incidents.resolve', 'incidents.assign', 'alerts.read', 'alerts.update', 
        'metrics.read', 'system.health'
    );

    -- Viewer permissions
    INSERT INTO role_permissions (role_id, permission_id)
    SELECT viewer_role_id, id FROM permissions WHERE name IN (
        'incidents.read', 'alerts.read', 'metrics.read', 'system.health'
    );
END $$;

-- Add assignee_id foreign key constraint to incidents table (if not exists)
-- This will reference the users table for incident assignments
DO $$
BEGIN
    -- Check if the constraint doesn't already exist
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_incidents_assignee_id' 
        AND table_name = 'incidents'
    ) THEN
        -- Add the foreign key constraint
        ALTER TABLE incidents 
        ADD CONSTRAINT fk_incidents_assignee_id 
        FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL;
        
        -- Also update the assignee_id column type to UUID if it's not already
        ALTER TABLE incidents ALTER COLUMN assignee_id TYPE UUID USING assignee_id::UUID;
    END IF;
EXCEPTION
    WHEN others THEN
        -- If there's an error (e.g., data incompatibility), just continue
        RAISE NOTICE 'Could not add foreign key constraint to assignee_id: %', SQLERRM;
END $$;