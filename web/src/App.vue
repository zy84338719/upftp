<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app'

const store = useAppStore()

onMounted(() => {
  // Check for initial config from server (injected by Go template)
  const configEl = document.getElementById('app-config')
  if (configEl) {
    try {
      const config = JSON.parse(configEl.textContent || '{}')
      store.initConfig(config)
    } catch (e) {
      console.error('Failed to parse config:', e)
    }
  }
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
