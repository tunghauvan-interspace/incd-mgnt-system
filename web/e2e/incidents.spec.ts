import { test, expect } from '@playwright/test'

test.describe('Incidents Management', () => {
  test.beforeEach(async ({ page }) => {
    // Mock the incidents API
    await page.route('/api/incidents', async (route) => {
      await route.fulfill({
        json: [
          {
            id: '1',
            title: 'Database Connection Timeout',
            description: 'Connection to primary database is timing out after 30 seconds',
            severity: 'critical',
            status: 'open',
            created_at: '2023-12-25T10:00:00Z',
            updated_at: '2023-12-25T10:00:00Z',
            labels: { service: 'database', environment: 'production' }
          },
          {
            id: '2',
            title: 'High Memory Usage',
            description: 'Memory usage is consistently above 90% on web servers',
            severity: 'high',
            status: 'acknowledged',
            assignee_id: 'admin',
            acknowledged_at: '2023-12-25T11:00:00Z',
            created_at: '2023-12-25T09:30:00Z',
            updated_at: '2023-12-25T11:00:00Z',
            labels: { service: 'webserver', environment: 'production' }
          },
          {
            id: '3',
            title: 'SSL Certificate Expiring',
            description: 'SSL certificate for api.example.com expires in 7 days',
            severity: 'medium',
            status: 'resolved',
            assignee_id: 'security-team',
            acknowledged_at: '2023-12-24T14:00:00Z',
            resolved_at: '2023-12-24T16:00:00Z',
            created_at: '2023-12-24T13:00:00Z',
            updated_at: '2023-12-24T16:00:00Z',
            labels: { service: 'api', environment: 'production' }
          }
        ]
      })
    })
  })

  test('should display incidents list', async ({ page }) => {
    await page.goto('/incidents')

    // Check page title and header
    await expect(page.locator('h2')).toContainText('Incidents')

    // Wait for incidents to load
    await expect(page.locator('.table')).toBeVisible()
    await expect(page.locator('tbody tr')).toHaveCount(3)

    // Check incident data
    const firstRow = page.locator('tbody tr').first()
    await expect(firstRow).toContainText('Database Connection Timeout')
    await expect(firstRow).toContainText('critical')
    await expect(firstRow).toContainText('open')

    // Check status badges
    await expect(page.locator('.status-open')).toBeVisible()
    await expect(page.locator('.status-acknowledged')).toBeVisible()
    await expect(page.locator('.status-resolved')).toBeVisible()

    // Check severity badges
    await expect(page.locator('.severity-critical')).toBeVisible()
    await expect(page.locator('.severity-high')).toBeVisible()
    await expect(page.locator('.severity-medium')).toBeVisible()
  })

  test('should open incident details modal', async ({ page }) => {
    await page.goto('/incidents')

    // Wait for table to load
    await expect(page.locator('.table')).toBeVisible()

    // Click on the first incident to open details
    await page.locator('tbody tr').first().click()

    // Check that modal opens
    await expect(page.locator('.modal-overlay')).toBeVisible()
    await expect(page.locator('.modal h3')).toContainText('Database Connection Timeout')

    // Check incident details in modal
    await expect(page.locator('.modal')).toContainText(
      'Connection to primary database is timing out'
    )
    await expect(page.locator('.modal')).toContainText('critical')
    await expect(page.locator('.modal')).toContainText('open')
  })

  test('should acknowledge an incident', async ({ page }) => {
    // Mock the acknowledge API call
    await page.route('/api/incidents/1/acknowledge', async (route) => {
      await route.fulfill({
        status: 200,
        body: ''
      })
    })

    // Mock updated incidents list after acknowledgment
    await page.route('/api/incidents', async (route) => {
      await route.fulfill({
        json: [
          {
            id: '1',
            title: 'Database Connection Timeout',
            description: 'Connection to primary database is timing out after 30 seconds',
            severity: 'critical',
            status: 'acknowledged',
            assignee_id: 'current-user',
            acknowledged_at: '2023-12-25T12:00:00Z',
            created_at: '2023-12-25T10:00:00Z',
            updated_at: '2023-12-25T12:00:00Z'
          }
        ]
      })
    })

    await page.goto('/incidents')

    // Click on first incident to open modal
    await page.locator('tbody tr').first().click()
    await expect(page.locator('.modal-overlay')).toBeVisible()

    // Click acknowledge button
    await page.click('.btn-warning')

    // Fill in assignee (if required)
    const assigneeInput = page.locator('input[placeholder="Enter assignee ID"]')
    if (await assigneeInput.isVisible()) {
      await assigneeInput.fill('current-user')
      await page.click('button[type="submit"]')
    }

    // Check that modal closes
    await expect(page.locator('.modal-overlay')).not.toBeVisible()

    // Verify the incident status updated (refresh may be needed)
    await page.reload()
    await expect(page.locator('.status-acknowledged')).toBeVisible()
  })

  test('should resolve an incident', async ({ page }) => {
    // Start with an acknowledged incident
    await page.route('/api/incidents', async (route) => {
      await route.fulfill({
        json: [
          {
            id: '2',
            title: 'High Memory Usage',
            description: 'Memory usage is consistently above 90% on web servers',
            severity: 'high',
            status: 'acknowledged',
            assignee_id: 'admin',
            acknowledged_at: '2023-12-25T11:00:00Z',
            created_at: '2023-12-25T09:30:00Z',
            updated_at: '2023-12-25T11:00:00Z'
          }
        ]
      })
    })

    // Mock the resolve API call
    await page.route('/api/incidents/2/resolve', async (route) => {
      await route.fulfill({
        status: 200,
        body: ''
      })
    })

    await page.goto('/incidents')

    // Click on the acknowledged incident
    await page.locator('tbody tr').first().click()
    await expect(page.locator('.modal-overlay')).toBeVisible()

    // Click resolve button
    await page.click('.btn-success')

    // Check that modal closes
    await expect(page.locator('.modal-overlay')).not.toBeVisible()
  })

  test('should refresh incidents list', async ({ page }) => {
    let apiCallCount = 0

    await page.route('/api/incidents', async (route) => {
      apiCallCount++
      await route.fulfill({
        json: [
          {
            id: '1',
            title: `Incident ${apiCallCount}`,
            description: 'Test incident description',
            severity: 'medium',
            status: 'open',
            created_at: '2023-12-25T10:00:00Z',
            updated_at: '2023-12-25T10:00:00Z'
          }
        ]
      })
    })

    await page.goto('/incidents')

    // Check initial content
    await expect(page.locator('tbody tr')).toContainText('Incident 1')

    // Click refresh button
    await page.click('.refresh-btn button')

    // Check that content updated
    await expect(page.locator('tbody tr')).toContainText('Incident 2')

    expect(apiCallCount).toBe(2)
  })

  test('should close modal with escape key', async ({ page }) => {
    await page.goto('/incidents')

    // Open modal by clicking incident
    await page.locator('tbody tr').first().click()
    await expect(page.locator('.modal-overlay')).toBeVisible()

    // Press Escape key
    await page.keyboard.press('Escape')

    // Check that modal closes
    await expect(page.locator('.modal-overlay')).not.toBeVisible()
  })

  test('should close modal by clicking overlay', async ({ page }) => {
    await page.goto('/incidents')

    // Open modal by clicking incident
    await page.locator('tbody tr').first().click()
    await expect(page.locator('.modal-overlay')).toBeVisible()

    // Click on the overlay (outside modal content)
    await page.locator('.modal-overlay').click({
      position: { x: 10, y: 10 } // Click near top-left of overlay
    })

    // Check that modal closes
    await expect(page.locator('.modal-overlay')).not.toBeVisible()
  })

  test('should handle empty incidents list', async ({ page }) => {
    await page.route('/api/incidents', async (route) => {
      await route.fulfill({
        json: []
      })
    })

    await page.goto('/incidents')

    // Should show empty state
    await expect(page.locator('.no-data')).toContainText('No incidents found')
  })

  test('should handle API errors', async ({ page }) => {
    await page.route('/api/incidents', async (route) => {
      await route.fulfill({
        status: 500,
        body: 'Internal Server Error'
      })
    })

    await page.goto('/incidents')

    // Should show error message
    await expect(page.locator('.error-message')).toContainText('Error loading incidents')
  })
})
