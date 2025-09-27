<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { alertAPI } from '@/services/api'
import { formatDate } from '@/utils/format'
import type { Alert } from '@/types/api'

const alerts = ref<Alert[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

const loadAlerts = async () => {
  try {
    loading.value = true
    error.value = null
    alerts.value = await alertAPI.getAlerts()
    // Sort by creation date (newest first)
    alerts.value.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  } catch (err) {
    error.value = 'Error loading alerts'
    console.error('Error loading alerts:', err)
  } finally {
    loading.value = false
  }
}

const refreshAlerts = () => {
  loadAlerts()
}

onMounted(() => {
  loadAlerts()
})
</script>

<template>
  <div class="alerts">
    <div class="page-header">
      <h2>Alerts</h2>
      <div class="actions">
        <button @click="refreshAlerts" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div class="alerts-container">
      <div class="card">
        <div v-if="loading" class="loading">
          Loading alerts...
        </div>
        
        <table v-else-if="alerts.length > 0" class="table">
          <thead>
            <tr>
              <th>Alert Name</th>
              <th>Status</th>
              <th>Started</th>
              <th>Incident</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="alert in alerts" :key="alert.id">
              <td>{{ alert.alert_name }}</td>
              <td>
                <span :class="`status-badge status-${alert.status.toLowerCase()}`">
                  {{ alert.status }}
                </span>
              </td>
              <td>{{ formatDate(alert.starts_at) }}</td>
              <td>
                <span v-if="alert.incident_id">{{ alert.incident_id.substring(0, 8) }}</span>
                <span v-else>-</span>
              </td>
              <td>
                <button class="btn btn-primary">Details</button>
              </td>
            </tr>
          </tbody>
        </table>
        
        <div v-else class="no-data">
          No alerts found
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.alerts {
  padding: 20px 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h2 {
  color: #2c3e50;
  margin: 0;
}

.alerts-container {
  margin-top: 1rem;
}

.error-message {
  background: #ffeaea;
  color: #e74c3c;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
  border: 1px solid #f8cecc;
}

.no-data {
  text-align: center;
  padding: 2rem;
  color: #666;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .table {
    font-size: 0.85rem;
  }
  
  .table th,
  .table td {
    padding: 8px;
  }
}
</style>