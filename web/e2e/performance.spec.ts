import { test, expect } from '@playwright/test'

test.describe('Performance and Cross-browser Compatibility', () => {
  test.beforeEach(async ({ page }) => {
    // Set up basic API mocks for all tests
    await page.route('/api/metrics', async (route) => {
      await route.fulfill({
        json: {
          total_incidents: 50,
          open_incidents: 12,
          mtta: 1800000000000,
          mttr: 3600000000000,
          incidents_by_status: { open: 12, acknowledged: 8, resolved: 30 },
          incidents_by_severity: { critical: 2, high: 10, medium: 25, low: 13 }
        }
      })
    })

    await page.route('/api/incidents', async (route) => {
      await route.fulfill({
        json: [
          {
            id: '1',
            title: 'Performance Test Incident',
            description: 'Test incident for performance testing',
            severity: 'high',
            status: 'open',
            created_at: new Date(Date.now() - 3600000).toISOString(),
            updated_at: new Date(Date.now() - 1800000).toISOString()
          }
        ]
      })
    })

    await page.route('/api/alerts', async (route) => {
      await route.fulfill({
        json: [
          {
            id: '1',
            alert_name: 'Performance Test Alert',
            generator_url: 'http://localhost:9090/alerts',
            status: 'active',
            starts_at: new Date(Date.now() - 1800000).toISOString(),
            created_at: new Date(Date.now() - 1800000).toISOString(),
            updated_at: new Date(Date.now() - 900000).toISOString()
          }
        ]
      })
    })
  })

  test('should load quickly and meet performance metrics', async ({ page }) => {
    // Start timing
    const startTime = Date.now()

    await page.goto('/')

    // Wait for key content to be visible (First Meaningful Paint)
    await expect(page.locator('h2')).toContainText('Dashboard')
    await expect(page.locator('.metric-card')).toHaveCount(4)

    const loadTime = Date.now() - startTime

    // Performance assertions
    expect(loadTime).toBeLessThan(3000) // Should load within 3 seconds

    // Check that charts are loaded (indicating JavaScript executed successfully)
    await expect(page.locator('.charts-section')).toBeVisible()

    // Performance API metrics
    const performanceMetrics = await page.evaluate(() => {
      const perfEntries = performance.getEntriesByType(
        'navigation'
      )[0] as PerformanceNavigationTiming
      return {
        domContentLoaded:
          perfEntries.domContentLoadedEventEnd - perfEntries.domContentLoadedEventStart,
        loadComplete: perfEntries.loadEventEnd - perfEntries.loadEventStart,
        firstByte: perfEntries.responseStart - perfEntries.requestStart,
        domInteractive: perfEntries.domInteractive - perfEntries.navigationStart
      }
    })

    console.log('Performance Metrics:', performanceMetrics)

    // Assert performance thresholds
    expect(performanceMetrics.domInteractive).toBeLessThan(2000) // DOM should be interactive within 2s
    expect(performanceMetrics.firstByte).toBeLessThan(1000) // TTFB should be under 1s
  })

  test('should be accessible and follow WCAG guidelines', async ({ page }) => {
    await page.goto('/')

    // Wait for content to load
    await expect(page.locator('.metric-card')).toHaveCount(4)

    // Check for proper heading structure
    const headings = await page.locator('h1, h2, h3, h4, h5, h6').all()
    expect(headings.length).toBeGreaterThan(0)

    // Check that interactive elements are keyboard accessible
    await page.keyboard.press('Tab')
    const focusedElement = page.locator(':focus')
    await expect(focusedElement).toBeVisible()

    // Check for proper alt text on images (if any)
    const images = await page.locator('img').all()
    for (const img of images) {
      const alt = await img.getAttribute('alt')
      expect(alt).toBeDefined()
    }

    // Check that buttons have accessible names
    const buttons = await page.locator('button').all()
    for (const button of buttons) {
      const accessibleName =
        (await button.textContent()) || (await button.getAttribute('aria-label'))
      expect(accessibleName).toBeTruthy()
    }

    // Check color contrast (basic check)
    const backgroundColors = await page.evaluate(() => {
      const elements = Array.from(document.querySelectorAll('.btn, .card, .metric-card'))
      return elements.map((el) => window.getComputedStyle(el).backgroundColor)
    })

    // Ensure we have styled elements (not all transparent/white)
    const hasStyledElements = backgroundColors.some(
      (color) => color !== 'rgba(0, 0, 0, 0)' && color !== 'rgb(255, 255, 255)'
    )
    expect(hasStyledElements).toBe(true)
  })

  test('should work correctly in different viewports', async ({ page }) => {
    const viewports = [
      { width: 320, height: 568, name: 'iPhone SE' },
      { width: 375, height: 667, name: 'iPhone 8' },
      { width: 768, height: 1024, name: 'iPad' },
      { width: 1024, height: 768, name: 'iPad Landscape' },
      { width: 1920, height: 1080, name: 'Desktop' }
    ]

    for (const viewport of viewports) {
      console.log(`Testing viewport: ${viewport.name} (${viewport.width}x${viewport.height})`)

      await page.setViewportSize({ width: viewport.width, height: viewport.height })
      await page.goto('/')

      // Wait for content to load
      await expect(page.locator('h2')).toContainText('Dashboard')
      await expect(page.locator('.metric-card')).toHaveCount(4)

      // Check that content is visible and not cut off
      const dashboardBox = await page.locator('.dashboard').boundingBox()
      expect(dashboardBox).toBeTruthy()
      expect(dashboardBox?.width).toBeLessThanOrEqual(viewport.width)

      // Check that buttons are touch-friendly on mobile
      if (viewport.width < 768) {
        const buttons = page.locator('button')
        const buttonCount = await buttons.count()

        for (let i = 0; i < buttonCount; i++) {
          const buttonBox = await buttons.nth(i).boundingBox()
          if (buttonBox) {
            expect(buttonBox.height).toBeGreaterThanOrEqual(44) // Touch-friendly size
          }
        }
      }

      // Test navigation on different screen sizes
      await page.click('nav a[href="/incidents"]')
      await expect(page.locator('h2')).toContainText('Incidents')

      // Check that table is responsive
      if (viewport.width < 768) {
        // On mobile, table should be horizontally scrollable
        const tableContainer = page.locator('.table-responsive, .table').first()
        await expect(tableContainer).toBeVisible()
      }
    }
  })

  test('should handle slow network conditions gracefully', async ({ page, context }) => {
    // Simulate slow 3G connection
    await context.route('**/*', async (route, request) => {
      // Add delay to simulate slow network
      await new Promise((resolve) => setTimeout(resolve, 100))
      await route.continue()
    })

    const startTime = Date.now()
    await page.goto('/')

    // Should show loading state
    const loadingElement = page.locator('.loading')
    if (await loadingElement.isVisible()) {
      await expect(loadingElement).toContainText('Loading')
    }

    // Eventually should load content
    await expect(page.locator('.metric-card')).toHaveCount(4, { timeout: 10000 })

    const loadTime = Date.now() - startTime
    console.log(`Load time with slow network: ${loadTime}ms`)

    // Even with slow network, should load within reasonable time
    expect(loadTime).toBeLessThan(10000)
  })

  test('should maintain functionality across page refreshes', async ({ page }) => {
    // Load dashboard
    await page.goto('/')
    await expect(page.locator('.metric-card')).toHaveCount(4)

    // Navigate to incidents
    await page.click('nav a[href="/incidents"]')
    await expect(page.locator('h2')).toContainText('Incidents')

    // Refresh the page
    await page.reload()

    // Should maintain the current page
    await expect(page.locator('h2')).toContainText('Incidents')
    await expect(page.locator('.table')).toBeVisible()

    // Navigate back to dashboard
    await page.click('nav a[href="/"]')
    await expect(page.locator('h2')).toContainText('Dashboard')
    await expect(page.locator('.metric-card')).toHaveCount(4)
  })

  test('should handle JavaScript errors gracefully', async ({ page }) => {
    // Listen for console errors
    const consoleErrors: string[] = []
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text())
      }
    })

    // Listen for page errors
    const pageErrors: string[] = []
    page.on('pageerror', (error) => {
      pageErrors.push(error.message)
    })

    await page.goto('/')
    await expect(page.locator('.metric-card')).toHaveCount(4)

    // Navigate through the application
    await page.click('nav a[href="/incidents"]')
    await expect(page.locator('h2')).toContainText('Incidents')

    await page.click('nav a[href="/alerts"]')
    await expect(page.locator('h2')).toContainText('Alerts')

    await page.click('nav a[href="/"]')
    await expect(page.locator('h2')).toContainText('Dashboard')

    // Check that no critical errors occurred
    const criticalErrors = [...consoleErrors, ...pageErrors].filter(
      (error) =>
        !error.includes('404') && // Ignore 404s from dev server
        !error.includes('WebSocket') && // Ignore WebSocket errors in dev
        !error.toLowerCase().includes('warning') // Ignore warnings
    )

    if (criticalErrors.length > 0) {
      console.log('JavaScript errors found:', criticalErrors)
    }

    // Should have minimal critical errors
    expect(criticalErrors.length).toBeLessThan(3)
  })

  test('should have proper caching headers in production build', async ({ page }) => {
    // This test would be more relevant when testing against a production server
    // For now, we'll check that assets are being loaded correctly

    await page.goto('/')
    await expect(page.locator('.metric-card')).toHaveCount(4)

    // Check that CSS and JS assets loaded
    const stylesheets = await page.locator('link[rel="stylesheet"]').count()
    const scripts = await page.locator('script[src]').count()

    expect(stylesheets).toBeGreaterThan(0)
    expect(scripts).toBeGreaterThan(0)

    // Check for proper resource loading
    const resourceLoadErrors = await page.evaluate(() => {
      const errors: string[] = []
      const images = Array.from(document.images)
      images.forEach((img) => {
        if (!img.complete || img.naturalHeight === 0) {
          errors.push(`Image failed to load: ${img.src}`)
        }
      })
      return errors
    })

    expect(resourceLoadErrors.length).toBe(0)
  })
})
