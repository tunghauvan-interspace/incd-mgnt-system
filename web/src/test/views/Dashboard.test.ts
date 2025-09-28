import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import Dashboard from '@/views/Dashboard.vue'
import type { Metrics } from '@/types/api'

// Mock the API module
vi.mock('@/services/api', () => ({
  metricsAPI: {
    getMetrics: vi.fn()
  }
}))

// Mock the formatDuration utility
vi.mock('@/utils/format', () => ({
  formatDuration: vi.fn((duration) => `${duration}ms`)
}))

const mockMetrics: Metrics = {
  total_incidents: 100,
  open_incidents: 15,
  mtta: 3600000000000, // 1 hour in nanoseconds
  mttr: 7200000000000, // 2 hours in nanoseconds
  incidents_by_status: { open: 15, acknowledged: 5, resolved: 80 },
  incidents_by_severity: { critical: 3, high: 12, medium: 50, low: 35 }
}

describe('Dashboard', () => {
  // Import the mocked functions after mocking
  let mockGetMetrics: any

  beforeEach(async () => {
    const { metricsAPI } = await import('@/services/api')
    mockGetMetrics = vi.mocked(metricsAPI.getMetrics)
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('renders dashboard header correctly', () => {
    mockGetMetrics.mockResolvedValueOnce(mockMetrics)

    const wrapper = mount(Dashboard, {
      global: {
        stubs: {
          DoughnutChart: {
            template: '<div class="mock-chart">{{ title }}</div>',
            props: ['data', 'title']
          }
        }
      }
    })

    expect(wrapper.find('h2').text()).toBe('Dashboard')
    expect(wrapper.find('.refresh-btn button').exists()).toBe(true)
  })

  it('shows loading state initially', async () => {
    let resolvePromise: (value: Metrics) => void
    const promise = new Promise<Metrics>((resolve) => {
      resolvePromise = resolve
    })
    mockGetMetrics.mockReturnValueOnce(promise)

    const wrapper = mount(Dashboard, {
      global: {
        stubs: {
          DoughnutChart: {
            template: '<div class="mock-chart">{{ title }}</div>',
            props: ['data', 'title']
          }
        }
      }
    })

    expect(wrapper.find('.loading').exists()).toBe(true)
    expect(wrapper.text()).toContain('Loading dashboard...')
    expect(wrapper.find('button').attributes('disabled')).toBeDefined()
    expect(wrapper.find('button').text()).toBe('Loading...')

    // Resolve the promise to avoid hanging test
    resolvePromise!(mockMetrics)
    await nextTick()
  })

  it('displays metrics correctly after loading', async () => {
    mockGetMetrics.mockResolvedValueOnce(mockMetrics)

    const wrapper = mount(Dashboard, {
      global: {
        stubs: {
          DoughnutChart: {
            template: '<div class="mock-chart">{{ title }}</div>',
            props: ['data', 'title']
          }
        }
      }
    })

    // Wait for the component to finish loading
    await new Promise((resolve) => setTimeout(resolve, 10))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.loading').exists()).toBe(false)
    expect(wrapper.find('.metrics-grid').exists()).toBe(true)

    const metricCards = wrapper.findAll('.metric-card')
    expect(metricCards).toHaveLength(4)

    // Check Total Incidents
    expect(metricCards[0].text()).toContain('Total Incidents')
    expect(metricCards[0].text()).toContain('100')

    // Check Open Incidents
    expect(metricCards[1].text()).toContain('Open Incidents')
    expect(metricCards[1].text()).toContain('15')
    expect(metricCards[1].find('.metric-value').classes()).toContain('critical')

    // Check MTTA
    expect(metricCards[2].text()).toContain('MTTA')
    expect(metricCards[2].text()).toContain('Mean Time To Acknowledge')

    // Check MTTR
    expect(metricCards[3].text()).toContain('MTTR')
    expect(metricCards[3].text()).toContain('Mean Time To Resolve')
  })

  it('shows error message when API fails', async () => {
    const errorMessage = 'Network Error'
    mockGetMetrics.mockRejectedValueOnce(new Error(errorMessage))

    const wrapper = mount(Dashboard, {
      global: {
        stubs: {
          DoughnutChart: {
            template: '<div class="mock-chart">{{ title }}</div>',
            props: ['data', 'title']
          }
        }
      }
    })

    // Wait for the error to be handled
    await new Promise((resolve) => setTimeout(resolve, 10))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.error-message').exists()).toBe(true)
    expect(wrapper.text()).toContain('Error loading dashboard data')
    expect(wrapper.find('.metrics-grid').exists()).toBe(false)
  })

  it('calls API on mount', () => {
    mockGetMetrics.mockResolvedValueOnce(mockMetrics)

    mount(Dashboard, {
      global: {
        stubs: {
          DoughnutChart: {
            template: '<div class="mock-chart">{{ title }}</div>',
            props: ['data', 'title']
          }
        }
      }
    })

    expect(mockGetMetrics).toHaveBeenCalledTimes(1)
  })

  it('refreshes dashboard when refresh button is clicked', async () => {
    mockGetMetrics.mockResolvedValue(mockMetrics)

    const wrapper = mount(Dashboard, {
      global: {
        stubs: {
          DoughnutChart: {
            template: '<div class="mock-chart">{{ title }}</div>',
            props: ['data', 'title']
          }
        }
      }
    })

    // Wait for initial load
    await new Promise((resolve) => setTimeout(resolve, 10))
    await wrapper.vm.$nextTick()

    // Clear the first call from mount
    vi.clearAllMocks()

    await wrapper.find('.refresh-btn button').trigger('click')
    await nextTick()

    expect(mockGetMetrics).toHaveBeenCalledTimes(1)
  })

  it('handles empty metrics gracefully', async () => {
    const emptyMetrics: Metrics = {
      total_incidents: 0,
      open_incidents: 0,
      mtta: 0,
      mttr: 0,
      incidents_by_status: {},
      incidents_by_severity: {}
    }

    mockGetMetrics.mockResolvedValueOnce(emptyMetrics)

    const wrapper = mount(Dashboard, {
      global: {
        stubs: {
          DoughnutChart: {
            template: '<div class="mock-chart">{{ title }}</div>',
            props: ['data', 'title']
          }
        }
      }
    })

    // Wait for loading to complete
    await new Promise((resolve) => setTimeout(resolve, 10))
    await wrapper.vm.$nextTick()

    const metricCards = wrapper.findAll('.metric-card')
    expect(metricCards[0].text()).toContain('0') // Total incidents
    expect(metricCards[1].text()).toContain('0') // Open incidents
  })

  it('shows charts section when metrics are loaded', async () => {
    mockGetMetrics.mockResolvedValueOnce(mockMetrics)

    const wrapper = mount(Dashboard, {
      global: {
        stubs: {
          DoughnutChart: {
            template: '<div class="mock-chart">{{ title }}</div>',
            props: ['data', 'title']
          }
        }
      }
    })

    // Wait for loading to complete
    await new Promise((resolve) => setTimeout(resolve, 10))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.charts-section').exists()).toBe(true)
    expect(wrapper.find('.charts-grid').exists()).toBe(true)
  })

  it('does not show charts section when loading', async () => {
    let resolvePromise: (value: Metrics) => void
    const promise = new Promise<Metrics>((resolve) => {
      resolvePromise = resolve
    })
    mockGetMetrics.mockReturnValueOnce(promise)

    const wrapper = mount(Dashboard, {
      global: {
        stubs: {
          DoughnutChart: {
            template: '<div class="mock-chart">{{ title }}</div>',
            props: ['data', 'title']
          }
        }
      }
    })

    expect(wrapper.find('.charts-section').exists()).toBe(false)

    // Clean up
    resolvePromise!(mockMetrics)
    await nextTick()
  })

  it('does not show charts section when there is an error', async () => {
    mockGetMetrics.mockRejectedValueOnce(new Error('Network Error'))

    const wrapper = mount(Dashboard, {
      global: {
        stubs: {
          DoughnutChart: {
            template: '<div class="mock-chart">{{ title }}</div>',
            props: ['data', 'title']
          }
        }
      }
    })

    // Wait for error handling
    await new Promise((resolve) => setTimeout(resolve, 10))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.charts-section').exists()).toBe(false)
  })
})
