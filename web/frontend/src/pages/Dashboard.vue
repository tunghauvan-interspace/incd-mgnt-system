<template>
  <div>
    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1 class="page-title">Dashboard</h1>
        <p class="page-description">Welcome back! Here's what's happening with your incidents and alerts.</p>
      </div>
      <div class="flex space-x-3">
        <button class="btn btn-secondary btn-sm">
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4-4m0 0l-4 4m4-4v12" />
          </svg>
          Export Report
        </button>
        <button class="btn btn-primary btn-sm">
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Create Incident
        </button>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="stats-grid">
      <div 
        v-for="stat in statsData" 
        :key="stat.title"
        class="stat-card"
        :class="`stat-card-${stat.color}`"
        @click="stat.onClick"
      >
        <div class="stat-card-header">
          <div 
            class="stat-card-icon"
            :class="`bg-${stat.color}-100 text-${stat.color}-600`"
          >
            {{ stat.icon }}
          </div>
          <div class="stat-card-content">
            <h3>{{ stat.title }}</h3>
            <div class="stat-value">{{ formatNumber(stat.value) }}</div>
          </div>
        </div>
        <div class="stat-card-footer">
          <div 
            class="stat-trend"
            :class="stat.trend >= 0 ? 'positive' : 'negative'"
          >
            <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path 
                v-if="stat.trend >= 0"
                fill-rule="evenodd" 
                d="M3.293 9.707a1 1 0 010-1.414l6-6a1 1 0 011.414 0l6 6a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L4.707 9.707a1 1 0 01-1.414 0z" 
                clip-rule="evenodd" 
              />
              <path 
                v-else
                fill-rule="evenodd" 
                d="M16.707 10.293a1 1 0 010 1.414l-6 6a1 1 0 01-1.414 0l-6-6a1 1 0 111.414-1.414L9 14.586V3a1 1 0 012 0v11.586l4.293-4.293a1 1 0 011.414 0z" 
                clip-rule="evenodd" 
              />
            </svg>
            {{ Math.abs(stat.trend) }}%
          </div>
          <span class="text-gray-500 ml-2">vs last week</span>
        </div>
      </div>
    </div>

    <!-- Charts and Recent Activity -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
      <!-- Incident Trends Chart -->
      <div class="lg:col-span-2">
        <div class="card">
          <div class="card-header">
            <h2 class="card-title">Incident Trends</h2>
            <select class="form-input w-auto text-sm">
              <option>Last 7 days</option>
              <option>Last 30 days</option>
              <option>Last 90 days</option>
            </select>
          </div>
          <div class="card-body">
            <!-- Chart placeholder -->
            <div class="h-64 bg-gray-50 rounded-lg flex items-center justify-center">
              <div class="text-center">
                <div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mx-auto mb-3">
                  <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                  </svg>
                </div>
                <p class="text-sm text-gray-500">Chart visualization will be integrated here</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Severity Distribution -->
      <div>
        <div class="card">
          <div class="card-header">
            <h2 class="card-title">Severity Distribution</h2>
          </div>
          <div class="card-body">
            <div class="space-y-4">
              <div 
                v-for="severity in severityData" 
                :key="severity.name"
                class="flex items-center justify-between"
              >
                <div class="flex items-center">
                  <div 
                    class="w-3 h-3 rounded-full mr-3"
                    :class="`bg-${severity.color}-500`"
                  ></div>
                  <span class="text-sm font-medium text-gray-700">{{ severity.name }}</span>
                </div>
                <div class="flex items-center space-x-2">
                  <span class="text-sm font-bold text-gray-900">{{ severity.count }}</span>
                  <span class="text-xs text-gray-500">({{ severity.percentage }}%)</span>
                </div>
              </div>
            </div>
            
            <!-- Simple progress bars -->
            <div class="mt-4 space-y-2">
              <div 
                v-for="severity in severityData" 
                :key="`bar-${severity.name}`"
                class="flex items-center"
              >
                <div class="flex-1 bg-gray-200 rounded-full h-2 mr-3">
                  <div 
                    :class="`bg-${severity.color}-500 h-2 rounded-full transition-all duration-300`"
                    :style="`width: ${severity.percentage}%`"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Activity Tables -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Recent Incidents -->
      <div class="table-container">
        <div class="p-6 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900">Recent Incidents</h2>
            <RouterLink to="/incidents" class="text-sm text-blue-600 hover:text-blue-700 font-medium">
              View all →
            </RouterLink>
          </div>
        </div>
        <div class="p-6">
          <div class="space-y-4" v-if="recentIncidents.length > 0">
            <div 
              v-for="incident in recentIncidents" 
              :key="incident.id"
              class="flex items-start justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors cursor-pointer"
              @click="$router.push(`/incidents/${incident.id}`)"
            >
              <div class="flex-1">
                <div class="flex items-center space-x-2 mb-2">
                  <h4 class="font-medium text-gray-900 text-sm">{{ incident.title }}</h4>
                  <span 
                    class="badge"
                    :class="`badge-${incident.severity}`"
                  >
                    {{ incident.severity }}
                  </span>
                </div>
                <div class="flex items-center text-xs text-gray-500 space-x-4">
                  <span>{{ formatRelativeTime(incident.createdAt) }}</span>
                  <span v-if="incident.assignee">{{ incident.assignee.name }}</span>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            <p>No recent incidents</p>
          </div>
        </div>
      </div>

      <!-- Active Alerts -->
      <div class="table-container">
        <div class="p-6 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900">Active Alerts</h2>
            <RouterLink to="/alerts" class="text-sm text-blue-600 hover:text-blue-700 font-medium">
              View all →
            </RouterLink>
          </div>
        </div>
        <div class="p-6">
          <div class="space-y-4" v-if="activeAlerts.length > 0">
            <div 
              v-for="alert in activeAlerts" 
              :key="alert.id"
              class="flex items-start justify-between p-4 bg-red-50 rounded-lg border border-red-100"
            >
              <div class="flex-1">
                <div class="flex items-center space-x-2 mb-2">
                  <h4 class="font-medium text-gray-900 text-sm">{{ alert.summary }}</h4>
                  <span 
                    class="badge"
                    :class="`badge-${alert.severity}`"
                  >
                    {{ alert.severity }}
                  </span>
                </div>
                <div class="flex items-center text-xs text-gray-500 space-x-4">
                  <span>{{ formatRelativeTime(alert.startsAt) }}</span>
                  <span v-if="alert.source">{{ alert.source }}</span>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            <p>No active alerts</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useIncidents } from '@/composables/useIncidents'
import { useAlertsStore } from '@/stores/alerts'
import { formatDate, formatRelativeTime, formatNumber } from '@/utils/formatters'

const { incidents, fetchIncidents } = useIncidents()
const alertsStore = useAlertsStore()

// Computed stats
const statsData = computed(() => [
  {
    title: 'Total Incidents',
    value: incidents.value.length,
    icon: '🎯',
    color: 'blue',
    trend: 12,
    onClick: () => console.log('Navigate to incidents')
  },
  {
    title: 'Open Incidents',
    value: incidents.value.filter(i => i.status === 'open').length,
    icon: '🔴',
    color: 'red',
    trend: -5,
    onClick: () => console.log('Navigate to open incidents')
  },
  {
    title: 'Critical Alerts',
    value: alertsStore.criticalAlerts.length,
    icon: '⚠️',
    color: 'orange',
    trend: 8,
    onClick: () => console.log('Navigate to critical alerts')
  },
  {
    title: 'Avg Response Time',
    value: 4.2,
    icon: '⏱️',
    color: 'green',
    trend: -15,
    onClick: () => console.log('View response metrics')
  }
])

const severityData = computed(() => {
  const total = incidents.value.length || 1
  const counts = {
    critical: incidents.value.filter(i => i.severity === 'critical').length,
    high: incidents.value.filter(i => i.severity === 'high').length,
    medium: incidents.value.filter(i => i.severity === 'medium').length,
    low: incidents.value.filter(i => i.severity === 'low').length
  }

  return [
    { name: 'Critical', count: counts.critical, percentage: Math.round((counts.critical / total) * 100), color: 'red' },
    { name: 'High', count: counts.high, percentage: Math.round((counts.high / total) * 100), color: 'orange' },
    { name: 'Medium', count: counts.medium, percentage: Math.round((counts.medium / total) * 100), color: 'yellow' },
    { name: 'Low', count: counts.low, percentage: Math.round((counts.low / total) * 100), color: 'green' }
  ]
})

const recentIncidents = computed(() => {
  return incidents.value
    .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
    .slice(0, 5)
})

const activeAlerts = computed(() => {
  return alertsStore.firingAlerts.slice(0, 5)
})

// Initialize data
onMounted(async () => {
  await fetchIncidents()
  await alertsStore.fetchAlerts()
})
</script>

<style scoped>
.stat-card {
  cursor: pointer;
}

.stat-card-blue::before {
  background: linear-gradient(90deg, var(--color-primary-500), var(--color-primary-600));
}

.stat-card-red::before {
  background: linear-gradient(90deg, var(--color-error), #dc2626);
}

.stat-card-orange::before {
  background: linear-gradient(90deg, var(--color-warning), var(--color-high));
}

.stat-card-green::before {
  background: linear-gradient(90deg, var(--color-success), var(--color-low));
}
</style>