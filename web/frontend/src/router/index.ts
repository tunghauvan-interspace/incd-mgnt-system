import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// Import layouts
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'

// Import pages
import Dashboard from '@/pages/Dashboard.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: DefaultLayout,
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: Dashboard,
        meta: { requiresAuth: true }
      },
      {
        path: '/incidents',
        name: 'Incidents',
        component: () => import('@/pages/Incidents/index.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/alerts',
        name: 'Alerts',
        component: () => import('@/pages/Alerts/index.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/users',
        name: 'Users',
        component: () => import('@/pages/Users/index.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/settings',
        name: 'Settings',
        component: () => import('@/pages/Settings/index.vue'),
        meta: { requiresAuth: true }
      }
    ]
  },
  {
    path: '/auth',
    component: AuthLayout,
    children: [
      {
        path: 'login',
        name: 'Login',
        component: () => import('@/pages/Auth/Login.vue'),
        meta: { guest: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Import auth store
import { useAuthStore } from '@/stores/auth'

// Navigation guard for authentication
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/auth/login')
  } else if (to.meta.guest && authStore.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router