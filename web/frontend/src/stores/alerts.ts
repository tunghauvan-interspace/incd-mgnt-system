import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Alert } from '@/types/alerts'
import { useApi } from '@/composables/useApi'

export const useAlertsStore = defineStore('alerts', () => {
  const { client } = useApi()

  // State
  const alerts = ref<Alert[]>([])
  const currentAlert = ref<Alert | null>(null)
  const isLoading = ref(false)

  // Getters
  const firingAlerts = computed(() => 
    alerts.value.filter(alert => alert.status === 'firing')
  )

  const resolvedAlerts = computed(() => 
    alerts.value.filter(alert => alert.status === 'resolved')
  )

  const criticalAlerts = computed(() => 
    alerts.value.filter(alert => alert.severity === 'critical')
  )

  // Actions
  const fetchAlerts = async () => {
    isLoading.value = true
    try {
      const response = await client.get('/alerts')
      alerts.value = response.data
    } catch (error) {
      console.error('Failed to fetch alerts:', error)
    } finally {
      isLoading.value = false
    }
  }

  const fetchAlert = async (id: string) => {
    isLoading.value = true
    try {
      const response = await client.get(`/alerts/${id}`)
      currentAlert.value = response.data
      return response.data
    } catch (error) {
      console.error('Failed to fetch alert:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const acknowledgeAlert = async (id: string) => {
    try {
      const response = await client.post(`/alerts/${id}/acknowledge`)
      const index = alerts.value.findIndex(a => a.id === id)
      if (index > -1) {
        alerts.value[index] = response.data
      }
      if (currentAlert.value?.id === id) {
        currentAlert.value = response.data
      }
      return response.data
    } catch (error) {
      console.error('Failed to acknowledge alert:', error)
      throw error
    }
  }

  const silenceAlert = async (id: string, duration: number) => {
    try {
      const response = await client.post(`/alerts/${id}/silence`, { duration })
      const index = alerts.value.findIndex(a => a.id === id)
      if (index > -1) {
        alerts.value[index] = response.data
      }
      if (currentAlert.value?.id === id) {
        currentAlert.value = response.data
      }
      return response.data
    } catch (error) {
      console.error('Failed to silence alert:', error)
      throw error
    }
  }

  return {
    // State
    alerts,
    currentAlert,
    isLoading,
    // Getters
    firingAlerts,
    resolvedAlerts,
    criticalAlerts,
    // Actions
    fetchAlerts,
    fetchAlert,
    acknowledgeAlert,
    silenceAlert
  }
})