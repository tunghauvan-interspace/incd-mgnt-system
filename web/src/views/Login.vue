<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center px-4 sm:px-6 lg:px-8 relative overflow-hidden">
    <!-- Enhanced background with animated particles -->
    <div class="absolute inset-0 overflow-hidden">
      <!-- Main gradient blobs -->
      <div class="absolute -top-40 -right-40 w-80 h-80 rounded-full bg-gradient-to-br from-blue-400/20 to-purple-400/20 blur-3xl animate-pulse"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 rounded-full bg-gradient-to-tr from-indigo-400/20 to-pink-400/20 blur-3xl animate-pulse" style="animation-delay: 1s;"></div>
      <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-96 h-96 rounded-full bg-gradient-to-r from-cyan-400/10 to-blue-400/10 blur-3xl animate-pulse" style="animation-delay: 2s;"></div>
      
      <!-- Floating particles -->
      <div class="absolute top-20 left-20 w-2 h-2 bg-blue-400/30 rounded-full animate-float"></div>
      <div class="absolute top-40 right-32 w-1 h-1 bg-purple-400/40 rounded-full animate-float" style="animation-delay: 0.5s;"></div>
      <div class="absolute bottom-32 left-48 w-1.5 h-1.5 bg-indigo-400/35 rounded-full animate-float" style="animation-delay: 1.2s;"></div>
      <div class="absolute top-60 right-20 w-1 h-1 bg-cyan-400/40 rounded-full animate-float" style="animation-delay: 2.1s;"></div>
      <div class="absolute bottom-60 right-40 w-2 h-2 bg-pink-400/30 rounded-full animate-float" style="animation-delay: 1.8s;"></div>
      
      <!-- Grid pattern overlay -->
      <div class="absolute inset-0 bg-grid-pattern opacity-5"></div>
    </div>
    
    <div class="max-w-7xl w-full relative z-10 container-modern">
      <div class="lg:grid-cols-2 lg:grid gap-12 items-center min-h-screen py-12">
        <!-- Left Panel - Enhanced Hero Content -->
        <div class="hidden lg:flex flex-col justify-center space-y-8">
          <div class="max-w-lg">
            <!-- Animated badge -->
            <div class="hero-badge group cursor-default">
              <svg class="w-4 h-4 transition-transform duration-300 group-hover:rotate-12" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
              <span class="animate-text-gradient bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 bg-clip-text text-transparent">Secure & Reliable</span>
            </div>
            
            <!-- Enhanced typography -->
            <h1 class="hero-title mt-6 text-5xl xl:text-6xl font-bold leading-tight tracking-tight">
              <span class="block animate-slide-up">Access Your</span>
              <span class="block animate-slide-up bg-gradient-to-r from-slate-900 via-blue-900 to-indigo-900 bg-clip-text text-transparent" style="animation-delay: 0.2s;">Incident Management</span>
              <span class="block animate-slide-up text-blue-600" style="animation-delay: 0.4s;">Dashboard</span>
            </h1>
            
            <p class="hero-desc mt-6 text-xl text-slate-600 leading-relaxed animate-fade-in" style="animation-delay: 0.6s;">
              Monitor, manage, and resolve incidents efficiently with our comprehensive platform. 
              Get real-time alerts, track response times, and maintain system reliability.
            </p>
            
            <!-- Enhanced feature list with stagger animation -->
            <div class="mt-10 space-y-5">
              <div class="flex items-center space-x-4 animate-slide-in-left" style="animation-delay: 0.8s;">
                <div class="flex-shrink-0">
                  <div class="w-10 h-10 bg-gradient-to-r from-blue-500 to-blue-600 rounded-xl flex items-center justify-center shadow-lg">
                    <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z" />
                    </svg>
                  </div>
                </div>
                <span class="text-lg font-medium text-slate-700">Real-time incident monitoring</span>
              </div>
              
              <div class="flex items-center space-x-4 animate-slide-in-left" style="animation-delay: 1s;">
                <div class="flex-shrink-0">
                  <div class="w-10 h-10 bg-gradient-to-r from-indigo-500 to-indigo-600 rounded-xl flex items-center justify-center shadow-lg">
                    <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M3 3a1 1 0 000 2v8a2 2 0 002 2h2.586l-1.293 1.293a1 1 0 101.414 1.414L10 15.414l2.293 2.293a1 1 0 001.414-1.414L12.414 15H15a2 2 0 002-2V5a1 1 0 100-2H3zm11 4a1 1 0 10-2 0v4a1 1 0 102 0V7zm-3 1a1 1 0 10-2 0v3a1 1 0 102 0V8zM8 9a1 1 0 00-2 0v2a1 1 0 102 0V9z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
                <span class="text-lg font-medium text-slate-700">Automated alert aggregation</span>
              </div>
              
              <div class="flex items-center space-x-4 animate-slide-in-left" style="animation-delay: 1.2s;">
                <div class="flex-shrink-0">
                  <div class="w-10 h-10 bg-gradient-to-r from-purple-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg">
                    <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v3h8v-3z" />
                    </svg>
                  </div>
                </div>
                <span class="text-lg font-medium text-slate-700">Team collaboration tools</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Panel - Enhanced Login Form -->
        <div class="flex flex-col justify-center animate-slide-in-right">
          <!-- Mobile logo/branding with enhanced animation -->
          <div class="lg:hidden text-center mb-10 animate-bounce-in">
            <div class="inline-flex items-center justify-center w-18 h-18 rounded-3xl bg-gradient-to-tr from-blue-600 via-blue-700 to-indigo-700 shadow-xl mb-6 relative overflow-hidden group">
              <!-- Shimmer effect -->
              <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -skew-x-12 transform translate-x-full group-hover:translate-x-[-200%] transition-transform duration-1000"></div>
              <svg class="w-10 h-10 text-white relative z-10" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
            </div>
            <h1 class="text-3xl font-bold bg-gradient-to-r from-slate-900 to-blue-900 bg-clip-text text-transparent">Incident Management</h1>
            <p class="text-slate-600 text-base mt-2">Secure & Reliable Platform</p>
          </div>
          
          <div class="card-modern max-w-lg mx-auto w-full">
            <!-- Enhanced Header with gradient text -->
            <div class="text-center mb-10">
              <h2 class="text-4xl font-bold bg-gradient-to-r from-slate-900 via-blue-900 to-indigo-900 bg-clip-text text-transparent mb-3">Welcome Back</h2>
              <p class="text-lg text-slate-600">Sign in to your account to continue</p>
              <!-- Progress indicator -->
              <div class="mt-4 flex justify-center">
                <div class="flex space-x-2">
                  <div class="w-8 h-1 bg-blue-600 rounded-full"></div>
                  <div class="w-8 h-1 bg-slate-200 rounded-full"></div>
                  <div class="w-8 h-1 bg-slate-200 rounded-full"></div>
                </div>
              </div>
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
/* Modern Web Platform Animations */
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

/* Enhanced card with modern effects */
.card-modern {
  @apply bg-white/80 backdrop-blur-xl rounded-3xl shadow-2xl p-10 border border-white/30;
  box-shadow: 
    0 32px 64px -12px rgba(0, 0, 0, 0.08),
    0 25px 50px -12px rgba(59, 130, 246, 0.15),
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
  background: linear-gradient(90deg, transparent, rgba(59, 130, 246, 0.5), transparent);
}

/* Grid background pattern */
.bg-grid-pattern {
  background-image: 
    linear-gradient(rgba(59, 130, 246, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(59, 130, 246, 0.1) 1px, transparent 1px);
  background-size: 20px 20px;
}

/* Enhanced component styles */
.hero-badge {
  @apply inline-flex items-center bg-gradient-to-r from-blue-50 to-indigo-50 text-blue-700 px-4 py-2 rounded-2xl text-sm font-semibold border border-blue-200/50 shadow-sm;
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
  @apply bg-gradient-to-r from-blue-600 via-blue-700 to-indigo-700 text-white font-semibold rounded-2xl shadow-xl hover:shadow-2xl transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2;
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
  box-shadow: 0 25px 50px -12px rgba(59, 130, 246, 0.4);
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
    0 0 0 3px rgba(59, 130, 246, 0.12), 
    0 10px 25px -5px rgba(59, 130, 246, 0.25),
    0 0 0 1px rgba(59, 130, 246, 0.5);
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

/* Custom scrollbar for modern look */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(148, 163, 184, 0.1);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, rgba(59, 130, 246, 0.8), rgba(99, 102, 241, 0.8));
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, rgba(59, 130, 246, 1), rgba(99, 102, 241, 1));
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
  
  .hero-title {
    font-size: 2.5rem;
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