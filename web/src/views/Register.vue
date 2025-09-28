<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 flex items-center justify-center px-4 sm:px-6 lg:px-8">
    <div class="max-w-5xl w-full">
      <div class="lg:grid-cols-2 lg:grid gap-8">
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

        <!-- Right Panel - Registration Form -->
        <div class="flex flex-col justify-center">
          <div class="card max-w-md mx-auto w-full">
            <!-- Header -->
            <div class="text-center mb-8">
              <h2 class="text-3xl font-bold text-slate-900">Create Account</h2>
              <p class="mt-2 text-slate-600">Join thousands of teams managing incidents effectively</p>
            </div>

            <!-- Error Message -->
            <div v-if="authStore.error" class="error-box mb-6" role="alert" aria-live="polite">
              {{ authStore.error }}
            </div>

            <!-- Registration Form -->
            <form @submit.prevent="handleRegister" class="space-y-6">
              <!-- Full Name Field -->
              <div class="form-group">
                <label for="fullName" class="block text-sm font-medium text-slate-700 mb-2">
                  Full Name *
                </label>
                <input
                  id="fullName"
                  v-model="registerForm.fullName"
                  type="text"
                  autocomplete="name"
                  required
                  :class="[
                    'w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200',
                    validation.errors.value.fullName 
                      ? 'border-red-300 bg-red-50' 
                      : 'border-slate-300 bg-white hover:border-slate-400'
                  ]"
                  :aria-invalid="!!validation.errors.value.fullName"
                  :aria-describedby="validation.errors.value.fullName ? 'fullName-error' : undefined"
                  placeholder="Enter your full name"
                  @input="validation.clearError('fullName')"
                />
                <p v-if="validation.errors.value.fullName" id="fullName-error" class="mt-1 text-sm text-red-600">
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
                <label class="flex items-start space-x-3">
                  <input
                    v-model="registerForm.acceptTerms"
                    type="checkbox"
                    required
                    class="mt-1 w-4 h-4 text-blue-600 border-slate-300 rounded focus:ring-blue-500 focus:ring-2"
                    :aria-describedby="validation.errors.value.acceptTerms ? 'terms-error' : undefined"
                  />
                  <span class="text-sm text-slate-700 leading-5">
                    I agree to the 
                    <a href="#" class="text-blue-600 hover:text-blue-700 hover:underline font-medium">Terms of Service</a>
                    and 
                    <a href="#" class="text-blue-600 hover:text-blue-700 hover:underline font-medium">Privacy Policy</a>
                  </span>
                </label>
                <p v-if="validation.errors.value.acceptTerms" id="terms-error" class="mt-1 text-sm text-red-600">
                  {{ validation.errors.value.acceptTerms }}
                </p>
              </div>

              <!-- Submit Button -->
              <button
                type="submit"
                :disabled="authStore.isLoading || !isFormValid"
                class="btn-primary w-full py-3 px-4 text-base font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
                :aria-describedby="authStore.isLoading ? 'register-status' : undefined"
              >
                <svg
                  v-if="authStore.isLoading"
                  class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
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
            <div class="mt-6 text-center">
              <p class="text-sm text-slate-600">
                Already have an account?
                <router-link
                  to="/login"
                  class="text-blue-600 hover:text-blue-700 hover:underline font-medium focus:outline-none focus:underline"
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
/* Component-specific styles if needed */
.hero-badge {
  @apply inline-flex items-center bg-green-50 text-green-700 px-3 py-1.5 rounded-full text-sm font-semibold;
}

.hero-badge > svg {
  @apply mr-2;
}

.hero-title {
  @apply text-slate-900 leading-tight;
}

.hero-desc {
  @apply text-slate-600 max-w-lg;
}

.card {
  @apply bg-white rounded-2xl shadow-xl p-8 border border-slate-100;
}

.form-group {
  @apply space-y-1;
}

.error-box {
  @apply bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg text-sm;
}

.btn-primary {
  @apply bg-gradient-to-r from-blue-600 to-blue-700 text-white font-semibold rounded-lg shadow-lg hover:shadow-xl hover:from-blue-700 hover:to-blue-800 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2;
}

/* Responsive adjustments */
@media (max-width: 1023px) {
  .lg\:grid-cols-2 {
    grid-template-columns: 1fr;
  }
  
  .hidden.lg\:flex {
    display: none;
  }
}
</style>