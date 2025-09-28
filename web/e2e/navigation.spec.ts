import { test, expect } from '@playwright/test';

test.describe('Navigation', () => {
  test('should navigate between pages', async ({ page }) => {
    // Mock API responses
    await page.route('/api/metrics', async route => {
      await route.fulfill({
        json: {
          total_incidents: 10,
          open_incidents: 3,
          mtta: 1800000000000,
          mttr: 3600000000000,
          incidents_by_status: { open: 3, acknowledged: 2, resolved: 5 },
          incidents_by_severity: { critical: 1, high: 2, medium: 4, low: 3 }
        }
      });
    });

    await page.route('/api/incidents', async route => {
      await route.fulfill({
        json: [
          {
            id: '1',
            title: 'Test Incident 1',
            description: 'Test Description 1',
            severity: 'high',
            status: 'open',
            created_at: '2023-12-25T10:00:00Z',
            updated_at: '2023-12-25T10:00:00Z'
          },
          {
            id: '2',
            title: 'Test Incident 2',
            description: 'Test Description 2',
            severity: 'medium',
            status: 'acknowledged',
            assignee_id: 'user1',
            acknowledged_at: '2023-12-25T11:00:00Z',
            created_at: '2023-12-25T09:30:00Z',
            updated_at: '2023-12-25T11:00:00Z'
          }
        ]
      });
    });

    await page.route('/api/alerts', async route => {
      await route.fulfill({
        json: [
          {
            id: '1',
            alert_name: 'High CPU Usage',
            generator_url: 'http://localhost:9090/alerts',
            status: 'active',
            starts_at: '2023-12-25T10:30:00Z',
            labels: { severity: 'warning', instance: 'server1' },
            annotations: { summary: 'CPU usage is above 80%' },
            created_at: '2023-12-25T10:30:00Z',
            updated_at: '2023-12-25T10:30:00Z'
          }
        ]
      });
    });

    await page.goto('/');

    // Start at Dashboard
    await expect(page.locator('h2')).toContainText('Dashboard');

    // Navigate to Incidents
    await page.click('nav a[href="/incidents"]');
    await expect(page.locator('h2')).toContainText('Incidents');
    await expect(page.locator('.table')).toBeVisible();

    // Navigate to Alerts
    await page.click('nav a[href="/alerts"]');
    await expect(page.locator('h2')).toContainText('Alerts');
    await expect(page.locator('.table')).toBeVisible();

    // Navigate back to Dashboard
    await page.click('nav a[href="/"]');
    await expect(page.locator('h2')).toContainText('Dashboard');
  });

  test('should maintain navigation state on refresh', async ({ page }) => {
    await page.route('/api/incidents', async route => {
      await route.fulfill({
        json: [
          {
            id: '1',
            title: 'Persistent Incident',
            description: 'Test Description',
            severity: 'high',
            status: 'open',
            created_at: '2023-12-25T10:00:00Z',
            updated_at: '2023-12-25T10:00:00Z'
          }
        ]
      });
    });

    // Navigate to incidents page
    await page.goto('/incidents');
    await expect(page.locator('h2')).toContainText('Incidents');

    // Refresh the page
    await page.reload();

    // Should still be on incidents page
    await expect(page.locator('h2')).toContainText('Incidents');
    await expect(page.locator('.table')).toBeVisible();
  });

  test('should handle direct URL navigation', async ({ page }) => {
    await page.route('/api/alerts', async route => {
      await route.fulfill({
        json: [
          {
            id: '1',
            alert_name: 'Direct Navigation Alert',
            generator_url: 'http://localhost:9090/alerts',
            status: 'active',
            starts_at: '2023-12-25T10:30:00Z',
            created_at: '2023-12-25T10:30:00Z',
            updated_at: '2023-12-25T10:30:00Z'
          }
        ]
      });
    });

    // Navigate directly to alerts page
    await page.goto('/alerts');
    
    // Should load the alerts page correctly
    await expect(page.locator('h2')).toContainText('Alerts');
    await expect(page.locator('.table')).toBeVisible();
  });
});