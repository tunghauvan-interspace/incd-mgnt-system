import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Incident } from '@/types/incidents'
import { useApi } from '@/composables/useApi'

export const useIncidentsStore = defineStore('incidents', () => {
  const { client } = useApi()

  // State
  const incidents = ref<Incident[]>([])
  const currentIncident = ref<Incident | null>(null)
  const isLoading = ref(false)
  const filters = ref({
    status: '',
    severity: '',
    assignee: '',
    search: ''
  })

  // Getters
  const filteredIncidents = computed(() => {
    return incidents.value.filter(incident => {
      if (filters.value.status && incident.status !== filters.value.status) return false
      if (filters.value.severity && incident.severity !== filters.value.severity) return false
      if (filters.value.assignee && incident.assignee?.id !== filters.value.assignee) return false
      if (filters.value.search && !incident.title.toLowerCase().includes(filters.value.search.toLowerCase())) return false
      return true
    })
  })

  const criticalIncidents = computed(() => 
    incidents.value.filter(i => i.severity === 'critical')
  )

  const openIncidents = computed(() => 
    incidents.value.filter(i => i.status === 'open')
  )

  // Actions
  const fetchIncidents = async () => {
    isLoading.value = true
    try {
      const response = await client.get('/incidents')
      incidents.value = response.data
    } catch (error) {
      console.error('Failed to fetch incidents:', error)
    } finally {
      isLoading.value = false
    }
  }

  const fetchIncident = async (id: string) => {
    isLoading.value = true
    try {
      const response = await client.get(`/incidents/${id}`)
      currentIncident.value = response.data
      return response.data
    } catch (error) {
      console.error('Failed to fetch incident:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const createIncident = async (data: Partial<Incident>) => {
    try {
      const response = await client.post('/incidents', data)
      incidents.value.unshift(response.data)
      return response.data
    } catch (error) {
      console.error('Failed to create incident:', error)
      throw error
    }
  }

  const updateIncident = async (id: string, data: Partial<Incident>) => {
    try {
      const response = await client.put(`/incidents/${id}`, data)
      const index = incidents.value.findIndex(i => i.id === id)
      if (index > -1) {
        incidents.value[index] = response.data
      }
      if (currentIncident.value?.id === id) {
        currentIncident.value = response.data
      }
      return response.data
    } catch (error) {
      console.error('Failed to update incident:', error)
      throw error
    }
  }

  const deleteIncident = async (id: string) => {
    try {
      await client.delete(`/incidents/${id}`)
      incidents.value = incidents.value.filter(i => i.id !== id)
      if (currentIncident.value?.id === id) {
        currentIncident.value = null
      }
    } catch (error) {
      console.error('Failed to delete incident:', error)
      throw error
    }
  }

  const setFilters = (newFilters: Partial<typeof filters.value>) => {
    filters.value = { ...filters.value, ...newFilters }
  }

  const clearFilters = () => {
    filters.value = {
      status: '',
      severity: '',
      assignee: '',
      search: ''
    }
  }

  return {
    // State
    incidents,
    currentIncident,
    isLoading,
    filters,
    // Getters
    filteredIncidents,
    criticalIncidents,
    openIncidents,
    // Actions
    fetchIncidents,
    fetchIncident,
    createIncident,
    updateIncident,
    deleteIncident,
    setFilters,
    clearFilters
  }
})