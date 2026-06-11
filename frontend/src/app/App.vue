<script setup lang="ts">
import AppHeader from '@/shared/components/AppHeader.vue'
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="app-shell">
    <AppHeader />
    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>
