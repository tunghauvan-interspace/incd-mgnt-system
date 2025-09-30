import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig } from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useNotificationsStore } from '@/stores/notifications'

let apiClient: AxiosInstance | null = null

export const useApi = () => {
  if (!apiClient) {
    apiClient = createApiClient()
  }
  
  return {
    client: apiClient
  }
}

const createApiClient = (): AxiosInstance => {
  const client = axios.create({
    baseURL: '/api',
    headers: {
      'Content-Type': 'application/json'
    },
    timeout: 10000
  })

  // Request interceptor for auth
  client.interceptors.request.use((config) => {
    const authStore = useAuthStore()
    
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    
    return config
  }, (error) => {
    return Promise.reject(error)
  })

  // Response interceptor for error handling
  client.interceptors.response.use(
    (response) => {
      return response
    },
    (error) => {
      const notificationsStore = useNotificationsStore()
      
      // Handle 401 Unauthorized
      if (error.response?.status === 401) {
        const authStore = useAuthStore()
        authStore.logout()
        window.location.href = '/auth/login'
        return Promise.reject(error)
      }

      // Handle other HTTP errors
      if (error.response) {
        const message = error.response.data?.message || `HTTP ${error.response.status}: ${error.response.statusText}`
        notificationsStore.error(message, 'API Error')
      } else if (error.request) {
        notificationsStore.error('Network error - please check your connection', 'Connection Error')
      } else {
        notificationsStore.error('An unexpected error occurred', 'Error')
      }

      return Promise.reject(error)
    }
  )

  return client
}

// Helper functions for common API operations
export const apiGet = <T = any>(url: string, config?: AxiosRequestConfig) => {
  const { client } = useApi()
  return client.get<T>(url, config)
}

export const apiPost = <T = any>(url: string, data?: any, config?: AxiosRequestConfig) => {
  const { client } = useApi()
  return client.post<T>(url, data, config)
}

export const apiPut = <T = any>(url: string, data?: any, config?: AxiosRequestConfig) => {
  const { client } = useApi()
  return client.put<T>(url, data, config)
}

export const apiDelete = <T = any>(url: string, config?: AxiosRequestConfig) => {
  const { client } = useApi()
  return client.delete<T>(url, config)
}