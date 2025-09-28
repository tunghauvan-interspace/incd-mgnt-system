<template>
  <div
    class="min-h-screen bg-gradient-to-br from-white via-gray-100 to-gray-200 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8"
  >
    <div class="max-w-5xl w-full grid grid-cols-1 gap-10 lg:grid-cols-[1.1fr_0.9fr] items-center">
      <section
        class="hidden h-full flex-col justify-center space-y-6 rounded-3xl bg-white/80 p-10 shadow-lg backdrop-blur-sm lg:flex"
      >
        <span
          class="inline-flex w-fit items-center gap-2 rounded-full bg-blue-50 px-4 py-2 text-sm font-medium text-blue-700"
        >
          <span class="h-2 w-2 rounded-full bg-blue-500"></span>
          Incident Command Center
        </span>
        <h1 class="text-4xl font-semibold leading-tight text-gray-900">
          Stay ahead of every alert with a dependable response workflow.
        </h1>
        <p class="text-base leading-relaxed text-gray-600">
          Triage incidents, collaborate with your team, and close loops faster with real-time
          updates and a command dashboard designed for on-call engineers.
        </p>
        <div class="space-y-3 text-sm text-gray-700">
          <div class="flex items-center gap-3">
            <span
              class="grid h-8 w-8 place-items-center rounded-full bg-blue-100 text-blue-600 text-lg"
              >✓</span
            >
            <p>Unified view of alerts, incidents, and escalations.</p>
          </div>
          <div class="flex items-center gap-3">
            <span
              class="grid h-8 w-8 place-items-center rounded-full bg-emerald-100 text-emerald-600 text-lg"
              >✓</span
            >
            <p>Actionable analytics with severity-based prioritization.</p>
          </div>
          <div class="flex items-center gap-3">
            <span
              class="grid h-8 w-8 place-items-center rounded-full bg-amber-100 text-amber-600 text-lg"
              >✓</span
            >
            <p>Seamless handoffs across shifts and notification channels.</p>
          </div>
        </div>
      </section>

      <div class="w-full rounded-3xl bg-white p-8 shadow-2xl ring-1 ring-black/5 sm:p-10">
        <div class="space-y-3 text-center">
          <span
            class="inline-flex items-center justify-center rounded-full bg-blue-50 px-3 py-1 text-xs font-medium uppercase tracking-wide text-blue-600"
          >
            Welcome back
          </span>
          <h2 class="text-2xl font-semibold text-gray-900 sm:text-3xl">Sign in to continue</h2>
          <p class="text-sm text-gray-600">
            Access your incident workspace and keep everything running smoothly.
          </p>
        </div>

        <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
          <div class="space-y-5">
            <div class="flex flex-col space-y-2">
              <label for="username" class="text-sm font-medium text-gray-700">Username</label>
              <input
                id="username"
                name="username"
                type="text"
                autocomplete="username"
                required
                class="w-full rounded-xl border border-gray-200 bg-white px-4 py-3 text-gray-900 shadow-inner focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
                placeholder="your username"
                v-model="form.username"
              />
            </div>
            <div class="flex flex-col space-y-2">
              <label for="password" class="text-sm font-medium text-gray-700">Password</label>
              <input
                id="password"
                name="password"
                type="password"
                autocomplete="current-password"
                required
                class="w-full rounded-xl border border-gray-200 bg-white px-4 py-3 text-gray-900 shadow-inner focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
                placeholder="••••••••"
                v-model="form.password"
              />
            </div>
          </div>

          <div
            v-if="error"
            class="rounded-2xl border border-red-200 bg-red-50/80 p-4 text-sm text-red-700"
          >
            {{ error }}
          </div>

          <div class="space-y-4">
            <button
              type="submit"
              :disabled="isLoading"
              class="relative flex w-full items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-blue-600 via-blue-500 to-blue-600 px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-blue-500/20 transition duration-200 hover:from-blue-500 hover:via-blue-500 hover:to-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-200 disabled:cursor-not-allowed disabled:opacity-60"
            >
              <svg
                v-if="isLoading"
                class="h-5 w-5 animate-spin"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                />
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                />
              </svg>
              <span>{{ isLoading ? 'Signing you in…' : 'Sign in' }}</span>
            </button>

            <p class="text-center text-sm text-gray-600">
              Need an account?
              <router-link to="/register" class="font-medium text-blue-600 hover:text-blue-500">
                Create one now
              </router-link>
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  username: '',
  password: ''
})

const error = ref<string | null>(null)
const isLoading = ref(false)

const handleLogin = async () => {
  if (isLoading.value) return

  isLoading.value = true
  error.value = null

  try {
    await authStore.login({
      username: form.value.username,
      password: form.value.password
    })

    // Redirect to dashboard on success
    router.push('/')
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Login failed. Please try again.'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* Fallback styles in case Tailwind isn't compiled in dev */
.min-h-screen {
  min-height: 100vh;
}
body, #app {
  background: linear-gradient(180deg, #ffffff 0%, #f3f4f6 100%);
}
.max-w-5xl {
  max-width: 72rem;
  margin: 0 auto;
}
.rounded-3xl {
  border-radius: 1.5rem;
}
.bg-white {
  background: #fff;
}
.shadow-2xl {
  box-shadow: 0 25px 50px -12px rgba(0,0,0,0.25);
}
.ring-1 {
  box-shadow: 0 0 0 1px rgba(0,0,0,0.04) inset;
}
.p-8 { padding: 2rem; }
.sm\:p-10 { padding: 2.5rem; }
.text-center { text-align: center; }
.text-gray-900 { color: #0f172a; }
.text-gray-600 { color: #475569; }
.text-gray-700 { color: #334155; }
.w-full { width: 100%; }
.rounded-xl { border-radius: 0.75rem; }
.border { border: 1px solid #e6e6e6; }
.border-gray-200 { border-color: #e6e6e6; }
.px-4 { padding-left: 1rem; padding-right: 1rem; }
.py-3 { padding-top: .75rem; padding-bottom: .75rem; }
.shadow-inner { box-shadow: inset 0 1px 2px rgba(0,0,0,0.03); }
input[type="text"], input[type="password"] {
  display: block;
  width: 100%;
  padding: .75rem 1rem;
  border-radius: .5rem;
  border: 1px solid #d1d5db;
  background: #fff;
  color: #0f172a;
  font-size: 14px;
}
label { display:block; margin-bottom: .35rem; font-weight: 600; }
.space-y-4 > * + * { margin-top: 1rem; }
.relative { position: relative; }
.flex { display: flex; }
.items-center { align-items: center; }
.justify-center { justify-content: center; }
.gap-2 { gap: .5rem; }
.py-3 { padding-top: .75rem; padding-bottom: .75rem; }
.px-6 { padding-left: 1.5rem; padding-right: 1.5rem; }
.font-semibold { font-weight: 600; }
.text-white { color: #fff; }
.btn-fallback {
  display:inline-flex; align-items:center; justify-content:center;
  width:100%; background: linear-gradient(90deg,#2563eb,#1d4ed8); color:#fff;
  padding:.75rem 1.25rem; border-radius:.75rem; border:none; cursor:pointer; font-weight:600;
  box-shadow: 0 8px 20px rgba(37,99,235,0.18);
}
.btn-fallback:disabled { opacity: .6; cursor: not-allowed; }
.card-hero { padding: 2.5rem; border-radius: 1rem; }
.hidden { display:none; }
@media (min-width: 1024px) {
  .hidden.lg\:flex { display:flex; }
}
</style>
