<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'default' | 'success' | 'warning' | 'error' | 'info' | 'muted'
    size?: 'sm' | 'md'
  }>(),
  {
    variant: 'default',
    size: 'sm',
  },
)
</script>

<template>
  <span class="ui-pill" :class="[`ui-pill--${variant}`, `ui-pill--${size}`]">
    <slot />
  </span>
</template>

<style src="./ui.css"></style>
