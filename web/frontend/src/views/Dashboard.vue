<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { metricsAPI } from '@/services/api'
import { formatDuration } from '@/utils/format'
import DoughnutChart from '@/components/DoughnutChart.vue'
import type { Metrics } from '@/types/api'

const metrics = ref<Metrics | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const loadDashboard = async () => {
  try {
    loading.value = true
    error.value = null
    metrics.value = await metricsAPI.getMetrics()
  } catch (err) {
    error.value = 'Error loading dashboard data'
    console.error('Error loading dashboard:', err)
  } finally {
    loading.value = false
  }
}

const refreshDashboard = () => {
  loadDashboard()
}

onMounted(() => {
  loadDashboard()
})
</script>

<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h2>Dashboard</h2>
      <div class="refresh-btn">
        <button @click="refreshDashboard" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div v-else-if="loading" class="loading">
      Loading dashboard...
    </div>

    <div v-else-if="metrics" class="metrics-grid">
      <div class="metric-card">
        <div class="metric-title">Total Incidents</div>
        <div class="metric-value">{{ metrics.total_incidents || 0 }}</div>
      </div>
      
      <div class="metric-card">
        <div class="metric-title">Open Incidents</div>
        <div class="metric-value critical">{{ metrics.open_incidents || 0 }}</div>
      </div>
      
      <div class="metric-card">
        <div class="metric-title">MTTA</div>
        <div class="metric-value">{{ formatDuration(metrics.mtta) }}</div>
        <div class="metric-subtitle">Mean Time To Acknowledge</div>
      </div>
      
      <div class="metric-card">
        <div class="metric-title">MTTR</div>
        <div class="metric-value">{{ formatDuration(metrics.mttr) }}</div>
        <div class="metric-subtitle">Mean Time To Resolve</div>
      </div>
    </div>

    <!-- Charts section -->
    <div v-if="metrics && !loading" class="charts-section">
      <div class="charts-grid">
        <div class="card">
          <DoughnutChart 
            :data="metrics.incidents_by_status || {}"
            title="Incidents by Status"
          />
        </div>
        
        <div class="card">
          <DoughnutChart 
            :data="metrics.incidents_by_severity || {}"
            title="Incidents by Severity"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  padding: 20px 0;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.dashboard-header h2 {
  color: #2c3e50;
  margin: 0;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.metric-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  text-align: center;
}

.metric-title {
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.metric-value {
  font-size: 2rem;
  font-weight: bold;
  color: #2c3e50;
  margin-bottom: 0.25rem;
}

.metric-value.critical {
  color: #e74c3c;
}

.metric-subtitle {
  font-size: 0.8rem;
  color: #999;
}

.charts-section {
  margin-top: 2rem;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 2rem;
}

.error-message {
  background: #ffeaea;
  color: #e74c3c;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
  border: 1px solid #f8cecc;
}

@media (max-width: 768px) {
  .dashboard-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .metrics-grid {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }
  
  .charts-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  
  .metric-value {
    font-size: 1.5rem;
  }
}
</style>