import { test, expect } from '@playwright/test'

test.describe('Authentication Flows', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the login page before each test
    await page.goto('/login')
  })

  test('should display login page correctly', async ({ page }) => {
    // Check page title
    await expect(page).toHaveTitle('Incident Management System')
    
    // Check main heading
    await expect(page.locator('h2')).toHaveText('Welcome Back')
    
    // Check form elements are present
    await expect(page.locator('input[type="text"]')).toBeVisible()
    await expect(page.locator('input[type="password"]')).toBeVisible()
    await expect(page.locator('button[type="submit"]')).toBeVisible()
    
    // Check navigation links
    await expect(page.locator('text=Forgot password?')).toBeVisible()
    await expect(page.locator('text=Sign up here')).toBeVisible()
  })

  test('should navigate to register page', async ({ page }) => {
    // Click register link
    await page.click('text=Sign up here')
    
    // Should navigate to register page
    await expect(page).toHaveURL('/register')
    await expect(page.locator('h2')).toHaveText('Create Account')
  })

  test('should display register page correctly', async ({ page }) => {
    await page.goto('/register')
    
    // Check page title and heading
    await expect(page).toHaveTitle('Incident Management System')
    await expect(page.locator('h2')).toHaveText('Create Account')
    
    // Check all form fields are present
    await expect(page.locator('#fullName')).toBeVisible()
    await expect(page.locator('#username')).toBeVisible()
    await expect(page.locator('#email')).toBeVisible()
    await expect(page.locator('#password')).toBeVisible()
    await expect(page.locator('#confirmPassword')).toBeVisible()
    await expect(page.locator('input[type="checkbox"]')).toBeVisible()
    
    // Check navigation link
    await expect(page.locator('text=Sign in here')).toBeVisible()
  })

  test('should navigate back to login from register', async ({ page }) => {
    await page.goto('/register')
    
    // Click login link
    await page.click('text=Sign in here')
    
    // Should navigate to login page
    await expect(page).toHaveURL('/login')
    await expect(page.locator('h2')).toHaveText('Welcome Back')
  })

  test('should validate login form', async ({ page }) => {
    // Submit empty form - button should be disabled
    const submitButton = page.locator('button[type="submit"]')
    await expect(submitButton).toBeDisabled()
    
    // Fill username only
    await page.fill('input[type="text"]', 'testuser')
    await expect(submitButton).toBeDisabled()
    
    // Fill password only (clear username first)
    await page.fill('input[type="text"]', '')
    await page.fill('input[type="password"]', 'password123')
    await expect(submitButton).toBeDisabled()
    
    // Fill both fields
    await page.fill('input[type="text"]', 'testuser')
    await page.fill('input[type="password"]', 'password123')
    await expect(submitButton).toBeEnabled()
  })

  test('should validate register form', async ({ page }) => {
    await page.goto('/register')
    
    const submitButton = page.locator('button[type="submit"]')
    
    // Submit empty form - button should be disabled
    await expect(submitButton).toBeDisabled()
    
    // Fill all fields except terms checkbox
    await page.fill('#fullName', 'John Doe')
    await page.fill('#username', 'johndoe')
    await page.fill('#email', 'john@example.com')
    await page.fill('#password', 'SecurePass123!')
    await page.fill('#confirmPassword', 'SecurePass123!')
    
    // Still disabled without terms
    await expect(submitButton).toBeDisabled()
    
    // Check terms checkbox
    await page.check('input[type="checkbox"]')
    
    // Now should be enabled
    await expect(submitButton).toBeEnabled()
  })

  test('should toggle password visibility', async ({ page }) => {
    const passwordInput = page.locator('input[type="password"]')
    const toggleButton = page.locator('button[aria-label="Show password"]')
    
    // Initially password type
    await expect(passwordInput).toHaveAttribute('type', 'password')
    
    // Fill password field
    await page.fill('input[type="password"]', 'testpassword')
    
    // Click toggle
    await toggleButton.click()
    
    // Should change to text type
    const revealedInput = page.locator('input[type="text"]').last()
    await expect(revealedInput).toHaveValue('testpassword')
  })

  test('should handle form submission with mock API error', async ({ page }) => {
    // Mock API to return error
    await page.route('/api/auth/login', (route) => {
      route.fulfill({
        status: 401,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'Invalid credentials' })
      })
    })
    
    // Fill and submit form
    await page.fill('input[type="text"]', 'wronguser')
    await page.fill('input[type="password"]', 'wrongpass')
    await page.click('button[type="submit"]')
    
    // Should show error message
    await expect(page.locator('[role="alert"]')).toContainText('Invalid credentials')
  })

  test('should handle successful login redirect', async ({ page }) => {
    // Mock successful API response
    await page.route('/api/auth/login', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          token: 'mock-token',
          refresh_token: 'mock-refresh',
          user: { id: '1', username: 'testuser', email: 'test@example.com' },
          expires_at: new Date(Date.now() + 3600000).toISOString()
        })
      })
    })
    
    // Mock dashboard route
    await page.route('/', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'text/html',
        body: '<html><body><h1>Dashboard</h1></body></html>'
      })
    })
    
    // Fill and submit form
    await page.fill('input[type="text"]', 'testuser')
    await page.fill('input[type="password"]', 'password123')
    await page.click('button[type="submit"]')
    
    // Should redirect to dashboard
    await expect(page).toHaveURL('/')
  })

  test('should be accessible', async ({ page }) => {
    // Check that form inputs have proper labels
    await expect(page.locator('label[for="username"]')).toBeVisible()
    await expect(page.locator('label[for="password"]')).toBeVisible()
    
    // Check ARIA attributes
    const usernameInput = page.locator('#username')
    const passwordInput = page.locator('#password')
    
    await expect(usernameInput).toHaveAttribute('aria-invalid', 'false')
    await expect(passwordInput).toHaveAttribute('aria-invalid', 'false')
    
    // Check autocomplete attributes
    await expect(usernameInput).toHaveAttribute('autocomplete', 'username')
    await expect(passwordInput).toHaveAttribute('autocomplete', 'current-password')
    
    // Test keyboard navigation
    await page.keyboard.press('Tab') // Should focus username input
    await expect(usernameInput).toBeFocused()
    
    await page.keyboard.press('Tab') // Should focus password input
    await expect(passwordInput).toBeFocused()
  })
})