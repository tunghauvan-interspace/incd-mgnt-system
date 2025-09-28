<script setup lang="ts">
import { RouterView } from 'vue-router'
import Navbar from './components/Navbar.vue'
</script>

<template>
  <div id="app">
    <Navbar />
    <main>
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  line-height: 1.6;
  color: var(--color-text-primary);
  background-color: var(--color-bg-secondary);
  min-height: 100vh;
}

main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}
</style>
