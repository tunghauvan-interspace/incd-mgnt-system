<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <AppHeader @toggle-sidebar="toggleMobileSidebar" />

    <!-- Sidebar -->
    <AppSidebar 
      :is-collapsed="sidebarCollapsed"
      :is-mobile-open="mobileSidebarOpen"
      @toggle-collapsed="toggleSidebarCollapsed"
      @close-mobile="closeMobileSidebar"
    />

    <!-- Main Content -->
    <main 
      :class="[
        'transition-all duration-300 ease-in-out',
        'pt-16', // Top padding for fixed header
        sidebarCollapsed ? 'lg:ml-16' : 'lg:ml-72'
      ]"
    >
      <div class="p-6">
        <RouterView />
      </div>
    </main>

    <!-- Notifications -->
    <NotificationContainer />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterView } from 'vue-router'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import NotificationContainer from '@/components/ui/NotificationContainer.vue'

// Sidebar state
const sidebarCollapsed = ref(false)
const mobileSidebarOpen = ref(false)

// Methods
const toggleSidebarCollapsed = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const toggleMobileSidebar = () => {
  mobileSidebarOpen.value = !mobileSidebarOpen.value
}

const closeMobileSidebar = () => {
  mobileSidebarOpen.value = false
}
</script>