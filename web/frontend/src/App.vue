<template>
  <div id="app">
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<style scoped>
#app {
  min-height: 100vh;
  font-family: var(--font-sans);
  background-color: var(--color-gray-50);
}
</style>