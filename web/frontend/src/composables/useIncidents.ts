import { ref, type Ref } from 'vue'
import { incidentAPI } from '@/services/api'
import { mockIncidents } from '@/utils/mockData'
import type { Incident, AcknowledgeIncidentRequest } from '@/types/api'

export interface UseIncidentsReturn {
  incidents: Ref<Incident[]>
  loading: Ref<boolean>
  error: Ref<string | null>
  loadIncidents: () => Promise<void>
  getIncident: (id: string) => Promise<Incident | null>
  acknowledgeIncident: (id: string, assigneeId: string) => Promise<void>
  resolveIncident: (id: string) => Promise<void>
  refresh: () => Promise<void>
}

/**
 * Composable for managing incidents data and operations
 * Provides reactive state management and API operations for incidents
 */
export function useIncidents(): UseIncidentsReturn {
  const incidents = ref<Incident[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  /**
   * Load all incidents from the API (with fallback to mock data for demo)
   */
  const loadIncidents = async (): Promise<void> => {
    try {
      loading.value = true
      error.value = null

      const data = await incidentAPI.getIncidents()

      // Sort incidents by creation date (newest first)
      incidents.value = data.sort(
        (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      )
    } catch (err) {
      // Use mock data for demo purposes when API fails
      console.warn('API failed, using mock data for demo:', err)
      incidents.value = mockIncidents.sort(
        (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      )
      error.value = null // Clear error since we have fallback data
    } finally {
      loading.value = false
    }
  }

  /**
   * Get a specific incident by ID
   */
  const getIncident = async (id: string): Promise<Incident | null> => {
    try {
      error.value = null
      return await incidentAPI.getIncident(id)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to get incident'
      console.error('Error getting incident:', err)
      return null
    }
  }

  /**
   * Acknowledge an incident
   */
  const acknowledgeIncident = async (id: string, assigneeId: string): Promise<void> => {
    try {
      error.value = null

      const request: AcknowledgeIncidentRequest = {
        assignee_id: assigneeId
      }

      await incidentAPI.acknowledgeIncident(id, request)

      // Update local state
      const incidentIndex = incidents.value.findIndex((incident) => incident.id === id)
      if (incidentIndex !== -1) {
        incidents.value[incidentIndex] = {
          ...incidents.value[incidentIndex],
          status: 'acknowledged',
          assignee_id: assigneeId,
          acknowledged_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        }
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to acknowledge incident'
      console.error('Error acknowledging incident:', err)
      throw err
    }
  }

  /**
   * Resolve an incident
   */
  const resolveIncident = async (id: string): Promise<void> => {
    try {
      error.value = null

      await incidentAPI.resolveIncident(id)

      // Update local state
      const incidentIndex = incidents.value.findIndex((incident) => incident.id === id)
      if (incidentIndex !== -1) {
        incidents.value[incidentIndex] = {
          ...incidents.value[incidentIndex],
          status: 'resolved',
          resolved_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        }
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to resolve incident'
      console.error('Error resolving incident:', err)
      throw err
    }
  }

  /**
   * Refresh incidents data (alias for loadIncidents)
   */
  const refresh = loadIncidents

  return {
    incidents,
    loading,
    error,
    loadIncidents,
    getIncident,
    acknowledgeIncident,
    resolveIncident,
    refresh
  }
}
