<template>
  <div class="min-h-screen bg-gradient-to-br from-emerald-50 via-teal-50 to-cyan-50 flex items-center justify-center px-4 sm:px-6 lg:px-8 relative overflow-hidden">
    <!-- Enhanced background with animated particles -->
    <div class="absolute inset-0 overflow-hidden">
      <!-- Main gradient blobs -->
      <div class="absolute -top-40 -right-40 w-80 h-80 rounded-full bg-gradient-to-br from-emerald-400/20 to-teal-400/20 blur-3xl animate-pulse"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 rounded-full bg-gradient-to-tr from-cyan-400/20 to-green-400/20 blur-3xl animate-pulse" style="animation-delay: 1s;"></div>
      <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-96 h-96 rounded-full bg-gradient-to-r from-teal-400/10 to-emerald-400/10 blur-3xl animate-pulse" style="animation-delay: 2s;"></div>
      
      <!-- Floating particles -->
      <div class="absolute top-20 left-20 w-2 h-2 bg-emerald-400/30 rounded-full animate-float"></div>
      <div class="absolute top-40 right-32 w-1 h-1 bg-teal-400/40 rounded-full animate-float" style="animation-delay: 0.5s;"></div>
      <div class="absolute bottom-32 left-48 w-1.5 h-1.5 bg-cyan-400/35 rounded-full animate-float" style="animation-delay: 1.2s;"></div>
      <div class="absolute top-60 right-20 w-1 h-1 bg-green-400/40 rounded-full animate-float" style="animation-delay: 2.1s;"></div>
      <div class="absolute bottom-60 right-40 w-2 h-2 bg-emerald-400/30 rounded-full animate-float" style="animation-delay: 1.8s;"></div>
      
      <!-- Grid pattern overlay -->
      <div class="absolute inset-0 bg-grid-pattern opacity-5"></div>
    </div>
    
    <div class="max-w-7xl w-full relative z-10 container-modern">
      <div class="lg:grid-cols-2 lg:grid gap-12 items-center min-h-screen py-12">
        <!-- Left Panel - Hero Content -->
        <div class="hidden lg:flex flex-col justify-center">
          <div class="max-w-md">
            <div class="hero-badge">
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
              </svg>
              Join Our Team
            </div>
            <h1 class="hero-title mt-3 text-4xl font-bold">
              Start Managing Incidents Today
            </h1>
            <p class="hero-desc mt-4 text-lg text-slate-600">
              Create your account to access powerful incident management tools. 
              Join organizations worldwide using our platform to maintain system reliability and respond to incidents efficiently.
            </p>
            <div class="mt-8 space-y-4">
              <div class="flex items-center space-x-3">
                <div class="flex-shrink-0">
                  <div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                    <svg class="w-4 h-4 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
                <span class="text-slate-700">Free account with full features</span>
              </div>
              <div class="flex items-center space-x-3">
                <div class="flex-shrink-0">
                  <div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                    <svg class="w-4 h-4 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
                <span class="text-slate-700">Role-based access control</span>
              </div>
              <div class="flex items-center space-x-3">
                <div class="flex-shrink-0">
                  <div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                    <svg class="w-4 h-4 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
                <span class="text-slate-700">24/7 incident monitoring</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Panel - Enhanced Registration Form -->
        <div class="flex flex-col justify-center animate-slide-in-right">
          <!-- Mobile logo/branding with enhanced animation -->
          <div class="lg:hidden text-center mb-10 animate-bounce-in">
            <div class="inline-flex items-center justify-center w-18 h-18 rounded-3xl bg-gradient-to-tr from-emerald-600 via-green-700 to-teal-700 shadow-xl mb-6 relative overflow-hidden group">
              <!-- Shimmer effect -->
              <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -skew-x-12 transform translate-x-full group-hover:translate-x-[-200%] transition-transform duration-1000"></div>
              <svg class="w-10 h-10 text-white relative z-10" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
              </svg>
            </div>
            <h1 class="text-3xl font-bold bg-gradient-to-r from-slate-900 to-emerald-900 bg-clip-text text-transparent">Join Our Team</h1>
            <p class="text-slate-600 text-base mt-2">Start Your Journey Today</p>
          </div>
          
          <div class="card-modern max-w-lg mx-auto w-full">
            <!-- Enhanced Header with gradient text -->
            <div class="text-center mb-10">
              <h2 class="text-4xl font-bold bg-gradient-to-r from-slate-900 via-emerald-900 to-teal-900 bg-clip-text text-transparent mb-3">Create Account</h2>
              <p class="text-lg text-slate-600">Join thousands of teams managing incidents effectively</p>
              <!-- Progress indicator for registration -->
              <div class="mt-4 flex justify-center">
                <div class="flex space-x-2">
                  <div class="w-8 h-1 bg-slate-200 rounded-full"></div>
                  <div class="w-8 h-1 bg-emerald-600 rounded-full"></div>
                  <div class="w-8 h-1 bg-slate-200 rounded-full"></div>
                </div>
              </div>
            </div>

            <!-- Error Message -->
            <div v-if="authStore.error" class="error-box mb-6" role="alert" aria-live="polite">
              {{ authStore.error }}
            </div>

            <!-- Registration Form -->
            <form @submit.prevent="handleRegister" class="space-y-5">
              <!-- Full Name Field -->
              <div class="form-group">
                <div class="relative">
                  <input
                    id="fullName"
                    v-model="registerForm.fullName"
                    type="text"
                    autocomplete="name"
                    required
                    :class="[
                      'w-full px-4 py-4 pt-6 border rounded-xl focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all duration-200 peer placeholder-transparent bg-white/50 backdrop-blur-sm',
                      validation.errors.value.fullName 
                        ? 'border-red-300 bg-red-50/50' 
                        : 'border-slate-200 hover:border-slate-300'
                    ]"
                    :aria-invalid="!!validation.errors.value.fullName"
                    :aria-describedby="validation.errors.value.fullName ? 'fullName-error' : undefined"
                    placeholder="Enter your full name"
                    @input="validation.clearError('fullName')"
                  />
                  <label 
                    for="fullName" 
                    class="absolute left-4 top-2 text-xs font-medium text-slate-600 transition-all duration-200 peer-placeholder-shown:text-base peer-placeholder-shown:text-slate-400 peer-placeholder-shown:top-4 peer-focus:top-2 peer-focus:text-xs peer-focus:text-green-600 peer-focus:font-medium"
                  >
                    Full Name *
                  </label>
                </div>
                <p v-if="validation.errors.value.fullName" id="fullName-error" class="mt-2 text-sm text-red-600 flex items-center">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                  {{ validation.errors.value.fullName }}
                </p>
              </div>

              <!-- Username Field -->
              <div class="form-group">
                <label for="username" class="block text-sm font-medium text-slate-700 mb-2">
                  Username *
                </label>
                <input
                  id="username"
                  v-model="registerForm.username"
                  type="text"
                  autocomplete="username"
                  required
                  :class="[
                    'w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200',
                    validation.errors.value.username 
                      ? 'border-red-300 bg-red-50' 
                      : 'border-slate-300 bg-white hover:border-slate-400'
                  ]"
                  :aria-invalid="!!validation.errors.value.username"
                  :aria-describedby="validation.errors.value.username ? 'username-error' : undefined"
                  placeholder="Choose a username"
                  @input="validation.clearError('username')"
                />
                <p v-if="validation.errors.value.username" id="username-error" class="mt-1 text-sm text-red-600">
                  {{ validation.errors.value.username }}
                </p>
              </div>

              <!-- Email Field -->
              <div class="form-group">
                <label for="email" class="block text-sm font-medium text-slate-700 mb-2">
                  Email Address *
                </label>
                <input
                  id="email"
                  v-model="registerForm.email"
                  type="email"
                  autocomplete="email"
                  required
                  :class="[
                    'w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200',
                    validation.errors.value.email 
                      ? 'border-red-300 bg-red-50' 
                      : 'border-slate-300 bg-white hover:border-slate-400'
                  ]"
                  :aria-invalid="!!validation.errors.value.email"
                  :aria-describedby="validation.errors.value.email ? 'email-error' : undefined"
                  placeholder="Enter your email address"
                  @input="validation.clearError('email')"
                />
                <p v-if="validation.errors.value.email" id="email-error" class="mt-1 text-sm text-red-600">
                  {{ validation.errors.value.email }}
                </p>
              </div>

              <!-- Password Field -->
              <div class="form-group">
                <label for="password" class="block text-sm font-medium text-slate-700 mb-2">
                  Password *
                </label>
                <div class="relative">
                  <input
                    id="password"
                    v-model="registerForm.password"
                    :type="showPassword ? 'text' : 'password'"
                    autocomplete="new-password"
                    required
                    :class="[
                      'w-full px-4 py-3 pr-12 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200',
                      validation.errors.value.password 
                        ? 'border-red-300 bg-red-50' 
                        : 'border-slate-300 bg-white hover:border-slate-400'
                    ]"
                    :aria-invalid="!!validation.errors.value.password"
                    :aria-describedby="validation.errors.value.password ? 'password-error password-help' : 'password-help'"
                    placeholder="Create a strong password"
                    @input="validation.clearError('password')"
                  />
                  <button
                    type="button"
                    @click="showPassword = !showPassword"
                    class="absolute inset-y-0 right-0 flex items-center pr-3 text-slate-400 hover:text-slate-600"
                    :aria-label="showPassword ? 'Hide password' : 'Show password'"
                  >
                    <svg v-if="!showPassword" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                      <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                    </svg>
                    <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.958 9.958 0 00-4.512 1.074l-1.78-1.781zm4.261 4.26l1.514 1.515a2.003 2.003 0 012.45 2.45l1.514 1.514a4 4 0 00-5.478-5.478z" clip-rule="evenodd" />
                      <path d="M12.454 16.697L9.75 13.992a4 4 0 01-3.742-3.741L2.335 6.578A9.98 9.98 0 00.458 10c1.274 4.057 5.065 7 9.542 7 .847 0 1.669-.105 2.454-.303z" />
                    </svg>
                  </button>
                </div>
                <p id="password-help" class="mt-1 text-xs text-slate-500">
                  Must be at least 8 characters with a mix of letters, numbers, and symbols
                </p>
                <p v-if="validation.errors.value.password" id="password-error" class="mt-1 text-sm text-red-600">
                  {{ validation.errors.value.password }}
                </p>
              </div>

              <!-- Confirm Password Field -->
              <div class="form-group">
                <label for="confirmPassword" class="block text-sm font-medium text-slate-700 mb-2">
                  Confirm Password *
                </label>
                <div class="relative">
                  <input
                    id="confirmPassword"
                    v-model="registerForm.confirmPassword"
                    :type="showConfirmPassword ? 'text' : 'password'"
                    autocomplete="new-password"
                    required
                    :class="[
                      'w-full px-4 py-3 pr-12 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200',
                      validation.errors.value.confirmPassword 
                        ? 'border-red-300 bg-red-50' 
                        : 'border-slate-300 bg-white hover:border-slate-400'
                    ]"
                    :aria-invalid="!!validation.errors.value.confirmPassword"
                    :aria-describedby="validation.errors.value.confirmPassword ? 'confirmPassword-error' : undefined"
                    placeholder="Confirm your password"
                    @input="validation.clearError('confirmPassword')"
                  />
                  <button
                    type="button"
                    @click="showConfirmPassword = !showConfirmPassword"
                    class="absolute inset-y-0 right-0 flex items-center pr-3 text-slate-400 hover:text-slate-600"
                    :aria-label="showConfirmPassword ? 'Hide password confirmation' : 'Show password confirmation'"
                  >
                    <svg v-if="!showConfirmPassword" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                      <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                    </svg>
                    <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.958 9.958 0 00-4.512 1.074l-1.78-1.781zm4.261 4.26l1.514 1.515a2.003 2.003 0 012.45 2.45l1.514 1.514a4 4 0 00-5.478-5.478z" clip-rule="evenodd" />
                      <path d="M12.454 16.697L9.75 13.992a4 4 0 01-3.742-3.741L2.335 6.578A9.98 9.98 0 00.458 10c1.274 4.057 5.065 7 9.542 7 .847 0 1.669-.105 2.454-.303z" />
                    </svg>
                  </button>
                </div>
                <p v-if="validation.errors.value.confirmPassword" id="confirmPassword-error" class="mt-1 text-sm text-red-600">
                  {{ validation.errors.value.confirmPassword }}
                </p>
              </div>

              <!-- Terms and Privacy -->
              <div class="form-group">
                <label class="flex items-start group cursor-pointer">
                  <div class="relative">
                    <input
                      v-model="registerForm.acceptTerms"
                      type="checkbox"
                      required
                      class="sr-only"
                      :aria-describedby="validation.errors.value.acceptTerms ? 'terms-error' : undefined"
                    />
                    <div class="w-5 h-5 rounded border-2 border-slate-300 group-hover:border-green-500 transition-colors duration-200 flex items-center justify-center" :class="registerForm.acceptTerms ? 'bg-green-600 border-green-600' : 'bg-white'">
                      <svg v-if="registerForm.acceptTerms" class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                      </svg>
                    </div>
                  </div>
                  <span class="ml-3 text-sm text-slate-700 leading-5 group-hover:text-slate-900 transition-colors duration-200">
                    I agree to the 
                    <a href="#" class="text-green-600 hover:text-green-700 hover:underline font-medium">Terms of Service</a>
                    and 
                    <a href="#" class="text-green-600 hover:text-green-700 hover:underline font-medium">Privacy Policy</a>
                  </span>
                </label>
                <p v-if="validation.errors.value.acceptTerms" id="terms-error" class="mt-2 text-sm text-red-600 flex items-center">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                  {{ validation.errors.value.acceptTerms }}
                </p>
              </div>

              <!-- Submit Button -->
              <button
                type="submit"
                :disabled="authStore.isLoading || !isFormValid"
                class="btn-primary w-full py-4 px-4 text-base font-semibold rounded-xl shadow-lg disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center transition-all duration-300 transform hover:scale-105 active:scale-95 bg-gradient-to-r from-green-600 via-green-700 to-emerald-700 hover:from-green-700 hover:via-green-800 hover:to-emerald-800"
                :aria-describedby="authStore.isLoading ? 'register-status' : undefined"
              >
                <svg
                  v-if="authStore.isLoading"
                  class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ authStore.isLoading ? 'Creating Account...' : 'Create Account' }}
              </button>
              <span v-if="authStore.isLoading" id="register-status" class="sr-only">Creating account, please wait</span>
            </form>

            <!-- Login Link -->
            <div class="mt-8 text-center">
              <p class="text-sm text-slate-600">
                Already have an account?
                <router-link
                  to="/login"
                  class="text-green-600 hover:text-green-700 hover:underline font-semibold focus:outline-none focus:underline ml-1 transition-colors duration-200"
                >
                  Sign in here
                </router-link>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useFormValidation } from '@/composables/useFormValidation'

const router = useRouter()
const authStore = useAuthStore()
const validation = useFormValidation()

// Form state
const registerForm = ref({
  fullName: '',
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  acceptTerms: false
})

const showPassword = ref(false)
const showConfirmPassword = ref(false)

// Form validation
const isFormValid = computed(() => {
  return registerForm.value.fullName.trim() !== '' &&
         registerForm.value.username.trim() !== '' &&
         registerForm.value.email.trim() !== '' &&
         registerForm.value.password.trim() !== '' &&
         registerForm.value.confirmPassword.trim() !== '' &&
         registerForm.value.acceptTerms &&
         validation.isValid.value
})

// Password strength validation
const validatePasswordStrength = (password: string): string | null => {
  if (password.length < 8) {
    return 'Password must be at least 8 characters long'
  }
  
  const hasLower = /[a-z]/.test(password)
  const hasUpper = /[A-Z]/.test(password)
  const hasNumber = /\d/.test(password)
  const hasSpecial = /[!@#$%^&*(),.?":{}|<>]/.test(password)
  
  const strength = [hasLower, hasUpper, hasNumber, hasSpecial].filter(Boolean).length
  
  if (strength < 3) {
    return 'Password should contain a mix of uppercase, lowercase, numbers, and symbols'
  }
  
  return null
}

// Password confirmation validation
const validatePasswordConfirmation = (confirmPassword: string): string | null => {
  if (confirmPassword !== registerForm.value.password) {
    return 'Passwords do not match'
  }
  return null
}

// Handle form submission
const handleRegister = async () => {
  // Clear previous errors
  validation.clearAllErrors()

  // Validate form
  const isValid = validation.validate(registerForm.value, {
    fullName: { required: true, minLength: 2, maxLength: 100 },
    username: { required: true, minLength: 3, maxLength: 50 },
    email: { required: true, email: true },
    password: { 
      required: true, 
      minLength: 8,
      custom: validatePasswordStrength
    },
    confirmPassword: {
      required: true,
      custom: validatePasswordConfirmation
    },
    acceptTerms: {
      required: true,
      custom: (value: any) => !value ? 'You must accept the terms and privacy policy' : null
    }
  })

  if (!isValid) {
    return
  }

  try {
    const registerData = {
      full_name: registerForm.value.fullName,
      username: registerForm.value.username,
      email: registerForm.value.email,
      password: registerForm.value.password
    }

    await authStore.register(registerData)
    
    // Redirect to dashboard on success
    router.push('/')
  } catch (error) {
    // Error is already handled by the auth store
    console.error('Registration failed:', error)
  }
}

// Focus management
onMounted(() => {
  // Focus the first input field for better UX
  const fullNameInput = document.getElementById('fullName')
  if (fullNameInput) {
    fullNameInput.focus()
  }
})
</script>

<style scoped>
/* Modern Web Platform Animations - Register Version */
@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(180deg); }
}

@keyframes slide-up {
  from { 
    opacity: 0; 
    transform: translateY(30px); 
  }
  to { 
    opacity: 1; 
    transform: translateY(0); 
  }
}

@keyframes slide-in-left {
  from { 
    opacity: 0; 
    transform: translateX(-50px); 
  }
  to { 
    opacity: 1; 
    transform: translateX(0); 
  }
}

@keyframes slide-in-right {
  from { 
    opacity: 0; 
    transform: translateX(50px); 
  }
  to { 
    opacity: 1; 
    transform: translateX(0); 
  }
}

@keyframes fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes bounce-in {
  0% {
    opacity: 0;
    transform: scale(0.3);
  }
  50% {
    opacity: 1;
    transform: scale(1.05);
  }
  70% {
    transform: scale(0.9);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes text-gradient {
  0%, 100% {
    background-size: 200% 200%;
    background-position: left center;
  }
  50% {
    background-size: 200% 200%;
    background-position: right center;
  }
}

/* Apply animations */
.animate-float {
  animation: float 6s ease-in-out infinite;
}

.animate-slide-up {
  animation: slide-up 0.8s ease-out both;
}

.animate-slide-in-left {
  animation: slide-in-left 0.8s ease-out both;
}

.animate-slide-in-right {
  animation: slide-in-right 0.8s ease-out both;
}

.animate-fade-in {
  animation: fade-in 1s ease-out both;
}

.animate-bounce-in {
  animation: bounce-in 1s ease-out both;
}

.animate-text-gradient {
  animation: text-gradient 3s ease infinite;
  background-size: 200% auto;
}

/* Modern container system */
.container-modern {
  @apply px-6 sm:px-8 lg:px-12;
}

/* Enhanced card with modern effects - Emerald theme */
.card-modern {
  @apply bg-white/80 backdrop-blur-xl rounded-3xl shadow-2xl p-10 border border-white/30;
  box-shadow: 
    0 32px 64px -12px rgba(0, 0, 0, 0.08),
    0 25px 50px -12px rgba(16, 185, 129, 0.15),
    0 0 0 1px rgba(255, 255, 255, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.5);
  position: relative;
  overflow: hidden;
}

.card-modern::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(16, 185, 129, 0.5), transparent);
}

/* Grid background pattern - Emerald theme */
.bg-grid-pattern {
  background-image: 
    linear-gradient(rgba(16, 185, 129, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(16, 185, 129, 0.1) 1px, transparent 1px);
  background-size: 20px 20px;
}

/* Enhanced component styles */
.hero-badge {
  @apply inline-flex items-center bg-gradient-to-r from-emerald-50 to-teal-50 text-emerald-700 px-4 py-2 rounded-2xl text-sm font-semibold border border-emerald-200/50 shadow-sm;
  backdrop-filter: blur(10px);
}

.hero-badge > svg {
  @apply mr-2;
}

.hero-title {
  @apply text-slate-900 leading-tight;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  font-variation-settings: 'wght' 800;
}

.hero-desc {
  @apply text-slate-600 max-w-2xl;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}

.form-group {
  @apply space-y-2;
}

.error-box {
  @apply bg-red-50/90 backdrop-blur-sm border border-red-200/60 text-red-800 px-5 py-4 rounded-2xl text-sm shadow-sm;
  position: relative;
}

.error-box::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, #ef4444, #dc2626);
  border-radius: 1rem 1rem 0 0;
}

.btn-primary {
  @apply bg-gradient-to-r from-emerald-600 via-green-700 to-teal-700 text-white font-semibold rounded-2xl shadow-xl hover:shadow-2xl transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2;
  transform: translateY(0);
  position: relative;
  overflow: hidden;
}

.btn-primary::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.btn-primary:hover:not(:disabled)::before {
  left: 100%;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-3px) scale(1.02);
  box-shadow: 0 25px 50px -12px rgba(16, 185, 129, 0.4);
}

.btn-primary:active:not(:disabled) {
  transform: translateY(-1px) scale(0.98);
}

/* Enhanced input styling with modern focus states */
input[type="text"], input[type="password"], input[type="email"] {
  backdrop-filter: blur(12px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

input:focus {
  box-shadow: 
    0 0 0 3px rgba(16, 185, 129, 0.12), 
    0 10px 25px -5px rgba(16, 185, 129, 0.25),
    0 0 0 1px rgba(16, 185, 129, 0.5);
  transform: translateY(-2px);
}

input:hover:not(:focus) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px -2px rgba(0, 0, 0, 0.1);
}

/* Floating label enhancements */
label {
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  font-variation-settings: 'wght' 500;
}

/* Custom checkbox styling - Enhanced */
.checkbox-group {
  position: relative;
}

/* Custom scrollbar for modern look - Emerald theme */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(148, 163, 184, 0.1);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, rgba(16, 185, 129, 0.8), rgba(20, 184, 166, 0.8));
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, rgba(16, 185, 129, 1), rgba(20, 184, 166, 1));
}

/* Responsive enhancements */
@media (max-width: 1023px) {
  .lg\:grid-cols-2 {
    grid-template-columns: 1fr;
  }
  
  .hidden.lg\:flex {
    display: none;
  }
  
  .card-modern {
    @apply mx-4;
  }
}

@media (max-width: 640px) {
  .card-modern {
    @apply p-8 mx-4;
  }
  
  .btn-primary {
    @apply py-4 text-lg;
  }
  
  input[type="text"], input[type="password"], input[type="email"] {
    @apply py-4 text-base;
  }
  
  .form-group {
    @apply space-y-2;
  }
  
  .container-modern {
    @apply px-4;
  }
}

/* Reduced motion for accessibility */
@media (prefers-reduced-motion: reduce) {
  .animate-float,
  .animate-slide-up,
  .animate-slide-in-left,
  .animate-slide-in-right,
  .animate-fade-in,
  .animate-bounce-in,
  .animate-text-gradient {
    animation: none;
  }
  
  input:focus {
    transform: none;
  }
  
  .btn-primary:hover:not(:disabled) {
    transform: none;
  }
}
</style>