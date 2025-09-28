export interface Incident {
  id: string
  title: string
  description: string
  severity: string
  status: string
  assignee_id?: string
  created_at: string
  updated_at: string
  acknowledged_at?: string
  resolved_at?: string
  labels?: Record<string, string>
}

export interface Alert {
  id: string
  alert_name: string
  generator_url: string
  status: string
  starts_at: string
  ends_at?: string
  labels?: Record<string, string>
  annotations?: Record<string, string>
  incident_id?: string
  created_at: string
  updated_at: string
}

export interface Metrics {
  total_incidents: number
  open_incidents: number
  mtta: number // Mean Time To Acknowledge in nanoseconds
  mttr: number // Mean Time To Resolve in nanoseconds
  incidents_by_status: Record<string, number>
  incidents_by_severity: Record<string, number>
}

export interface AcknowledgeIncidentRequest {
  assignee_id: string
}

export interface LoginRequest {
  username: string
  email?: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  full_name: string
  password: string
}

export interface AuthResponse {
  token: string
  refresh_token: string
  user: User
  expires_at: string
}

export interface User {
  id: string
  username: string
  email: string
  full_name: string
  roles: Role[]
  is_active: boolean
  created_at: string
  updated_at: string
  last_login?: string
}

export interface Role {
  id: string
  name: string
  display_name: string
  description?: string
  permissions: Permission[]
  created_at: string
  updated_at: string
}

export interface Permission {
  id: string
  name: string
  resource: string
  action: string
  description?: string
}
