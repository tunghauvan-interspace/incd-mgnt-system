// Mock data for demonstrating the component architecture
import type { Incident, Alert, Metrics } from '@/types/api'

export const mockIncidents: Incident[] = [
  {
    id: 'inc-001-2024',
    title: 'Database Connection Timeout',
    description: 'Multiple database connection timeouts reported in the production environment.',
    severity: 'critical',
    status: 'open',
    assignee_id: 'user123',
    created_at: '2024-01-15T10:30:00Z',
    updated_at: '2024-01-15T10:30:00Z',
    labels: {
      environment: 'production',
      service: 'database',
      priority: 'urgent'
    }
  },
  {
    id: 'inc-002-2024',
    title: 'API Rate Limit Exceeded',
    description: 'External API calls are being rate limited, affecting user experience.',
    severity: 'high',
    status: 'acknowledged',
    assignee_id: 'user456',
    created_at: '2024-01-15T09:15:00Z',
    updated_at: '2024-01-15T11:00:00Z',
    acknowledged_at: '2024-01-15T11:00:00Z',
    labels: {
      environment: 'production',
      service: 'api',
      component: 'rate-limiter'
    }
  },
  {
    id: 'inc-003-2024',
    title: 'Memory Usage Alert',
    description: 'Application memory usage has exceeded 85% threshold.',
    severity: 'medium',
    status: 'resolved',
    assignee_id: 'user789',
    created_at: '2024-01-14T14:20:00Z',
    updated_at: '2024-01-14T16:45:00Z',
    acknowledged_at: '2024-01-14T14:30:00Z',
    resolved_at: '2024-01-14T16:45:00Z',
    labels: {
      environment: 'production',
      service: 'application',
      type: 'performance'
    }
  },
  {
    id: 'inc-004-2024',
    title: 'SSL Certificate Expiring',
    description: 'SSL certificate for the main domain will expire in 7 days.',
    severity: 'low',
    status: 'open',
    created_at: '2024-01-13T08:00:00Z',
    updated_at: '2024-01-13T08:00:00Z',
    labels: {
      type: 'security',
      component: 'ssl'
    }
  }
]

export const mockAlerts: Alert[] = [
  {
    id: 'alert-001',
    alert_name: 'High CPU Usage',
    generator_url: 'http://prometheus:9090/graph',
    status: 'firing',
    starts_at: '2024-01-15T10:00:00Z',
    labels: {
      severity: 'warning',
      instance: 'web-server-01',
      job: 'web-servers'
    },
    annotations: {
      summary: 'CPU usage is above 90%',
      description: 'The CPU usage on web-server-01 has been above 90% for more than 5 minutes'
    },
    created_at: '2024-01-15T10:00:00Z',
    updated_at: '2024-01-15T10:00:00Z'
  },
  {
    id: 'alert-002',
    alert_name: 'Disk Space Low',
    generator_url: 'http://prometheus:9090/graph',
    status: 'resolved',
    starts_at: '2024-01-14T15:00:00Z',
    ends_at: '2024-01-14T16:30:00Z',
    labels: {
      severity: 'critical',
      instance: 'db-server-01',
      job: 'database-servers'
    },
    annotations: {
      summary: 'Disk space is critically low',
      description: 'Available disk space on db-server-01 is below 5%'
    },
    incident_id: 'inc-002-2024',
    created_at: '2024-01-14T15:00:00Z',
    updated_at: '2024-01-14T16:30:00Z'
  }
]

export const mockMetrics: Metrics = {
  total_incidents: 4,
  open_incidents: 2,
  mtta: 1800000000000, // 30 minutes in nanoseconds
  mttr: 7200000000000, // 2 hours in nanoseconds
  incidents_by_status: {
    open: 2,
    acknowledged: 1,
    resolved: 1
  },
  incidents_by_severity: {
    critical: 1,
    high: 1,
    medium: 1,
    low: 1
  }
}
