<template>
  <div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <div class="text-center">
        <h1 class="text-3xl font-bold text-gray-900">
          Incident Management
        </h1>
        <p class="mt-2 text-sm text-gray-600">
          Manage incidents and alerts efficiently
        </p>
      </div>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
        <RouterView />
      </div>
    </div>

    <!-- Notifications -->
    <NotificationContainer />
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import NotificationContainer from '@/components/ui/NotificationContainer.vue'
</script>