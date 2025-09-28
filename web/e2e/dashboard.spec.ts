import { test, expect } from '@playwright/test'

test.describe('Dashboard Page', () => {
  test('should display dashboard with metrics', async ({ page }) => {
    // Mock the API response for metrics
    await page.route('/api/metrics', async (route) => {
      await route.fulfill({
        json: {
          total_incidents: 100,
          open_incidents: 15,
          mtta: 3600000000000, // 1 hour in nanoseconds
          mttr: 7200000000000, // 2 hours in nanoseconds
          incidents_by_status: { open: 15, acknowledged: 5, resolved: 80 },
          incidents_by_severity: { critical: 3, high: 12, medium: 50, low: 35 }
        }
      })
    })

    await page.goto('/')

    // Check page title
    await expect(page).toHaveTitle(/Incident Management/)

    // Check header navigation
    await expect(page.locator('h2')).toContainText('Dashboard')

    // Wait for metrics to load
    await expect(page.locator('.metric-card')).toHaveCount(4)

    // Check metric values
    await expect(page.locator('.metric-card').nth(0)).toContainText('Total Incidents')
    await expect(page.locator('.metric-card').nth(0)).toContainText('100')

    await expect(page.locator('.metric-card').nth(1)).toContainText('Open Incidents')
    await expect(page.locator('.metric-card').nth(1)).toContainText('15')

    await expect(page.locator('.metric-card').nth(2)).toContainText('MTTA')
    await expect(page.locator('.metric-card').nth(3)).toContainText('MTTR')

    // Check charts are present
    await expect(page.locator('.charts-section')).toBeVisible()
  })

  test('should handle refresh functionality', async ({ page }) => {
    let apiCallCount = 0

    await page.route('/api/metrics', async (route) => {
      apiCallCount++
      await route.fulfill({
        json: {
          total_incidents: 100 + apiCallCount * 10,
          open_incidents: 15,
          mtta: 3600000000000,
          mttr: 7200000000000,
          incidents_by_status: { open: 15, acknowledged: 5, resolved: 80 },
          incidents_by_severity: { critical: 3, high: 12, medium: 50, low: 35 }
        }
      })
    })

    await page.goto('/')

    // Wait for initial load
    await expect(page.locator('.metric-card').nth(0)).toContainText('110')

    // Click refresh button
    await page.click('.refresh-btn button')

    // Check that metrics updated
    await expect(page.locator('.metric-card').nth(0)).toContainText('120')

    expect(apiCallCount).toBe(2)
  })

  test('should show loading state', async ({ page }) => {
    // Create a promise that we can resolve manually
    let resolveMetrics: (value: any) => void
    const metricsPromise = new Promise((resolve) => {
      resolveMetrics = resolve
    })

    await page.route('/api/metrics', async (route) => {
      await metricsPromise
      await route.fulfill({
        json: {
          total_incidents: 50,
          open_incidents: 10,
          mtta: 1800000000000,
          mttr: 3600000000000,
          incidents_by_status: { open: 10, acknowledged: 5, resolved: 35 },
          incidents_by_severity: { critical: 2, high: 8, medium: 25, low: 15 }
        }
      })
    })

    await page.goto('/')

    // Check loading state
    await expect(page.locator('.loading')).toContainText('Loading dashboard...')
    await expect(page.locator('button')).toContainText('Loading...')
    await expect(page.locator('button')).toBeDisabled()

    // Resolve the metrics promise
    resolveMetrics(true)

    // Wait for loading to complete
    await expect(page.locator('.metric-card').nth(0)).toContainText('50')
    await expect(page.locator('.loading')).not.toBeVisible()
  })

  test('should handle API errors gracefully', async ({ page }) => {
    await page.route('/api/metrics', async (route) => {
      await route.fulfill({
        status: 500,
        body: 'Internal Server Error'
      })
    })

    await page.goto('/')

    // Check error message
    await expect(page.locator('.error-message')).toContainText('Error loading dashboard data')
    await expect(page.locator('.metrics-grid')).not.toBeVisible()
  })
})

test.describe('Mobile Responsiveness', () => {
  test('should be mobile responsive', async ({ page }) => {
    await page.route('/api/metrics', async (route) => {
      await route.fulfill({
        json: {
          total_incidents: 25,
          open_incidents: 5,
          mtta: 1800000000000,
          mttr: 3600000000000,
          incidents_by_status: { open: 5, acknowledged: 2, resolved: 18 },
          incidents_by_severity: { critical: 1, high: 4, medium: 12, low: 8 }
        }
      })
    })

    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 })
    await page.goto('/')

    // Check that the dashboard is responsive
    await expect(page.locator('.dashboard')).toBeVisible()
    await expect(page.locator('.metric-card')).toHaveCount(4)

    // Check that metric cards stack properly on mobile
    const metricCards = page.locator('.metric-card')
    const firstCard = metricCards.nth(0)
    const secondCard = metricCards.nth(1)

    const firstCardBox = await firstCard.boundingBox()
    const secondCardBox = await secondCard.boundingBox()

    // On mobile, cards should stack vertically (second card should be below first)
    expect(secondCardBox?.y).toBeGreaterThan(firstCardBox?.y || 0)
  })
})
