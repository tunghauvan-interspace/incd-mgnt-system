import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import type { LoginCredentials, RegisterData } from '@/types/auth'

export const useAuth = () => {
  const authStore = useAuthStore()
  
  // Local state for forms
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed properties
  const user = computed(() => authStore.user)
  const isAuthenticated = computed(() => authStore.isAuthenticated)

  // Login function
  const login = async (credentials: LoginCredentials) => {
    isLoading.value = true
    error.value = null
    
    try {
      const result = await authStore.login(credentials)
      if (!result.success) {
        error.value = result.error || 'Login failed'
        return false
      }
      return true
    } catch (err) {
      error.value = 'An unexpected error occurred'
      return false
    } finally {
      isLoading.value = false
    }
  }

  // Register function
  const register = async (data: RegisterData) => {
    isLoading.value = true
    error.value = null
    
    try {
      const result = await authStore.register(data)
      if (!result.success) {
        error.value = result.error || 'Registration failed'
        return false
      }
      return true
    } catch (err) {
      error.value = 'An unexpected error occurred'
      return false
    } finally {
      isLoading.value = false
    }
  }

  // Logout function
  const logout = () => {
    authStore.logout()
  }

  // Clear error
  const clearError = () => {
    error.value = null
  }

  return {
    // State
    isLoading,
    error,
    // Computed
    user,
    isAuthenticated,
    // Actions
    login,
    register,
    logout,
    clearError
  }
}