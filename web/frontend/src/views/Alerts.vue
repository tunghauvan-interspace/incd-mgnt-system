<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { alertAPI } from '@/services/api'
import { formatDate } from '@/utils/format'
import Modal from '@/components/Modal.vue'
import type { Alert } from '@/types/api'

const alerts = ref<Alert[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const selectedAlert = ref<Alert | null>(null)
const showModal = ref(false)

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

const showAlertDetails = (alert: Alert) => {
  selectedAlert.value = alert
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  selectedAlert.value = null
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
                <button class="btn btn-primary" @click="showAlertDetails(alert)">
                  Details
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        
        <div v-else class="no-data">
          No alerts found
        </div>
      </div>
    </div>

    <!-- Alert Details Modal -->
    <Modal :show="showModal" :title="`Alert Details - ${selectedAlert?.alert_name || ''}`" @close="closeModal">
      <div v-if="selectedAlert" class="alert-details">
        <div class="detail-row">
          <strong>ID:</strong> {{ selectedAlert.id }}
        </div>
        <div class="detail-row">
          <strong>Alert Name:</strong> {{ selectedAlert.alert_name }}
        </div>
        <div class="detail-row">
          <strong>Status:</strong>
          <span :class="`status-badge status-${selectedAlert.status.toLowerCase()}`">
            {{ selectedAlert.status }}
          </span>
        </div>
        <div class="detail-row">
          <strong>Started:</strong> {{ formatDate(selectedAlert.starts_at) }}
        </div>
        <div class="detail-row" v-if="selectedAlert.ends_at">
          <strong>Ended:</strong> {{ formatDate(selectedAlert.ends_at) }}
        </div>
        <div class="detail-row" v-if="selectedAlert.generator_url">
          <strong>Generator URL:</strong>
          <a :href="selectedAlert.generator_url" target="_blank" rel="noopener noreferrer">
            {{ selectedAlert.generator_url }}
          </a>
        </div>
        <div class="detail-row" v-if="selectedAlert.incident_id">
          <strong>Incident ID:</strong> {{ selectedAlert.incident_id }}
        </div>
        <div class="detail-row" v-if="selectedAlert.labels && Object.keys(selectedAlert.labels).length > 0">
          <strong>Labels:</strong>
          <div class="labels-container">
            <span 
              v-for="[key, value] in Object.entries(selectedAlert.labels)" 
              :key="key" 
              class="label-tag"
            >
              <strong>{{ key }}:</strong> {{ value }}
            </span>
          </div>
        </div>
        <div class="detail-row" v-if="selectedAlert.annotations && Object.keys(selectedAlert.annotations).length > 0">
          <strong>Annotations:</strong>
          <div class="annotations-container">
            <div 
              v-for="[key, value] in Object.entries(selectedAlert.annotations)" 
              :key="key" 
              class="annotation-item"
            >
              <strong>{{ key }}:</strong>
              <div class="annotation-value">{{ value }}</div>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <button class="btn btn-secondary" @click="closeModal">
          Close
        </button>
      </template>
    </Modal>
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

.alert-details {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.detail-row {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
}

.detail-row strong {
  min-width: 120px;
  flex-shrink: 0;
}

.detail-row a {
  color: #3498db;
  text-decoration: none;
  word-break: break-all;
}

.detail-row a:hover {
  text-decoration: underline;
}

.labels-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.label-tag {
  background: #e9ecef;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.85rem;
}

.annotations-container {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.annotation-item {
  background: #f8f9fa;
  padding: 0.5rem;
  border-radius: 4px;
  border-left: 3px solid #17a2b8;
}

.annotation-value {
  margin-top: 0.25rem;
  font-family: 'Courier New', monospace;
  font-size: 0.9rem;
  white-space: pre-wrap;
  word-break: break-word;
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
  
  .detail-row {
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .detail-row strong {
    min-width: auto;
  }
}
</style>