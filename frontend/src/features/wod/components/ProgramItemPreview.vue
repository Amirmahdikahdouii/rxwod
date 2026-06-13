<script setup lang="ts">
import { computed } from 'vue'
import type { MovementInput } from '@/features/wod/model/wodTypes'

const props = defineProps<{
  movement: MovementInput
}>()

const previewLine = computed(() => {
  const parts: string[] = []
  if (props.movement.label?.trim()) {
    parts.push(`${props.movement.label.trim()}.`)
  }
  if (props.movement.name.trim()) {
    parts.push(props.movement.name.trim())
  }
  if (props.movement.prescription?.trim()) {
    parts.push(props.movement.prescription.trim())
  }
  if (props.movement.sets && props.movement.reps) {
    parts.push(`${props.movement.sets} x ${props.movement.reps}`)
  }
  return parts.join(' ')
})
</script>

<template>
  <p v-if="previewLine" class="item-preview">{{ previewLine }}</p>
</template>
