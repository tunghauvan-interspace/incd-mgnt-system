<template>
  <div class="fixed top-4 right-4 z-50 max-w-md w-full">
    <TransitionGroup name="notification" tag="div" class="space-y-2">
      <div
        v-for="notification in notifications"
        :key="notification.id"
        :class="getNotificationClass(notification.type)"
        class="p-4 rounded-lg shadow-lg border-l-4"
      >
        <div class="flex justify-between items-start">
          <div class="flex-1">
            <h4 class="font-medium">{{ notification.title }}</h4>
            <p class="text-sm mt-1">{{ notification.message }}</p>
          </div>
          <button
            @click="removeNotification(notification.id)"
            class="ml-4 text-gray-400 hover:text-gray-600 flex-shrink-0"
          >
            ×
          </button>
        </div>
        <div v-if="notification.action" class="mt-3">
          <button
            @click="notification.action.handler"
            class="text-sm font-medium text-blue-600 hover:text-blue-500"
          >
            {{ notification.action.label }}
          </button>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { useNotifications } from '@/composables/useNotifications'

const { notifications, remove } = useNotifications()

const removeNotification = (id: string) => {
  remove(id)
}

const getNotificationClass = (type: string) => {
  const classes = {
    success: 'bg-green-50 border-green-400 text-green-800',
    error: 'bg-red-50 border-red-400 text-red-800',
    warning: 'bg-yellow-50 border-yellow-400 text-yellow-800',
    info: 'bg-blue-50 border-blue-400 text-blue-800'
  }
  return classes[type] || classes.info
}
</script>

<style scoped>
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>