<template>
  <div>
    <div class="mb-6 flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Incidents</h1>
        <p class="text-gray-600">Manage and track incident responses</p>
      </div>
      <button class="btn btn-primary">
        Create Incident
      </button>
    </div>

    <!-- Filters -->
    <div class="card mb-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
          <select v-model="filters.status" class="w-full p-2 border border-gray-300 rounded-md">
            <option value="">All Statuses</option>
            <option value="open">Open</option>
            <option value="acknowledged">Acknowledged</option>
            <option value="investigating">Investigating</option>
            <option value="resolved">Resolved</option>
            <option value="closed">Closed</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Severity</label>
          <select v-model="filters.severity" class="w-full p-2 border border-gray-300 rounded-md">
            <option value="">All Severities</option>
            <option value="critical">Critical</option>
            <option value="high">High</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
          <input
            v-model="filters.search"
            type="text"
            placeholder="Search incidents..."
            class="w-full p-2 border border-gray-300 rounded-md"
          />
        </div>
        <div class="flex items-end">
          <button @click="clearFilters" class="btn text-gray-600 border border-gray-300">
            Clear Filters
          </button>
        </div>
      </div>
    </div>

    <!-- Incidents List -->
    <div class="card">
      <div v-if="isLoading" class="text-center py-8">
        <p class="text-gray-500">Loading incidents...</p>
      </div>
      
      <div v-else-if="filteredIncidents.length === 0" class="text-center py-8">
        <p class="text-gray-500">No incidents found</p>
      </div>
      
      <div v-else class="space-y-4">
        <div
          v-for="incident in filteredIncidents"
          :key="incident.id"
          class="p-4 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer"
          @click="$router.push(`/incidents/${incident.id}`)"
        >
          <div class="flex justify-between items-start">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <h3 class="text-lg font-medium text-gray-900">{{ incident.title }}</h3>
                <span :class="getSeverityClass(incident.severity)" class="px-2 py-1 text-xs rounded-full">
                  {{ incident.severity }}
                </span>
                <span :class="getStatusClass(incident.status)" class="px-2 py-1 text-xs rounded-full">
                  {{ incident.status }}
                </span>
              </div>
              <p class="text-gray-600 mb-2">{{ incident.description }}</p>
              <div class="flex items-center text-sm text-gray-500 space-x-4">
                <span>Created: {{ formatDate(incident.createdAt) }}</span>
                <span v-if="incident.assignee">Assigned to: {{ incident.assignee.name }}</span>
                <span>{{ incident.alertCount }} alerts</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useIncidents } from '@/composables/useIncidents'
import { formatDate } from '@/utils/formatters'

const { incidents, fetchIncidents, isLoading, setFilters, clearFilters: clearStoreFilters } = useIncidents()

const filters = ref({
  status: '',
  severity: '',
  search: ''
})

const filteredIncidents = computed(() => {
  return incidents.value.filter(incident => {
    if (filters.value.status && incident.status !== filters.value.status) return false
    if (filters.value.severity && incident.severity !== filters.value.severity) return false
    if (filters.value.search && !incident.title.toLowerCase().includes(filters.value.search.toLowerCase())) return false
    return true
  })
})

const getSeverityClass = (severity: string) => {
  const classes = {
    critical: 'status-critical',
    high: 'status-high',
    medium: 'status-medium',
    low: 'status-low'
  }
  return classes[severity] || 'bg-gray-100 text-gray-800'
}

const getStatusClass = (status: string) => {
  const classes = {
    open: 'bg-red-100 text-red-800',
    acknowledged: 'bg-yellow-100 text-yellow-800',
    investigating: 'bg-blue-100 text-blue-800',
    resolved: 'bg-green-100 text-green-800',
    closed: 'bg-gray-100 text-gray-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const clearFilters = () => {
  filters.value = {
    status: '',
    severity: '',
    search: ''
  }
  clearStoreFilters()
}

// Watch filters and update store
watch(filters, (newFilters) => {
  setFilters(newFilters)
}, { deep: true })

onMounted(() => {
  fetchIncidents()
})
</script>