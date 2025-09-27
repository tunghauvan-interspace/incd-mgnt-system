<script setup lang="ts">
import { RouterLink, useRoute } from 'vue-router'

const route = useRoute()
</script>

<template>
  <nav class="navbar">
    <div class="nav-container">
      <div class="nav-logo">
        <h1>Incident Management System</h1>
      </div>
      <ul class="nav-menu">
        <li class="nav-item">
          <RouterLink 
            to="/" 
            class="nav-link" 
            :class="{ active: route.name === 'dashboard' }"
          >
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
          <RouterLink 
            to="/alerts" 
            class="nav-link" 
            :class="{ active: route.name === 'alerts' }"
          >
            Alerts
          </RouterLink>
        </li>
      </ul>
    </div>
  </nav>
</template>

<style scoped>
.navbar {
  background: #2c3e50;
  color: white;
  padding: 1rem 0;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.nav-logo h1 {
  font-size: 1.5rem;
  font-weight: 600;
}

.nav-menu {
  display: flex;
  list-style: none;
}

.nav-item {
  margin-left: 2rem;
}

.nav-link {
  color: white;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.nav-link:hover,
.nav-link.active {
  background-color: #34495e;
}

@media (max-width: 768px) {
  .nav-container {
    flex-direction: column;
    gap: 1rem;
  }
  
  .nav-menu {
    gap: 1rem;
  }
  
  .nav-item {
    margin-left: 0;
  }
  
  .nav-logo h1 {
    font-size: 1.2rem;
  }
}
</style>