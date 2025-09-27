import { ref, type Ref } from 'vue'
import { alertAPI } from '@/services/api'
import type { Alert } from '@/types/api'

export interface UseAlertsReturn {
  alerts: Ref<Alert[]>
  loading: Ref<boolean>
  error: Ref<string | null>
  loadAlerts: () => Promise<void>
  getAlert: (id: string) => Promise<Alert | null>
  refresh: () => Promise<void>
}

/**
 * Composable for managing alerts data and operations
 * Provides reactive state management and API operations for alerts
 */
export function useAlerts(): UseAlertsReturn {
  const alerts = ref<Alert[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  /**
   * Load all alerts from the API
   */
  const loadAlerts = async (): Promise<void> => {
    try {
      loading.value = true
      error.value = null

      const data = await alertAPI.getAlerts()

      // Sort alerts by creation date (newest first)
      alerts.value = data.sort(
        (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      )
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load alerts'
      console.error('Error loading alerts:', err)
      alerts.value = []
    } finally {
      loading.value = false
    }
  }

  /**
   * Get a specific alert by ID (simulated - API doesn't have this endpoint yet)
   */
  const getAlert = async (id: string): Promise<Alert | null> => {
    try {
      error.value = null

      // Since the API doesn't have a getAlert endpoint, find in local data
      const alert = alerts.value.find((a) => a.id === id)
      if (alert) {
        return alert
      }

      // If not found locally, try to load all alerts and find it
      await loadAlerts()
      return alerts.value.find((a) => a.id === id) || null
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to get alert'
      console.error('Error getting alert:', err)
      return null
    }
  }

  /**
   * Refresh alerts data (alias for loadAlerts)
   */
  const refresh = loadAlerts

  return {
    alerts,
    loading,
    error,
    loadAlerts,
    getAlert,
    refresh
  }
}
