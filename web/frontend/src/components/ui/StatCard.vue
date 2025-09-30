<template>
  <div 
    class="card"
    :class="getColorClass(color)"
  >
    <div class="flex items-center">
      <div class="flex-shrink-0">
        <span class="text-2xl">{{ icon }}</span>
      </div>
      <div class="ml-4">
        <p class="text-sm font-medium text-gray-600">{{ title }}</p>
        <p class="text-2xl font-bold text-gray-900">{{ formattedValue }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { formatNumber } from '@/utils/formatters'

interface Props {
  title: string
  value: number
  icon: string
  color?: 'blue' | 'red' | 'green' | 'orange' | 'gray'
}

const props = withDefaults(defineProps<Props>(), {
  color: 'blue'
})

const formattedValue = computed(() => formatNumber(props.value))

const getColorClass = (color: string) => {
  const classes = {
    blue: 'border-blue-200 bg-blue-50',
    red: 'border-red-200 bg-red-50',
    green: 'border-green-200 bg-green-50',
    orange: 'border-orange-200 bg-orange-50',
    gray: 'border-gray-200 bg-gray-50'
  }
  return classes[color] || classes.blue
}
</script>