<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { incidentAPI } from '@/services/api'
import { formatDate, calculateDuration } from '@/utils/format'
import type { Incident } from '@/types/api'

const incidents = ref<Incident[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

const loadIncidents = async () => {
  try {
    loading.value = true
    error.value = null
    incidents.value = await incidentAPI.getIncidents()
    // Sort by creation date (newest first)
    incidents.value.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  } catch (err) {
    error.value = 'Error loading incidents'
    console.error('Error loading incidents:', err)
  } finally {
    loading.value = false
  }
}

const refreshIncidents = () => {
  loadIncidents()
}

onMounted(() => {
  loadIncidents()
})
</script>

<template>
  <div class="incidents">
    <div class="page-header">
      <h2>Incidents</h2>
      <div class="actions">
        <button @click="refreshIncidents" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div class="incidents-container">
      <div class="card">
        <div v-if="loading" class="loading">
          Loading incidents...
        </div>
        
        <table v-else-if="incidents.length > 0" class="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Title</th>
              <th>Severity</th>
              <th>Status</th>
              <th>Created</th>
              <th>Duration</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="incident in incidents" :key="incident.id">
              <td>{{ incident.id.substring(0, 8) }}</td>
              <td>{{ incident.title }}</td>
              <td>
                <span :class="`severity-badge severity-${incident.severity.toLowerCase()}`">
                  {{ incident.severity }}
                </span>
              </td>
              <td>
                <span :class="`status-badge status-${incident.status.toLowerCase().replace(' ', '-')}`">
                  {{ incident.status }}
                </span>
              </td>
              <td>{{ formatDate(incident.created_at) }}</td>
              <td>{{ calculateDuration(incident.created_at, incident.resolved_at) }}</td>
              <td>
                <button class="btn btn-primary">Details</button>
              </td>
            </tr>
          </tbody>
        </table>
        
        <div v-else class="no-data">
          No incidents found
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.incidents {
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

.incidents-container {
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