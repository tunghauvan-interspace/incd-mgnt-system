import { ref, type Ref } from 'vue'
import { metricsAPI } from '@/services/api'
import type { Metrics } from '@/types/api'

export interface UseMetricsReturn {
  metrics: Ref<Metrics | null>
  loading: Ref<boolean>
  error: Ref<string | null>
  loadMetrics: () => Promise<void>
  refresh: () => Promise<void>
}

/**
 * Composable for managing dashboard metrics data
 * Provides reactive state management and API operations for metrics
 */
export function useMetrics(): UseMetricsReturn {
  const metrics = ref<Metrics | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  /**
   * Load metrics from the API
   */
  const loadMetrics = async (): Promise<void> => {
    try {
      loading.value = true
      error.value = null

      const data = await metricsAPI.getMetrics()
      metrics.value = data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load metrics'
      console.error('Error loading metrics:', err)
      metrics.value = null
    } finally {
      loading.value = false
    }
  }

  /**
   * Refresh metrics data (alias for loadMetrics)
   */
  const refresh = loadMetrics

  return {
    metrics,
    loading,
    error,
    loadMetrics,
    refresh
  }
}
