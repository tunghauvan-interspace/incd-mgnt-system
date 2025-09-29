export interface Alert {
  id: string
  fingerprint: string
  status: 'firing' | 'resolved' | 'silenced'
  severity: 'critical' | 'high' | 'medium' | 'low'
  summary: string
  description: string
  source: string
  startsAt: string
  endsAt?: string
  labels: { [key: string]: string }
  annotations: { [key: string]: string }
  generatorURL?: string
  silencedUntil?: string
  incidentId?: string
  createdAt: string
  updatedAt: string
}

export interface AlertGroup {
  id: string
  fingerprint: string
  alerts: Alert[]
  status: Alert['status']
  severity: Alert['severity']
  summary: string
  labels: { [key: string]: string }
  createdAt: string
  updatedAt: string
}

export interface AlertFilter {
  status?: string
  severity?: string
  source?: string
  search?: string
  dateFrom?: string
  dateTo?: string
}

export interface AlertStats {
  total: number
  firing: number
  resolved: number
  silenced: number
  critical: number
}

export interface SilenceAlertData {
  duration: number // in minutes
  reason?: string
}