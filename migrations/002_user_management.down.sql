-- Drop foreign key constraint from incidents table if it exists
ALTER TABLE incidents DROP CONSTRAINT IF EXISTS fk_incidents_assignee_id;

-- Drop user management tables in reverse dependency order
DROP TABLE IF EXISTS user_activities;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;