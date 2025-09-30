<template>
  <div>
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Alerts</h1>
      <p class="text-gray-600">Monitor and manage active alerts</p>
    </div>

    <!-- Alert Summary -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <StatCard
        title="Firing Alerts"
        :value="firingAlerts.length"
        icon="🔥"
        color="red"
      />
      <StatCard
        title="Critical Alerts"
        :value="criticalAlerts.length"
        icon="⚠️"
        color="orange"
      />
      <StatCard
        title="Resolved Today"
        :value="resolvedToday"
        icon="✅"
        color="green"
      />
    </div>

    <!-- Alerts List -->
    <div class="card">
      <div v-if="isLoading" class="text-center py-8">
        <p class="text-gray-500">Loading alerts...</p>
      </div>
      
      <div v-else-if="alerts.length === 0" class="text-center py-8">
        <p class="text-gray-500">No alerts found</p>
      </div>
      
      <div v-else class="space-y-4">
        <div
          v-for="alert in alerts"
          :key="alert.id"
          class="p-4 border border-gray-200 rounded-lg"
          :class="getAlertBorderClass(alert.status)"
        >
          <div class="flex justify-between items-start">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <h3 class="text-lg font-medium text-gray-900">{{ alert.summary }}</h3>
                <span :class="getSeverityClass(alert.severity)" class="px-2 py-1 text-xs rounded-full">
                  {{ alert.severity }}
                </span>
                <span :class="getStatusClass(alert.status)" class="px-2 py-1 text-xs rounded-full">
                  {{ alert.status }}
                </span>
              </div>
              <p class="text-gray-600 mb-2">{{ alert.description }}</p>
              <div class="flex items-center text-sm text-gray-500 space-x-4">
                <span>Started: {{ formatDate(alert.startsAt) }}</span>
                <span v-if="alert.source">Source: {{ alert.source }}</span>
                <span v-if="alert.incidentId">Incident: {{ alert.incidentId }}</span>
              </div>
            </div>
            <div class="flex space-x-2" v-if="alert.status === 'firing'">
              <button 
                @click="acknowledgeAlert(alert.id)"
                class="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded hover:bg-blue-200"
              >
                Acknowledge
              </button>
              <button 
                @click="silenceAlert(alert.id)"
                class="px-3 py-1 text-sm bg-yellow-100 text-yellow-700 rounded hover:bg-yellow-200"
              >
                Silence
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAlertsStore } from '@/stores/alerts'
import { formatDate } from '@/utils/formatters'
import StatCard from '@/components/ui/StatCard.vue'

const alertsStore = useAlertsStore()

const isLoading = ref(false)

const alerts = computed(() => alertsStore.alerts)
const firingAlerts = computed(() => alertsStore.firingAlerts)
const criticalAlerts = computed(() => alertsStore.criticalAlerts)

const resolvedToday = computed(() => {
  const today = new Date().toDateString()
  return alertsStore.resolvedAlerts.filter(alert => 
    alert.endsAt && new Date(alert.endsAt).toDateString() === today
  ).length
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
    firing: 'bg-red-100 text-red-800',
    resolved: 'bg-green-100 text-green-800',
    silenced: 'bg-gray-100 text-gray-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getAlertBorderClass = (status: string) => {
  const classes = {
    firing: 'border-red-200 bg-red-50',
    resolved: 'border-green-200',
    silenced: 'border-gray-200'
  }
  return classes[status] || 'border-gray-200'
}

const acknowledgeAlert = async (alertId: string) => {
  try {
    await alertsStore.acknowledgeAlert(alertId)
  } catch (error) {
    console.error('Failed to acknowledge alert:', error)
  }
}

const silenceAlert = async (alertId: string) => {
  try {
    // Default silence for 1 hour (60 minutes)
    await alertsStore.silenceAlert(alertId, 60)
  } catch (error) {
    console.error('Failed to silence alert:', error)
  }
}

onMounted(async () => {
  isLoading.value = true
  try {
    await alertsStore.fetchAlerts()
  } finally {
    isLoading.value = false
  }
})
</script>