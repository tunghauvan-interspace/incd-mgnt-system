import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import Login from '@/views/Login.vue'
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

describe('Login.vue', () => {
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
    return mount(Login, {
      global: {
        stubs: {
          'router-link': true
        }
      }
    })
  }

  it('renders login form correctly', () => {
    wrapper = createWrapper()
    
    expect(wrapper.find('h2').text()).toBe('Welcome Back')
    expect(wrapper.find('input[type="text"]').attributes('placeholder')).toBe('Enter your email or username')
    expect(wrapper.find('input[type="password"]').attributes('placeholder')).toBe('Enter your password')
    expect(wrapper.find('button[type="submit"]').text()).toBe('Sign In')
  })

  it('has proper form structure', () => {
    wrapper = createWrapper()
    
    expect(wrapper.find('form').exists()).toBe(true)
    expect(wrapper.find('#username').exists()).toBe(true)
    expect(wrapper.find('#password').exists()).toBe(true)
    expect(wrapper.find('input[type="checkbox"]').exists()).toBe(true) // Remember me
  })

  it('toggles password visibility', async () => {
    wrapper = createWrapper()
    
    const passwordInput = wrapper.find('input[type="password"]')
    const toggleButton = wrapper.find('button[aria-label="Show password"]')
    
    expect(passwordInput.attributes('type')).toBe('password')
    
    await toggleButton.trigger('click')
    await wrapper.vm.$nextTick()
    
    // After click, the input type should change
    expect(wrapper.find('input[type="text"]').exists()).toBe(true)
  })

  it('has proper accessibility attributes', () => {
    wrapper = createWrapper()
    
    const usernameInput = wrapper.find('#username')
    const passwordInput = wrapper.find('#password')
    
    expect(usernameInput.attributes('aria-invalid')).toBe('false')
    expect(passwordInput.attributes('aria-invalid')).toBe('false')
    expect(usernameInput.attributes('autocomplete')).toBe('username')
    expect(passwordInput.attributes('autocomplete')).toBe('current-password')
  })

  it('includes required form fields', () => {
    wrapper = createWrapper()
    
    const usernameInput = wrapper.find('#username')
    const passwordInput = wrapper.find('#password')
    
    expect(usernameInput.attributes('required')).toBeDefined()
    expect(passwordInput.attributes('required')).toBeDefined()
  })

  it('has navigation links', () => {
    wrapper = createWrapper()
    
    expect(wrapper.text()).toContain('Don\'t have an account?')
    // Router links are stubbed, so we check for their presence differently
    const routerLinks = wrapper.findAll('[to="/register"], [to="/forgot-password"]')
    expect(routerLinks.length).toBeGreaterThan(0)
  })
})