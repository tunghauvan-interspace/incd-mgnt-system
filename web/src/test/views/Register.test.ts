import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import Register from '@/views/Register.vue'
import { useRouter } from 'vue-router'

// Mock vue-router
vi.mock('vue-router', () => ({
  useRouter: vi.fn()
}))

// Mock the form validation composable
vi.mock('@/composables/useFormValidation', () => ({
  useFormValidation: () => ({
    errors: { value: {} },
    isValid: { value: true },
    validate: vi.fn(() => true),
    clearError: vi.fn(),
    clearAllErrors: vi.fn()
  })
}))

// Mock auth store
vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({
    user: null,
    token: null,
    isLoading: false,
    error: null,
    login: vi.fn(),
    register: vi.fn()
  })
}))

describe('Register.vue', () => {
  let wrapper: any
  let mockRouter: any
  
  beforeEach(() => {
    mockRouter = {
      push: vi.fn()
    }
    vi.mocked(useRouter).mockReturnValue(mockRouter)
    setActivePinia(createPinia())
  })

  const createWrapper = () => {
    return mount(Register, {
      global: {
        stubs: {
          'router-link': true
        }
      }
    })
  }

  it('renders registration form correctly', () => {
    wrapper = createWrapper()
    
    expect(wrapper.find('h2').text()).toBe('Create Account')
    expect(wrapper.find('#fullName').attributes('placeholder')).toBe('Enter your full name')
    expect(wrapper.find('#username').attributes('placeholder')).toBe('Choose a username')
    expect(wrapper.find('#email').attributes('placeholder')).toBe('Enter your email address')
    expect(wrapper.find('#password').attributes('placeholder')).toBe('Create a strong password')
    expect(wrapper.find('#confirmPassword').attributes('placeholder')).toBe('Confirm your password')
    expect(wrapper.find('button[type="submit"]').text()).toBe('Create Account')
  })

  it('has all required form fields', () => {
    wrapper = createWrapper()
    
    expect(wrapper.find('#fullName').exists()).toBe(true)
    expect(wrapper.find('#username').exists()).toBe(true)
    expect(wrapper.find('#email').exists()).toBe(true)
    expect(wrapper.find('#password').exists()).toBe(true)
    expect(wrapper.find('#confirmPassword').exists()).toBe(true)
    expect(wrapper.find('input[type="checkbox"]').exists()).toBe(true)
  })

  it('toggles password visibility', async () => {
    wrapper = createWrapper()
    
    const passwordInput = wrapper.find('#password')
    const toggleButton = wrapper.find('button[aria-label="Show password"]')
    
    expect(passwordInput.attributes('type')).toBe('password')
    
    await toggleButton.trigger('click')
    await wrapper.vm.$nextTick()
    
    expect(passwordInput.attributes('type')).toBe('text')
  })

  it('toggles confirm password visibility', async () => {
    wrapper = createWrapper()
    
    const confirmPasswordInput = wrapper.find('#confirmPassword')
    const toggleButton = wrapper.find('button[aria-label="Show password confirmation"]')
    
    expect(confirmPasswordInput.attributes('type')).toBe('password')
    
    await toggleButton.trigger('click')
    await wrapper.vm.$nextTick()
    
    expect(confirmPasswordInput.attributes('type')).toBe('text')
  })

  it('shows password strength requirements', () => {
    wrapper = createWrapper()
    
    const passwordHelp = wrapper.find('#password-help')
    expect(passwordHelp.text()).toContain('Must be at least 8 characters')
    expect(passwordHelp.text()).toContain('mix of letters, numbers, and symbols')
  })

  it('validates password strength', () => {
    wrapper = createWrapper()
    const component = wrapper.vm
    
    // Test weak password
    expect(component.validatePasswordStrength('weak')).toContain('at least 8 characters')
    
    // Test password without variety
    expect(component.validatePasswordStrength('password')).toContain('mix of uppercase, lowercase')
    
    // Test strong password
    expect(component.validatePasswordStrength('SecurePass123!')).toBe(null)
  })

  it('validates password confirmation', () => {
    wrapper = createWrapper()
    const component = wrapper.vm
    
    // Set the original password
    component.registerForm.password = 'SecurePass123!'
    
    // Test non-matching password
    expect(component.validatePasswordConfirmation('different')).toContain('do not match')
    
    // Test matching password
    expect(component.validatePasswordConfirmation('SecurePass123!')).toBe(null)
  })

  it('has proper accessibility attributes', () => {
    wrapper = createWrapper()
    
    const fullNameInput = wrapper.find('#fullName')
    const usernameInput = wrapper.find('#username')
    const emailInput = wrapper.find('#email')
    const passwordInput = wrapper.find('#password')
    
    expect(fullNameInput.attributes('aria-invalid')).toBe('false')
    expect(usernameInput.attributes('aria-invalid')).toBe('false')
    expect(emailInput.attributes('aria-invalid')).toBe('false')
    expect(passwordInput.attributes('aria-invalid')).toBe('false')
    
    expect(fullNameInput.attributes('autocomplete')).toBe('name')
    expect(usernameInput.attributes('autocomplete')).toBe('username')
    expect(emailInput.attributes('autocomplete')).toBe('email')
    expect(passwordInput.attributes('autocomplete')).toBe('new-password')
  })

  it('requires terms acceptance for form submission', () => {
    wrapper = createWrapper()
    
    const checkbox = wrapper.find('input[type="checkbox"]')
    expect(checkbox.attributes('required')).toBeDefined()
    expect(wrapper.text()).toContain('Terms of Service')
    expect(wrapper.text()).toContain('Privacy Policy')
  })

  it('has navigation links', () => {
    wrapper = createWrapper()
    
    expect(wrapper.text()).toContain('Already have an account?')
    // Router links are stubbed, so we check for their presence differently
    const routerLinks = wrapper.findAll('[to="/login"]')
    expect(routerLinks.length).toBeGreaterThan(0)
  })

  it('has all required fields marked', () => {
    wrapper = createWrapper()
    
    expect(wrapper.find('label[for="fullName"]').text()).toContain('*')
    expect(wrapper.find('label[for="username"]').text()).toContain('*')
    expect(wrapper.find('label[for="email"]').text()).toContain('*')
    expect(wrapper.find('label[for="password"]').text()).toContain('*')
    expect(wrapper.find('label[for="confirmPassword"]').text()).toContain('*')
  })
})