<script setup lang="ts">
import { onBeforeUnmount, onMounted } from 'vue'

const props = defineProps<{
  message: string
  variant?: 'success' | 'error'
  duration?: number
}>()

const emit = defineEmits<{
  dismiss: []
}>()

let timer: ReturnType<typeof setTimeout> | null = null

onMounted(() => {
  timer = setTimeout(() => emit('dismiss'), props.duration ?? 3500)
})

onBeforeUnmount(() => {
  if (timer) {
    clearTimeout(timer)
  }
})
</script>

<template>
  <div
    class="toast"
    :class="variant === 'error' ? 'toast--error' : 'toast--success'"
    role="status"
    aria-live="polite"
  >
    {{ message }}
  </div>
</template>
