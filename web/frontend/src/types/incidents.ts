import type { User } from './auth'

export interface Incident {
  id: string
  title: string
  description: string
  status: 'open' | 'acknowledged' | 'investigating' | 'resolved' | 'closed'
  severity: 'critical' | 'high' | 'medium' | 'low'
  assignee?: User
  reporter: User
  createdAt: string
  updatedAt: string
  resolvedAt?: string
  labels: string[]
  tags: { [key: string]: string }
  alertCount: number
  timelineEvents: IncidentTimelineEvent[]
}

export interface IncidentTimelineEvent {
  id: string
  type: 'created' | 'status_changed' | 'assigned' | 'comment' | 'resolved'
  description: string
  user: User
  timestamp: string
  metadata?: { [key: string]: any }
}

export interface IncidentFilter {
  status?: string
  severity?: string
  assignee?: string
  search?: string
  dateFrom?: string
  dateTo?: string
}

export interface IncidentStats {
  total: number
  open: number
  resolved: number
  critical: number
  averageResolutionTime: number
}

export interface CreateIncidentData {
  title: string
  description: string
  severity: Incident['severity']
  assigneeId?: string
  labels?: string[]
  tags?: { [key: string]: string }
}

export interface UpdateIncidentData {
  title?: string
  description?: string
  status?: Incident['status']
  severity?: Incident['severity']
  assigneeId?: string
  labels?: string[]
  tags?: { [key: string]: string }
}