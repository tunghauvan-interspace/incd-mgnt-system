// Component prop interfaces and type definitions

import type { Ref } from 'vue'
import type { Incident, Alert, Metrics } from './api'

// Button component types
export interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  type?: 'button' | 'submit' | 'reset'
  block?: boolean
}

// Modal component types
export interface ModalProps {
  show: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  closable?: boolean
  backdropClosable?: boolean
}

// Table component types
export interface TableColumn<T = Record<string, any>> {
  key: keyof T
  label: string
  sortable?: boolean
  width?: string
  align?: 'left' | 'center' | 'right'
  formatter?: (value: any, row: T) => string
}

export interface TableProps<T = Record<string, any>> {
  columns: TableColumn<T>[]
  data: T[]
  loading?: boolean
  emptyMessage?: string
  sortBy?: keyof T | null
  sortOrder?: 'asc' | 'desc'
  hoverable?: boolean
  striped?: boolean
  bordered?: boolean
  compact?: boolean
}

// Badge component types
export interface StatusBadgeProps {
  status: 'open' | 'acknowledged' | 'resolved' | string
  size?: 'sm' | 'md' | 'lg'
}

export interface SeverityBadgeProps {
  severity: 'critical' | 'high' | 'medium' | 'low' | 'info' | string
  size?: 'sm' | 'md' | 'lg'
  showIcon?: boolean
}

// Navbar component types
export interface NavItem {
  name: string
  path: string
  label: string
  icon?: string
}

// Chart component types
export interface ChartData {
  labels: string[]
  datasets: {
    data: number[]
    backgroundColor?: string[]
    borderColor?: string[]
    borderWidth?: number
  }[]
}

export interface DoughnutChartProps {
  data: ChartData
  title?: string
  height?: number
  width?: number
}

// Form component types
export interface FormField {
  name: string
  label: string
  type: 'text' | 'email' | 'password' | 'textarea' | 'select' | 'checkbox' | 'radio'
  required?: boolean
  placeholder?: string
  options?: { value: string; label: string }[]
  validation?: {
    pattern?: RegExp
    message?: string
    minLength?: number
    maxLength?: number
  }
}

// Composable return types
export interface UseApiState<T> {
  data: Ref<T | null>
  loading: Ref<boolean>
  error: Ref<string | null>
  refresh: () => Promise<void>
}

export interface UseIncidentsReturn extends UseApiState<Incident[]> {
  acknowledgeIncident: (id: string, assigneeId: string) => Promise<void>
  resolveIncident: (id: string) => Promise<void>
  getIncident: (id: string) => Promise<Incident | null>
}

export interface UseAlertsReturn extends UseApiState<Alert[]> {
  getAlert: (id: string) => Promise<Alert | null>
}

export interface UseMetricsReturn extends UseApiState<Metrics> {
  // Metrics-specific methods can be added here
}

// Event handler types
export interface TableRowClickEvent<T> {
  row: T
  index: number
}

export interface SortEvent {
  column: string
  order: 'asc' | 'desc'
}

// Filter types
export interface IncidentFilter {
  status?: string[]
  severity?: string[]
  assignee?: string
  dateRange?: {
    start: Date
    end: Date
  }
}

export interface AlertFilter {
  status?: string[]
  labels?: Record<string, string>
  dateRange?: {
    start: Date
    end: Date
  }
}

// Utility types
export type EmitEvents<T = Record<string, any>> = {
  [K in keyof T]: T[K] extends (...args: infer Args) => any ? Args : never
}

// Component instance types (for template refs) - simplified
// export type ButtonInstance = InstanceType<typeof import('@/components/Button.vue')['default']>
// export type ModalInstance = InstanceType<typeof import('@/components/Modal.vue')['default']>
// export type TableInstance = InstanceType<typeof import('@/components/DataTable.vue')['default']>

// Re-export API types for convenience
export type { Incident, Alert, Metrics, AcknowledgeIncidentRequest } from './api'
