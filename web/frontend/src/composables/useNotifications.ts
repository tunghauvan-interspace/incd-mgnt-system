import { computed } from 'vue'
import { useNotificationsStore } from '@/stores/notifications'

export const useNotifications = () => {
  const notificationsStore = useNotificationsStore()

  const notifications = computed(() => notificationsStore.notifications)

  const success = (message: string, title?: string) => {
    return notificationsStore.success(message, title)
  }

  const error = (message: string, title?: string) => {
    return notificationsStore.error(message, title)
  }

  const warning = (message: string, title?: string) => {
    return notificationsStore.warning(message, title)
  }

  const info = (message: string, title?: string) => {
    return notificationsStore.info(message, title)
  }

  const remove = (id: string) => {
    notificationsStore.removeNotification(id)
  }

  const clear = () => {
    notificationsStore.clearAllNotifications()
  }

  return {
    notifications,
    success,
    error,
    warning,
    info,
    remove,
    clear
  }
}