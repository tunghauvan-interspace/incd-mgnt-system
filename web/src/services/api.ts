import axios from 'axios'
import type {
  Incident,
  Alert,
  Metrics,
  AcknowledgeIncidentRequest,
  LoginRequest,
  RegisterRequest,
  AuthResponse
} from '@/types/api'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor for logging and auth
api.interceptors.request.use(
  (config) => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`)

    // Add auth token if available
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

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

export const authAPI = {
  // Login
  login: (data: LoginRequest): Promise<AuthResponse> =>
    api.post<AuthResponse>('/auth/login', data).then((res) => res.data),

  // Register
  register: (data: RegisterRequest): Promise<AuthResponse> =>
    api.post<AuthResponse>('/auth/register', data).then((res) => res.data),

  // Refresh token
  refreshToken: (refreshToken: string): Promise<AuthResponse> =>
    api.post<AuthResponse>('/auth/refresh', { refresh_token: refreshToken }).then((res) => res.data)
}

export default api
