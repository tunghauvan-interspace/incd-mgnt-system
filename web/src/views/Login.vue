<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center px-4 sm:px-6 lg:px-8 relative overflow-hidden">
    <!-- Background decorative elements -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-40 -right-40 w-80 h-80 rounded-full bg-gradient-to-br from-blue-400/20 to-purple-400/20 blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 rounded-full bg-gradient-to-tr from-indigo-400/20 to-pink-400/20 blur-3xl"></div>
      <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-96 h-96 rounded-full bg-gradient-to-r from-cyan-400/10 to-blue-400/10 blur-3xl"></div>
    </div>
    
    <div class="max-w-5xl w-full relative z-10">
      <div class="lg:grid-cols-2 lg:grid gap-8">
        <!-- Left Panel - Hero Content -->
        <div class="hidden lg:flex flex-col justify-center">
          <div class="max-w-md">
            <div class="hero-badge">
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
              Secure & Reliable
            </div>
            <h1 class="hero-title mt-3 text-4xl font-bold">
              Access Your Incident Management Dashboard
            </h1>
            <p class="hero-desc mt-4 text-lg text-slate-600">
              Monitor, manage, and resolve incidents efficiently with our comprehensive incident management platform. 
              Get real-time alerts, track response times, and maintain system reliability.
            </p>
            <div class="mt-8 space-y-4">
              <div class="flex items-center space-x-3">
                <div class="flex-shrink-0">
                  <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                    <svg class="w-4 h-4 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
                <span class="text-slate-700">Real-time incident monitoring</span>
              </div>
              <div class="flex items-center space-x-3">
                <div class="flex-shrink-0">
                  <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                    <svg class="w-4 h-4 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
                <span class="text-slate-700">Automated alert aggregation</span>
              </div>
              <div class="flex items-center space-x-3">
                <div class="flex-shrink-0">
                  <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                    <svg class="w-4 h-4 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
                <span class="text-slate-700">Team collaboration tools</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Panel - Login Form -->
        <div class="flex flex-col justify-center">
          <!-- Mobile logo/branding -->
          <div class="lg:hidden text-center mb-8">
            <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-tr from-blue-600 to-purple-600 shadow-lg mb-4">
              <svg class="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
            </div>
            <h1 class="text-2xl font-bold text-slate-900">Incident Management</h1>
            <p class="text-slate-600 text-sm">Secure & Reliable Platform</p>
          </div>
          
          <div class="card max-w-md mx-auto w-full backdrop-blur-lg bg-white/70 border border-white/20 shadow-2xl">
            <!-- Header -->
            <div class="text-center mb-8">
              <h2 class="text-3xl font-bold bg-gradient-to-r from-slate-900 to-slate-700 bg-clip-text text-transparent">Welcome Back</h2>
              <p class="mt-2 text-slate-600">Sign in to your account to continue</p>
            </div>

            <!-- Error Message -->
            <div v-if="authStore.error" class="error-box mb-6" role="alert" aria-live="polite">
              {{ authStore.error }}
            </div>

            <!-- Login Form -->
            <form @submit.prevent="handleLogin" class="space-y-6">
              <!-- Username/Email Field -->
              <div class="form-group">
                <div class="relative">
                  <input
                    id="username"
                    v-model="loginForm.identifier"
                    type="text"
                    autocomplete="username"
                    required
                    :class="[
                      'w-full px-4 py-4 pt-6 border rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 peer placeholder-transparent bg-white/50 backdrop-blur-sm',
                      validation.errors.value.identifier 
                        ? 'border-red-300 bg-red-50/50' 
                        : 'border-slate-200 hover:border-slate-300'
                    ]"
                    :aria-invalid="!!validation.errors.value.identifier"
                    :aria-describedby="validation.errors.value.identifier ? 'username-error' : undefined"
                    placeholder="Enter your email or username"
                    @input="validation.clearError('identifier')"
                  />
                  <label 
                    for="username" 
                    class="absolute left-4 top-2 text-xs font-medium text-slate-600 transition-all duration-200 peer-placeholder-shown:text-base peer-placeholder-shown:text-slate-400 peer-placeholder-shown:top-4 peer-focus:top-2 peer-focus:text-xs peer-focus:text-blue-600 peer-focus:font-medium"
                  >
                    Email or Username
                  </label>
                </div>
                <p v-if="validation.errors.value.identifier" id="username-error" class="mt-2 text-sm text-red-600 flex items-center">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                  {{ validation.errors.value.identifier }}
                </p>
              </div>

              <!-- Password Field -->
              <div class="form-group">
                <div class="relative">
                  <input
                    id="password"
                    v-model="loginForm.password"
                    :type="showPassword ? 'text' : 'password'"
                    autocomplete="current-password"
                    required
                    :class="[
                      'w-full px-4 py-4 pt-6 pr-12 border rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 peer placeholder-transparent bg-white/50 backdrop-blur-sm',
                      validation.errors.value.password 
                        ? 'border-red-300 bg-red-50/50' 
                        : 'border-slate-200 hover:border-slate-300'
                    ]"
                    :aria-invalid="!!validation.errors.value.password"
                    :aria-describedby="validation.errors.value.password ? 'password-error' : undefined"
                    placeholder="Enter your password"
                    @input="validation.clearError('password')"
                  />
                  <label 
                    for="password" 
                    class="absolute left-4 top-2 text-xs font-medium text-slate-600 transition-all duration-200 peer-placeholder-shown:text-base peer-placeholder-shown:text-slate-400 peer-placeholder-shown:top-4 peer-focus:top-2 peer-focus:text-xs peer-focus:text-blue-600 peer-focus:font-medium"
                  >
                    Password
                  </label>
                  <button
                    type="button"
                    @click="showPassword = !showPassword"
                    class="absolute inset-y-0 right-0 flex items-center pr-4 text-slate-400 hover:text-slate-600 transition-colors duration-200"
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
                <p v-if="validation.errors.value.password" id="password-error" class="mt-2 text-sm text-red-600 flex items-center">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                  {{ validation.errors.value.password }}
                </p>
              </div>

              <!-- Remember Me & Forgot Password -->
              <div class="flex items-center justify-between">
                <label class="flex items-center group cursor-pointer">
                  <div class="relative">
                    <input
                      v-model="loginForm.rememberMe"
                      type="checkbox"
                      class="sr-only"
                    />
                    <div class="w-5 h-5 rounded border-2 border-slate-300 group-hover:border-blue-500 transition-colors duration-200 flex items-center justify-center" :class="loginForm.rememberMe ? 'bg-blue-600 border-blue-600' : 'bg-white'">
                      <svg v-if="loginForm.rememberMe" class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                      </svg>
                    </div>
                  </div>
                  <span class="ml-3 text-sm text-slate-700 group-hover:text-slate-900 transition-colors duration-200">Remember me</span>
                </label>
                <router-link
                  to="/forgot-password"
                  class="text-sm text-blue-600 hover:text-blue-700 hover:underline focus:outline-none focus:underline font-medium transition-colors duration-200"
                >
                  Forgot password?
                </router-link>
              </div>

              <!-- Submit Button -->
              <button
                type="submit"
                :disabled="authStore.isLoading || !isFormValid"
                class="btn-primary w-full py-4 px-4 text-base font-semibold rounded-xl shadow-lg disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center transition-all duration-300 transform hover:scale-105 active:scale-95 bg-gradient-to-r from-blue-600 via-blue-700 to-indigo-700 hover:from-blue-700 hover:via-blue-800 hover:to-indigo-800"
                :aria-describedby="authStore.isLoading ? 'login-status' : undefined"
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
                {{ authStore.isLoading ? 'Signing In...' : 'Sign In' }}
              </button>
              <span v-if="authStore.isLoading" id="login-status" class="sr-only">Signing in, please wait</span>
            </form>

            <!-- Register Link -->
            <div class="mt-8 text-center">
              <p class="text-sm text-slate-600">
                Don't have an account?
                <router-link
                  to="/register"
                  class="text-blue-600 hover:text-blue-700 hover:underline font-semibold focus:outline-none focus:underline ml-1 transition-colors duration-200"
                >
                  Sign up here
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
const loginForm = ref({
  identifier: '', // can be username or email
  password: '',
  rememberMe: false
})

const showPassword = ref(false)

// Form validation
const isFormValid = computed(() => {
  return loginForm.value.identifier.trim() !== '' && 
         loginForm.value.password.trim() !== '' &&
         validation.isValid.value
})

// Handle form submission
const handleLogin = async () => {
  // Clear previous errors
  validation.clearAllErrors()

  // Validate form
  const isValid = validation.validate(loginForm.value, {
    identifier: { required: true, minLength: 3 },
    password: { required: true, minLength: 6 }
  })

  if (!isValid) {
    return
  }

  try {
    const identifier = loginForm.value.identifier.trim()
    
    // Determine if identifier is email or username
    const isEmail = identifier.includes('@')
    const loginData = isEmail 
      ? { email: identifier, password: loginForm.value.password, username: '' }
      : { username: identifier, password: loginForm.value.password, email: '' }

    await authStore.login(loginData)
    
    // Redirect to dashboard on success
    router.push('/')
  } catch (error) {
    // Error is already handled by the auth store
    console.error('Login failed:', error)
  }
}

// Focus management
onMounted(() => {
  // Focus the first input field for better UX
  const usernameInput = document.getElementById('username')
  if (usernameInput) {
    usernameInput.focus()
  }
})
</script>

<style scoped>
/* Component-specific styles */
.hero-badge {
  @apply inline-flex items-center bg-blue-50 text-blue-700 px-3 py-1.5 rounded-full text-sm font-semibold;
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
  @apply bg-white/70 backdrop-blur-lg rounded-3xl shadow-2xl p-8 border border-white/20;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.1), 0 0 0 1px rgba(255, 255, 255, 0.2);
}

.form-group {
  @apply space-y-1;
}

.error-box {
  @apply bg-red-50/80 backdrop-blur-sm border border-red-200/50 text-red-800 px-4 py-3 rounded-xl text-sm shadow-sm;
}

.btn-primary {
  @apply bg-gradient-to-r from-blue-600 to-indigo-600 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl hover:from-blue-700 hover:to-indigo-700 transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2;
  transform: translateY(0);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 25px 50px -12px rgba(59, 130, 246, 0.25);
}

.btn-primary:active:not(:disabled) {
  transform: translateY(0);
}

/* Enhanced input styling */
input[type="text"], input[type="password"], input[type="email"] {
  backdrop-filter: blur(10px);
}

input:focus {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1), 0 10px 25px -5px rgba(0, 0, 0, 0.1);
}

/* Responsive adjustments */
@media (max-width: 1023px) {
  .lg\:grid-cols-2 {
    grid-template-columns: 1fr;
  }
  
  .hidden.lg\:flex {
    display: none;
  }
  
  .card {
    @apply mx-4;
  }
}

/* Mobile-specific enhancements */
@media (max-width: 640px) {
  .card {
    @apply p-6 mx-3;
  }
  
  .btn-primary {
    @apply py-4 text-lg;
  }
  
  input[type="text"], input[type="password"], input[type="email"] {
    @apply py-4 text-base;
  }
}
</style>