<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { Doughnut } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  CategoryScale,
  LinearScale
} from 'chart.js'

// Register Chart.js components
ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale)

interface Props {
  data: Record<string, number>
  title: string
}

const props = defineProps<Props>()

const chartData = ref({
  labels: [] as string[],
  datasets: [{
    data: [] as number[],
    backgroundColor: [
      '#e74c3c', // red for open/high
      '#f39c12', // orange for acknowledged/medium  
      '#27ae60', // green for resolved/low
      '#3498db', // blue for info/other
      '#9b59b6', // purple for additional categories
      '#1abc9c'  // teal for additional categories
    ],
    borderWidth: 1,
    borderColor: '#fff'
  }]
})

const chartOptions = ref({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    title: {
      display: true,
      text: props.title,
      font: {
        size: 16
      }
    },
    legend: {
      position: 'bottom' as const,
      labels: {
        padding: 15,
        usePointStyle: true
      }
    },
    tooltip: {
      callbacks: {
        label: function(context: any) {
          const label = context.label || ''
          const value = context.parsed || 0
          const total = context.dataset.data.reduce((a: number, b: number) => a + b, 0)
          const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : '0.0'
          return `${label}: ${value} (${percentage}%)`
        }
      }
    }
  },
  cutout: '50%' // Makes it a donut chart
})

// Update chart data when props change
watch(() => props.data, (newData) => {
  updateChartData(newData)
}, { immediate: true })

function updateChartData(data: Record<string, number>) {
  const labels = Object.keys(data)
  const values = Object.values(data)
  
  chartData.value = {
    labels,
    datasets: [{
      ...chartData.value.datasets[0],
      data: values
    }]
  }
}

onMounted(() => {
  updateChartData(props.data)
})
</script>

<template>
  <div class="chart-container">
    <Doughnut
      :data="chartData"
      :options="chartOptions"
      class="chart"
    />
  </div>
</template>

<style scoped>
.chart-container {
  position: relative;
  height: 350px;
  width: 100%;
}

.chart {
  max-height: 350px;
}

@media (max-width: 768px) {
  .chart-container {
    height: 250px;
  }
  
  .chart {
    max-height: 250px;
  }
}
</style>