<script setup lang="ts">
import { RouterLink } from 'vue-router'
import ThemeToggle from '@/shared/components/ThemeToggle.vue'
</script>

<template>
  <header class="app-header">
    <div class="app-header__inner">
      <RouterLink to="/" class="app-header__brand">RXWOD</RouterLink>

      <nav class="app-header__nav" aria-label="Main navigation">
        <RouterLink to="/" class="app-header__link">Create</RouterLink>
        <RouterLink to="/wods" class="app-header__link">Saved</RouterLink>
      </nav>

      <div class="app-header__actions">
        <ThemeToggle />
      </div>
    </div>
  </header>
</template>
