<script setup lang="ts">
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <nav class="navbar">
    <div class="nav-container">
      <div class="nav-logo">
        <h1>Incident Management System</h1>
      </div>
      <div class="nav-content">
        <ul v-if="authStore.isAuthenticated" class="nav-menu">
          <li class="nav-item">
            <RouterLink to="/" class="nav-link" :class="{ active: route.name === 'dashboard' }">
              Dashboard
            </RouterLink>
          </li>
          <li class="nav-item">
            <RouterLink
              to="/incidents"
              class="nav-link"
              :class="{ active: route.name === 'incidents' }"
            >
              Incidents
            </RouterLink>
          </li>
          <li class="nav-item">
            <RouterLink to="/alerts" class="nav-link" :class="{ active: route.name === 'alerts' }">
              Alerts
            </RouterLink>
          </li>
        </ul>
        <div v-if="authStore.isAuthenticated" class="nav-user">
          <span class="user-info">
            Welcome, {{ authStore.user?.full_name || authStore.user?.username }}
          </span>
          <button @click="handleLogout" class="logout-btn">Logout</button>
        </div>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.navbar {
  background: var(--color-gray-800);
  color: var(--color-text-white);
  padding: 1rem 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.nav-content {
  display: flex;
  align-items: center;
  gap: 2rem;
}

.nav-logo h1 {
  font-size: 1.5rem;
  font-weight: 600;
}

.nav-menu {
  display: flex;
  list-style: none;
  margin: 0;
}

.nav-item {
  margin-left: 2rem;
}

.nav-link {
  color: var(--color-text-white);
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.nav-link:hover,
.nav-link.active {
  background-color: var(--color-gray-700);
}

.nav-user {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-info {
  color: var(--color-gray-300);
  font-size: 0.9rem;
}

.logout-btn {
  background: var(--color-danger);
  color: var(--color-text-white);
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.3s;
}

.logout-btn:hover {
  background: var(--color-danger-hover);
}

@media (max-width: 768px) {
  .nav-container {
    flex-direction: column;
    gap: 1rem;
  }

  .nav-content {
    flex-direction: column;
    gap: 1rem;
    width: 100%;
  }

  .nav-menu {
    gap: 1rem;
    justify-content: center;
  }

  .nav-item {
    margin-left: 0;
  }

  .nav-logo h1 {
    font-size: 1.2rem;
  }

  .nav-user {
    flex-direction: column;
    gap: 0.5rem;
  }
}
</style>
