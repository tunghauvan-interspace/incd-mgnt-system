import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: Dashboard
    },
    {
      path: '/incidents',
      name: 'incidents',
      // route level code-splitting
      // this generates a separate chunk (Incidents.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/Incidents.vue')
    },
    {
      path: '/alerts',
      name: 'alerts',
      component: () => import('../views/Alerts.vue')
    }
  ]
})

export default router