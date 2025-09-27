<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { formatDate, calculateDuration } from '@/utils/format'
import Modal from '@/components/Modal.vue'
import Button from '@/components/Button.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import SeverityBadge from '@/components/SeverityBadge.vue'
import DataTable from '@/components/DataTable.vue'
import { useIncidents } from '@/composables/useIncidents'
import type { Incident } from '@/types/api'
import type { TableColumn } from '@/types/components'

const selectedIncident = ref<Incident | null>(null)
const showModal = ref(false)

const { incidents, loading, error, loadIncidents, acknowledgeIncident, resolveIncident } = useIncidents()

// Table configuration
const columns: TableColumn<Incident>[] = [
  { key: 'id', label: 'ID', sortable: true, width: '120px' },
  { key: 'title', label: 'Title', sortable: true },
  { key: 'severity', label: 'Severity', sortable: true, width: '120px', align: 'center' },
  { key: 'status', label: 'Status', sortable: true, width: '120px', align: 'center' },
  { key: 'created_at', label: 'Created', sortable: true, width: '180px' },
  { key: 'duration', label: 'Duration', width: '120px' },
  { key: 'actions', label: 'Actions', width: '200px', align: 'center' }
]

const showIncidentDetails = (incident: Incident) => {
  selectedIncident.value = incident
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  selectedIncident.value = null
}

const handleAcknowledge = async (incidentId: string) => {
  try {
    await acknowledgeIncident(incidentId, 'current_user')
  } catch (err) {
    console.error('Error acknowledging incident:', err)
    alert('Failed to acknowledge incident')
  }
}

const handleResolve = async (incidentId: string) => {
  try {
    await resolveIncident(incidentId)
    if (selectedIncident.value?.id === incidentId) {
      closeModal() // Close modal if this incident was being viewed
    }
  } catch (err) {
    console.error('Error resolving incident:', err)
    alert('Failed to resolve incident')
  }
}

const handleRowClick = (incident: Incident) => {
  showIncidentDetails(incident)
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
        <Button 
          @click="loadIncidents" 
          :loading="loading"
          variant="primary"
        >
          {{ loading ? 'Loading...' : 'Refresh' }}
        </Button>
      </div>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div class="incidents-container">
      <div class="card">
        <DataTable
          :columns="columns"
          :data="incidents"
          :loading="loading"
          empty-message="No incidents found"
          hoverable
          @row-click="handleRowClick"
        >
          <!-- Custom ID column -->
          <template #cell-id="{ value }">
            <code class="incident-id">{{ value.substring(0, 8) }}</code>
          </template>
          
          <!-- Custom Severity column -->
          <template #cell-severity="{ value }">
            <SeverityBadge :severity="value" size="sm" />
          </template>
          
          <!-- Custom Status column -->
          <template #cell-status="{ value }">
            <StatusBadge :status="value" size="sm" />
          </template>
          
          <!-- Custom Created column -->
          <template #cell-created_at="{ value }">
            {{ formatDate(value) }}
          </template>
          
          <!-- Custom Duration column -->
          <template #cell-duration="{ row }">
            {{ calculateDuration(row.created_at, row.resolved_at) }}
          </template>
          
          <!-- Custom Actions column -->
          <template #cell-actions="{ row }">
            <div class="actions-group">
              <Button 
                size="sm" 
                variant="primary"
                @click.stop="showIncidentDetails(row)"
              >
                Details
              </Button>
              <Button
                v-if="row.status === 'open'"
                size="sm"
                variant="warning"
                @click.stop="handleAcknowledge(row.id)"
              >
                Acknowledge
              </Button>
              <Button
                v-if="row.status !== 'resolved'"
                size="sm"
                variant="success"
                @click.stop="handleResolve(row.id)"
              >
                Resolve
              </Button>
            </div>
          </template>
        </DataTable>
      </div>
    </div>

    <!-- Incident Details Modal -->
    <Modal
      :show="showModal"
      :title="`Incident Details - ${selectedIncident?.id.substring(0, 8) || ''}`"
      @close="closeModal"
    >
      <div v-if="selectedIncident" class="incident-details">
        <div class="detail-row">
          <strong>ID:</strong> 
          <code class="incident-id">{{ selectedIncident.id }}</code>
        </div>
        <div class="detail-row">
          <strong>Title:</strong> 
          {{ selectedIncident.title }}
        </div>
        <div class="detail-row">
          <strong>Description:</strong>
          <p class="description">
            {{ selectedIncident.description || 'No description available' }}
          </p>
        </div>
        <div class="detail-row">
          <strong>Severity:</strong>
          <SeverityBadge 
            :severity="selectedIncident.severity" 
            size="md" 
            show-icon 
          />
        </div>
        <div class="detail-row">
          <strong>Status:</strong>
          <StatusBadge 
            :status="selectedIncident.status" 
            size="md" 
          />
        </div>
        <div class="detail-row">
          <strong>Created:</strong> 
          {{ formatDate(selectedIncident.created_at) }}
        </div>
        <div class="detail-row" v-if="selectedIncident.acknowledged_at">
          <strong>Acknowledged:</strong> 
          {{ formatDate(selectedIncident.acknowledged_at) }}
        </div>
        <div class="detail-row" v-if="selectedIncident.resolved_at">
          <strong>Resolved:</strong> 
          {{ formatDate(selectedIncident.resolved_at) }}
        </div>
        <div class="detail-row" v-if="selectedIncident.assignee_id">
          <strong>Assignee:</strong> 
          {{ selectedIncident.assignee_id }}
        </div>
        <div
          class="detail-row"
          v-if="selectedIncident.labels && Object.keys(selectedIncident.labels).length > 0"
        >
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
        <Button
          v-if="selectedIncident?.status === 'open'"
          variant="warning"
          @click="selectedIncident && handleAcknowledge(selectedIncident.id)"
        >
          Acknowledge
        </Button>
        <Button
          v-if="selectedIncident?.status !== 'resolved'"
          variant="success"
          @click="selectedIncident && handleResolve(selectedIncident.id)"
        >
          Resolve
        </Button>
        <Button variant="secondary" @click="closeModal">
          Close
        </Button>
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
