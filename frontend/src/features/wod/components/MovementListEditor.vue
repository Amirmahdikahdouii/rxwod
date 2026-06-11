<script setup lang="ts">
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseNumberInput from '@/shared/components/BaseNumberInput.vue'
import type { MovementInput } from '@/features/wod/model/wodTypes'

defineProps<{
  movements: MovementInput[]
}>()

const emit = defineEmits<{
  add: []
  remove: [index: number]
  'update:movement': [index: number, field: keyof MovementInput, value: string | number | undefined]
}>()
</script>

<template>
  <div class="stack">
    <div class="row row--align-center row--between">
      <h3 class="section-title" style="margin: 0">Movements</h3>
      <button type="button" class="secondary" @click="emit('add')">Add Movement</button>
    </div>

    <div
      v-for="(movement, index) in movements"
      :key="index"
      class="movement-item stack"
      style="margin-left: 0.5rem"
    >
      <span class="movement-item__number">{{ index + 1 }}</span>
      <div class="row">
        <BaseInput
          :label="`Movement ${index + 1}`"
          :model-value="movement.name"
          placeholder="e.g. Burpee"
          @update:model-value="emit('update:movement', index, 'name', $event)"
        />
        <BaseNumberInput
          label="Reps"
          :min="1"
          :model-value="movement.reps"
          @update:model-value="emit('update:movement', index, 'reps', $event)"
        />
      </div>
      <button type="button" class="secondary" @click="emit('remove', index)">Remove</button>
    </div>
  </div>
</template>
