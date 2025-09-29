export interface Notification {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  title: string
  message: string
  timestamp: Date
  persistent?: boolean
  action?: {
    label: string
    handler: () => void
  }
}

export interface NotificationChannel {
  id: string
  name: string
  type: 'slack' | 'email' | 'telegram' | 'webhook'
  enabled: boolean
  config: { [key: string]: any }
}

export interface NotificationRule {
  id: string
  name: string
  enabled: boolean
  triggers: {
    incidentCreated: boolean
    incidentResolved: boolean
    alertFiring: boolean
    alertResolved: boolean
  }
  channels: string[]
  filters: {
    severity?: string[]
    labels?: { [key: string]: string }
  }
}