import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Notification } from '@/types/notifications'

export const useNotificationsStore = defineStore('notifications', () => {
  // State
  const notifications = ref<Notification[]>([])

  // Actions
  const addNotification = (notification: Omit<Notification, 'id' | 'timestamp'>) => {
    const newNotification: Notification = {
      id: Date.now().toString(),
      timestamp: new Date(),
      ...notification
    }
    
    notifications.value.unshift(newNotification)
    
    // Auto-remove after timeout (except for persistent notifications)
    if (!notification.persistent) {
      const timeout = notification.type === 'error' ? 8000 : 5000
      setTimeout(() => {
        removeNotification(newNotification.id)
      }, timeout)
    }
    
    return newNotification
  }

  const removeNotification = (id: string) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
  }

  const clearAllNotifications = () => {
    notifications.value = []
  }

  // Convenience methods
  const success = (message: string, title?: string) => {
    return addNotification({
      type: 'success',
      title: title || 'Success',
      message
    })
  }

  const error = (message: string, title?: string) => {
    return addNotification({
      type: 'error',
      title: title || 'Error',
      message,
      persistent: true
    })
  }

  const warning = (message: string, title?: string) => {
    return addNotification({
      type: 'warning',
      title: title || 'Warning',
      message
    })
  }

  const info = (message: string, title?: string) => {
    return addNotification({
      type: 'info',
      title: title || 'Information',
      message
    })
  }

  return {
    // State
    notifications,
    // Actions
    addNotification,
    removeNotification,
    clearAllNotifications,
    // Convenience methods
    success,
    error,
    warning,
    info
  }
})