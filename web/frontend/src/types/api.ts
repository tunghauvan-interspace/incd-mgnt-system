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