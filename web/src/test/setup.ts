import { vi } from 'vitest'
import type { config } from '@vue/test-utils'

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn()
  }))
})

// Mock CSS variables
Object.defineProperty(window, 'getComputedStyle', {
  value: () => ({
    getPropertyValue: (prop: string) => {
      const cssVars: Record<string, string> = {
        '--color-primary': '#3498db',
        '--color-success': '#27ae60',
        '--color-warning': '#f39c12',
        '--color-danger': '#e74c3c',
        '--spacing-xs': '4px',
        '--spacing-sm': '8px',
        '--spacing-md': '16px',
        '--radius-base': '4px'
      }
      return cssVars[prop] || ''
    }
  })
})

// Mock Chart.js and canvas context
const mockCanvasContext: any = {
  fillRect: vi.fn(),
  clearRect: vi.fn(),
  getImageData: vi.fn(() => ({
    data: new Array(4).fill(0)
  })),
  putImageData: vi.fn(),
  createImageData: vi.fn(() => []),
  setTransform: vi.fn(),
  drawImage: vi.fn(),
  save: vi.fn(),
  restore: vi.fn(),
  beginPath: vi.fn(),
  moveTo: vi.fn(),
  lineTo: vi.fn(),
  closePath: vi.fn(),
  stroke: vi.fn(),
  fill: vi.fn(),
  measureText: vi.fn(() => ({ width: 0 })),
  arc: vi.fn(),
  fillText: vi.fn(),
  getContext: vi.fn(() => mockCanvasContext),
  toDataURL: vi.fn(() => '')
}

Object.defineProperty(HTMLCanvasElement.prototype, 'getContext', {
  value: vi.fn(() => mockCanvasContext)
})

// Global test configuration for Vue Test Utils
const globalConfig: typeof config = {
  global: {
    stubs: {
      // Stub Chart.js components as they require canvas
      DoughnutChart: {
        template: '<div class="mock-doughnut-chart">{{ title }}</div>',
        props: ['data', 'title']
      }
    }
  }
}

export { globalConfig as config }
