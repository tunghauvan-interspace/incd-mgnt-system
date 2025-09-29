// API endpoints
export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/auth/login',
    REGISTER: '/auth/register',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    ME: '/auth/me'
  },
  INCIDENTS: {
    LIST: '/incidents',
    CREATE: '/incidents',
    GET: (id: string) => `/incidents/${id}`,
    UPDATE: (id: string) => `/incidents/${id}`,
    DELETE: (id: string) => `/incidents/${id}`,
    TIMELINE: (id: string) => `/incidents/${id}/timeline`,
    ASSIGN: (id: string) => `/incidents/${id}/assign`,
    STATUS: (id: string) => `/incidents/${id}/status`
  },
  ALERTS: {
    LIST: '/alerts',
    GET: (id: string) => `/alerts/${id}`,
    ACKNOWLEDGE: (id: string) => `/alerts/${id}/acknowledge`,
    SILENCE: (id: string) => `/alerts/${id}/silence`,
    RESOLVE: (id: string) => `/alerts/${id}/resolve`
  },
  USERS: {
    LIST: '/users',
    CREATE: '/users',
    GET: (id: string) => `/users/${id}`,
    UPDATE: (id: string) => `/users/${id}`,
    DELETE: (id: string) => `/users/${id}`
  },
  NOTIFICATIONS: {
    CHANNELS: '/notifications/channels',
    RULES: '/notifications/rules',
    SEND: '/notifications/send'
  }
} as const

// Status options
export const INCIDENT_STATUSES = [
  { value: 'open', label: 'Open', color: 'red' },
  { value: 'acknowledged', label: 'Acknowledged', color: 'yellow' },
  { value: 'investigating', label: 'Investigating', color: 'blue' },
  { value: 'resolved', label: 'Resolved', color: 'green' },
  { value: 'closed', label: 'Closed', color: 'gray' }
] as const

export const SEVERITY_LEVELS = [
  { value: 'critical', label: 'Critical', color: 'red' },
  { value: 'high', label: 'High', color: 'orange' },
  { value: 'medium', label: 'Medium', color: 'yellow' },
  { value: 'low', label: 'Low', color: 'green' }
] as const

export const ALERT_STATUSES = [
  { value: 'firing', label: 'Firing', color: 'red' },
  { value: 'resolved', label: 'Resolved', color: 'green' },
  { value: 'silenced', label: 'Silenced', color: 'gray' }
] as const

export const USER_ROLES = [
  { value: 'admin', label: 'Administrator' },
  { value: 'user', label: 'User' },
  { value: 'viewer', label: 'Viewer' }
] as const

// Time constants
export const TIME_RANGES = {
  LAST_HOUR: 'last_hour',
  LAST_24_HOURS: 'last_24_hours',
  LAST_7_DAYS: 'last_7_days',
  LAST_30_DAYS: 'last_30_days',
  CUSTOM: 'custom'
} as const

// Pagination
export const DEFAULT_PAGE_SIZE = 20
export const PAGE_SIZE_OPTIONS = [10, 20, 50, 100]

// Chart colors
export const CHART_COLORS = {
  PRIMARY: '#3b82f6',
  SUCCESS: '#10b981',
  WARNING: '#f59e0b',
  ERROR: '#ef4444',
  INFO: '#6366f1',
  CRITICAL: '#dc2626',
  HIGH: '#ea580c',
  MEDIUM: '#d97706',
  LOW: '#65a30d'
} as const

// Local storage keys
export const STORAGE_KEYS = {
  AUTH_TOKEN: 'auth_token',
  USER_PREFERENCES: 'user_preferences',
  THEME: 'theme',
  SIDEBAR_STATE: 'sidebar_state'
} as const