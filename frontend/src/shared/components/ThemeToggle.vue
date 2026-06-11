<script setup lang="ts">
import { useTheme } from '@/shared/composables/useTheme'

const { isDark, toggleTheme } = useTheme()
</script>

<template>
  <button
    type="button"
    class="theme-toggle"
    :aria-label="isDark ? 'Switch to light theme' : 'Switch to dark theme'"
    @click="toggleTheme"
  >
    {{ isDark ? '☀' : '☾' }}
  </button>
</template>
