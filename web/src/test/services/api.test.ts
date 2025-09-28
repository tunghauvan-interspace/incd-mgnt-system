import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import axios from 'axios'
import type { Incident, Alert, Metrics } from '@/types/api'

// Mock axios entirely first
vi.mock('axios', () => ({
  default: {
    create: vi.fn()
  }
}))

const mockAxiosInstance = {
  get: vi.fn(),
  put: vi.fn(),
  post: vi.fn(),
  delete: vi.fn(),
  interceptors: {
    request: { use: vi.fn() },
    response: { use: vi.fn() }
  }
}

describe('API Services', () => {
  beforeEach(() => {
    vi.mocked(axios.create).mockReturnValue(mockAxiosInstance as any)
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('incidentAPI', () => {
    it('getIncidents should fetch all incidents', async () => {
      const mockIncidents: Incident[] = [
        {
          id: '1',
          title: 'Test Incident',
          description: 'Test Description',
          severity: 'high',
          status: 'open',
          created_at: '2023-01-01T00:00:00Z',
          updated_at: '2023-01-01T00:00:00Z'
        }
      ]

      mockAxiosInstance.get.mockResolvedValueOnce({ data: mockIncidents })

      // Import after mocking
      const { incidentAPI } = await import('@/services/api')
      const result = await incidentAPI.getIncidents()

      expect(mockAxiosInstance.get).toHaveBeenCalledWith('/incidents')
      expect(result).toEqual(mockIncidents)
    })

    it('getIncident should fetch a specific incident', async () => {
      const mockIncident: Incident = {
        id: '1',
        title: 'Test Incident',
        description: 'Test Description',
        severity: 'high',
        status: 'open',
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-01-01T00:00:00Z'
      }

      mockAxiosInstance.get.mockResolvedValueOnce({ data: mockIncident })

      const { incidentAPI } = await import('@/services/api')
      const result = await incidentAPI.getIncident('1')

      expect(mockAxiosInstance.get).toHaveBeenCalledWith('/incidents/1')
      expect(result).toEqual(mockIncident)
    })

    it('acknowledgeIncident should acknowledge an incident', async () => {
      const acknowledgeData = { assignee_id: 'user123' }
      mockAxiosInstance.put.mockResolvedValueOnce({})

      const { incidentAPI } = await import('@/services/api')
      await incidentAPI.acknowledgeIncident('1', acknowledgeData)

      expect(mockAxiosInstance.put).toHaveBeenCalledWith('/incidents/1/acknowledge', acknowledgeData)
    })

    it('resolveIncident should resolve an incident', async () => {
      mockAxiosInstance.put.mockResolvedValueOnce({})

      const { incidentAPI } = await import('@/services/api')
      await incidentAPI.resolveIncident('1')

      expect(mockAxiosInstance.put).toHaveBeenCalledWith('/incidents/1/resolve')
    })
  })

  describe('alertAPI', () => {
    it('getAlerts should fetch all alerts', async () => {
      const mockAlerts: Alert[] = [
        {
          id: '1',
          alert_name: 'Test Alert',
          generator_url: 'http://example.com',
          status: 'active',
          starts_at: '2023-01-01T00:00:00Z',
          created_at: '2023-01-01T00:00:00Z',
          updated_at: '2023-01-01T00:00:00Z'
        }
      ]

      mockAxiosInstance.get.mockResolvedValueOnce({ data: mockAlerts })

      const { alertAPI } = await import('@/services/api')
      const result = await alertAPI.getAlerts()

      expect(mockAxiosInstance.get).toHaveBeenCalledWith('/alerts')
      expect(result).toEqual(mockAlerts)
    })
  })

  describe('metricsAPI', () => {
    it('getMetrics should fetch metrics', async () => {
      const mockMetrics: Metrics = {
        total_incidents: 10,
        open_incidents: 5,
        mtta: 3600000000000, // 1 hour in nanoseconds
        mttr: 7200000000000, // 2 hours in nanoseconds
        incidents_by_status: { open: 5, closed: 5 },
        incidents_by_severity: { high: 3, medium: 4, low: 3 }
      }

      mockAxiosInstance.get.mockResolvedValueOnce({ data: mockMetrics })

      const { metricsAPI } = await import('@/services/api')
      const result = await metricsAPI.getMetrics()

      expect(mockAxiosInstance.get).toHaveBeenCalledWith('/metrics')
      expect(result).toEqual(mockMetrics)
    })
  })

  describe('Error handling', () => {
    it('should handle API errors gracefully', async () => {
      const errorResponse = new Error('Network Error')
      mockAxiosInstance.get.mockRejectedValueOnce(errorResponse)

      const { incidentAPI } = await import('@/services/api')
      await expect(incidentAPI.getIncidents()).rejects.toThrow('Network Error')
    })
  })
})