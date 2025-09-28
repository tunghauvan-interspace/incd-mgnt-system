import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, AuthResponse, LoginRequest, RegisterRequest } from '@/types/api'
import { authAPI } from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userRoles = computed(() => user.value?.roles || [])
  const userPermissions = computed(() => {
    const permissions = new Set<string>()
    userRoles.value.forEach((role) => {
      role.permissions.forEach((permission) => {
        permissions.add(permission.name)
      })
    })
    return Array.from(permissions)
  })

  const hasRole = (roleName: string) => {
    return userRoles.value.some((role) => role.name === roleName)
  }

  const hasPermission = (permissionName: string) => {
    return userPermissions.value.includes(permissionName)
  }

  const hasAnyRole = (roleNames: string[]) => {
    return roleNames.some((roleName) => hasRole(roleName))
  }

  const hasAnyPermission = (permissionNames: string[]) => {
    return permissionNames.some((permissionName) => hasPermission(permissionName))
  }

  const login = async (credentials: LoginRequest) => {
    isLoading.value = true
    error.value = null

    try {
      const response: AuthResponse = await authAPI.login(credentials)
      setAuthData(response)
      return response
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Login failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const register = async (userData: RegisterRequest) => {
    isLoading.value = true
    error.value = null

    try {
      const response: AuthResponse = await authAPI.register(userData)
      setAuthData(response)
      return response
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Registration failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const refreshAuthToken = async () => {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    try {
      const response: AuthResponse = await authAPI.refreshToken(refreshToken.value)
      setAuthData(response)
      return response
    } catch (err: any) {
      logout()
      throw err
    }
  }

  const logout = () => {
    user.value = null
    token.value = null
    refreshToken.value = null
    error.value = null

    // Clear localStorage
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_refresh_token')
    localStorage.removeItem('auth_user')
  }

  const setAuthData = (response: AuthResponse) => {
    user.value = response.user
    token.value = response.token
    refreshToken.value = response.refresh_token

    // Store in localStorage
    localStorage.setItem('auth_token', response.token)
    localStorage.setItem('auth_refresh_token', response.refresh_token)
    localStorage.setItem('auth_user', JSON.stringify(response.user))
  }

  const loadAuthFromStorage = () => {
    const storedToken = localStorage.getItem('auth_token')
    const storedRefreshToken = localStorage.getItem('auth_refresh_token')
    const storedUser = localStorage.getItem('auth_user')

    if (storedToken && storedUser) {
      token.value = storedToken
      refreshToken.value = storedRefreshToken
      try {
        user.value = JSON.parse(storedUser)
      } catch (err) {
        console.error('Failed to parse stored user data:', err)
        logout()
      }
    }
  }

  // Initialize auth from storage on store creation
  loadAuthFromStorage()

  return {
    // State
    user,
    token,
    refreshToken,
    isLoading,
    error,

    // Getters
    isAuthenticated,
    userRoles,
    userPermissions,

    // Actions
    hasRole,
    hasPermission,
    hasAnyRole,
    hasAnyPermission,
    login,
    register,
    refreshAuthToken,
    logout,
    loadAuthFromStorage
  }
})
