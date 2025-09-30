import { ref, computed } from 'vue'
import { useIncidentsStore } from '@/stores/incidents'
import type { Incident, CreateIncidentData, UpdateIncidentData } from '@/types/incidents'

export const useIncidents = () => {
  const incidentsStore = useIncidentsStore()
  
  // Local state
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed properties
  const incidents = computed(() => incidentsStore.filteredIncidents)
  const currentIncident = computed(() => incidentsStore.currentIncident)
  const filters = computed(() => incidentsStore.filters)
  const criticalIncidents = computed(() => incidentsStore.criticalIncidents)
  const openIncidents = computed(() => incidentsStore.openIncidents)

  // Actions
  const fetchIncidents = async () => {
    isLoading.value = true
    error.value = null
    
    try {
      await incidentsStore.fetchIncidents()
    } catch (err) {
      error.value = 'Failed to fetch incidents'
      console.error('Error fetching incidents:', err)
    } finally {
      isLoading.value = false
    }
  }

  const fetchIncident = async (id: string) => {
    isLoading.value = true
    error.value = null
    
    try {
      return await incidentsStore.fetchIncident(id)
    } catch (err) {
      error.value = 'Failed to fetch incident'
      console.error('Error fetching incident:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const createIncident = async (data: CreateIncidentData) => {
    isLoading.value = true
    error.value = null
    
    try {
      const incident = await incidentsStore.createIncident(data)
      return incident
    } catch (err) {
      error.value = 'Failed to create incident'
      console.error('Error creating incident:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const updateIncident = async (id: string, data: UpdateIncidentData) => {
    isLoading.value = true
    error.value = null
    
    try {
      const incident = await incidentsStore.updateIncident(id, data)
      return incident
    } catch (err) {
      error.value = 'Failed to update incident'
      console.error('Error updating incident:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const deleteIncident = async (id: string) => {
    isLoading.value = true
    error.value = null
    
    try {
      await incidentsStore.deleteIncident(id)
    } catch (err) {
      error.value = 'Failed to delete incident'
      console.error('Error deleting incident:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Filter functions
  const setFilters = (newFilters: Partial<typeof filters.value>) => {
    incidentsStore.setFilters(newFilters)
  }

  const clearFilters = () => {
    incidentsStore.clearFilters()
  }

  const clearError = () => {
    error.value = null
  }

  return {
    // State
    isLoading,
    error,
    // Computed
    incidents,
    currentIncident,
    filters,
    criticalIncidents,
    openIncidents,
    // Actions
    fetchIncidents,
    fetchIncident,
    createIncident,
    updateIncident,
    deleteIncident,
    setFilters,
    clearFilters,
    clearError
  }
}