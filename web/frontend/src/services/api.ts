import axios from 'axios'
import type { Incident, Alert, Metrics, AcknowledgeIncidentRequest } from '@/types/api'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor for logging
api.interceptors.request.use(
  (config) => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`)
    return config
  },
  (error) => {
    console.error('API Request Error:', error)
    return Promise.reject(error)
  }
)

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    console.error('API Response Error:', error)
    return Promise.reject(error)
  }
)

export const incidentAPI = {
  // Get all incidents
  getIncidents: (): Promise<Incident[]> =>
    api.get<Incident[]>('/incidents').then((res) => res.data),

  // Get incident by ID
  getIncident: (id: string): Promise<Incident> =>
    api.get<Incident>(`/incidents/${id}`).then((res) => res.data),

  // Acknowledge incident
  acknowledgeIncident: (id: string, data: AcknowledgeIncidentRequest): Promise<void> =>
    api.put(`/incidents/${id}/acknowledge`, data).then(() => {}),

  // Resolve incident
  resolveIncident: (id: string): Promise<void> => api.put(`/incidents/${id}/resolve`).then(() => {})
}

export const alertAPI = {
  // Get all alerts
  getAlerts: (): Promise<Alert[]> => api.get<Alert[]>('/alerts').then((res) => res.data)
}

export const metricsAPI = {
  // Get metrics
  getMetrics: (): Promise<Metrics> => api.get<Metrics>('/metrics').then((res) => res.data)
}

export default api
