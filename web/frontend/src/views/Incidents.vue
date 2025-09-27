<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { incidentAPI } from '@/services/api'
import { formatDate, calculateDuration } from '@/utils/format'
import Modal from '@/components/Modal.vue'
import type { Incident } from '@/types/api'

const incidents = ref<Incident[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const selectedIncident = ref<Incident | null>(null)
const showModal = ref(false)

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

const showIncidentDetails = (incident: Incident) => {
  selectedIncident.value = incident
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  selectedIncident.value = null
}

const acknowledgeIncident = async (incidentId: string) => {
  try {
    await incidentAPI.acknowledgeIncident(incidentId, { assignee_id: 'current_user' })
    await loadIncidents() // Refresh the list
  } catch (err) {
    console.error('Error acknowledging incident:', err)
    alert('Failed to acknowledge incident')
  }
}

const resolveIncident = async (incidentId: string) => {
  try {
    await incidentAPI.resolveIncident(incidentId)
    await loadIncidents() // Refresh the list
    if (selectedIncident.value?.id === incidentId) {
      closeModal() // Close modal if this incident was being viewed
    }
  } catch (err) {
    console.error('Error resolving incident:', err)
    alert('Failed to resolve incident')
  }
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
                <div class="actions-group">
                  <button class="btn btn-primary btn-sm" @click="showIncidentDetails(incident)">
                    Details
                  </button>
                  <button 
                    v-if="incident.status === 'open'"
                    class="btn btn-warning btn-sm"
                    @click="acknowledgeIncident(incident.id)"
                  >
                    Acknowledge
                  </button>
                  <button 
                    v-if="incident.status !== 'resolved'"
                    class="btn btn-success btn-sm"
                    @click="resolveIncident(incident.id)"
                  >
                    Resolve
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
        
        <div v-else class="no-data">
          No incidents found
        </div>
      </div>
    </div>

    <!-- Incident Details Modal -->
    <Modal :show="showModal" :title="`Incident Details - ${selectedIncident?.id.substring(0, 8) || ''}`" @close="closeModal">
      <div v-if="selectedIncident" class="incident-details">
        <div class="detail-row">
          <strong>ID:</strong> {{ selectedIncident.id }}
        </div>
        <div class="detail-row">
          <strong>Title:</strong> {{ selectedIncident.title }}
        </div>
        <div class="detail-row">
          <strong>Description:</strong> 
          <p class="description">{{ selectedIncident.description || 'No description available' }}</p>
        </div>
        <div class="detail-row">
          <strong>Severity:</strong>
          <span :class="`severity-badge severity-${selectedIncident.severity.toLowerCase()}`">
            {{ selectedIncident.severity }}
          </span>
        </div>
        <div class="detail-row">
          <strong>Status:</strong>
          <span :class="`status-badge status-${selectedIncident.status.toLowerCase().replace(' ', '-')}`">
            {{ selectedIncident.status }}
          </span>
        </div>
        <div class="detail-row">
          <strong>Created:</strong> {{ formatDate(selectedIncident.created_at) }}
        </div>
        <div class="detail-row" v-if="selectedIncident.acknowledged_at">
          <strong>Acknowledged:</strong> {{ formatDate(selectedIncident.acknowledged_at) }}
        </div>
        <div class="detail-row" v-if="selectedIncident.resolved_at">
          <strong>Resolved:</strong> {{ formatDate(selectedIncident.resolved_at) }}
        </div>
        <div class="detail-row" v-if="selectedIncident.assignee_id">
          <strong>Assignee:</strong> {{ selectedIncident.assignee_id }}
        </div>
        <div class="detail-row" v-if="selectedIncident.labels && Object.keys(selectedIncident.labels).length > 0">
          <strong>Labels:</strong>
          <div class="labels-container">
            <span 
              v-for="[key, value] in Object.entries(selectedIncident.labels)" 
              :key="key" 
              class="label-tag"
            >
              <strong>{{ key }}:</strong> {{ value }}
            </span>
          </div>
        </div>
      </div>

      <template #footer>
        <button 
          v-if="selectedIncident?.status === 'open'"
          class="btn btn-warning"
          @click="selectedIncident && acknowledgeIncident(selectedIncident.id)"
        >
          Acknowledge
        </button>
        <button 
          v-if="selectedIncident?.status !== 'resolved'"
          class="btn btn-success"
          @click="selectedIncident && resolveIncident(selectedIncident.id)"
        >
          Resolve
        </button>
        <button class="btn btn-secondary" @click="closeModal">
          Close
        </button>
      </template>
    </Modal>
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

.actions-group {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 0.85rem;
}

.incident-details {
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
  min-width: 100px;
  flex-shrink: 0;
}

.description {
  margin: 0.5rem 0 0 0;
  background: #f8f9fa;
  padding: 0.75rem;
  border-radius: 4px;
  border-left: 3px solid #3498db;
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
  
  .actions-group {
    flex-direction: column;
  }
  
  .btn-sm {
    padding: 6px 10px;
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