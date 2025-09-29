<template>
  <div>
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
      <p class="text-gray-600">Overview of incidents and system health</p>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <StatCard
        title="Total Incidents"
        :value="stats.total"
        icon="🎯"
        color="blue"
      />
      <StatCard
        title="Open Incidents"
        :value="stats.open"
        icon="🔴"
        color="red"
      />
      <StatCard
        title="Critical Alerts"
        :value="stats.critical"
        icon="⚠️"
        color="orange"
      />
      <StatCard
        title="Resolved Today"
        :value="stats.resolved"
        icon="✅"
        color="green"
      />
    </div>

    <!-- Recent Incidents -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="card">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Recent Incidents</h2>
        <div class="space-y-3">
          <div v-for="incident in recentIncidents" :key="incident.id" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-sm">{{ incident.title }}</p>
              <p class="text-xs text-gray-500">{{ formatDate(incident.createdAt) }}</p>
            </div>
            <span :class="getSeverityClass(incident.severity)" class="px-2 py-1 text-xs rounded-full">
              {{ incident.severity }}
            </span>
          </div>
        </div>
      </div>

      <!-- Active Alerts -->
      <div class="card">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Active Alerts</h2>
        <div class="space-y-3">
          <div v-for="alert in activeAlerts" :key="alert.id" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-sm">{{ alert.summary }}</p>
              <p class="text-xs text-gray-500">{{ formatDate(alert.startsAt) }}</p>
            </div>
            <span :class="getSeverityClass(alert.severity)" class="px-2 py-1 text-xs rounded-full">
              {{ alert.severity }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useIncidents } from '@/composables/useIncidents'
import { useAlertsStore } from '@/stores/alerts'
import { formatDate } from '@/utils/formatters'
import StatCard from '@/components/ui/StatCard.vue'

const { incidents, fetchIncidents } = useIncidents()
const alertsStore = useAlertsStore()

const stats = ref({
  total: 0,
  open: 0,
  critical: 0,
  resolved: 0
})

const recentIncidents = ref([])
const activeAlerts = ref([])

const getSeverityClass = (severity: string) => {
  const classes = {
    critical: 'status-critical',
    high: 'status-high',
    medium: 'status-medium',
    low: 'status-low'
  }
  return classes[severity] || 'bg-gray-100 text-gray-800'
}

onMounted(async () => {
  await fetchIncidents()
  await alertsStore.fetchAlerts()
  
  // Update stats
  stats.value.total = incidents.value.length
  stats.value.open = incidents.value.filter(i => i.status === 'open').length
  stats.value.critical = incidents.value.filter(i => i.severity === 'critical').length
  stats.value.resolved = incidents.value.filter(i => 
    i.status === 'resolved' && 
    new Date(i.resolvedAt || '').toDateString() === new Date().toDateString()
  ).length

  // Get recent incidents
  recentIncidents.value = incidents.value
    .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
    .slice(0, 5)

  // Get active alerts
  activeAlerts.value = alertsStore.firingAlerts.slice(0, 5)
})
</script>