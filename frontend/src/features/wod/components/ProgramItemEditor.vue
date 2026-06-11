<script setup lang="ts">
import { ref } from 'vue'
import ProgramItemPreview from '@/features/wod/components/ProgramItemPreview.vue'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseNumberInput from '@/shared/components/BaseNumberInput.vue'
import BaseSelect from '@/shared/components/BaseSelect.vue'
import BaseTextarea from '@/shared/components/BaseTextarea.vue'
import type { LoadUnit, MovementInput } from '@/features/wod/model/wodTypes'

defineProps<{
  movements: MovementInput[]
  openFormat: boolean
}>()

const emit = defineEmits<{
  add: []
  remove: [index: number]
  'update:movement': [index: number, field: keyof MovementInput, value: string | number | undefined]
}>()

const expandedAdvanced = ref<Record<number, boolean>>({})

const loadUnitOptions = [
  { value: 'kg', label: 'kg' },
  { value: 'lb', label: 'lb' },
  { value: 'bodyweight', label: 'bodyweight' },
]

function toggleAdvanced(index: number) {
  expandedAdvanced.value[index] = !expandedAdvanced.value[index]
}
</script>

<template>
  <div class="stack">
    <div class="row row--align-center row--between">
      <h3 class="section-title" style="margin: 0">Items</h3>
      <button type="button" class="secondary" @click="emit('add')">Add Item</button>
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
          label="Label"
          :model-value="movement.label ?? ''"
          placeholder="A"
          @update:model-value="emit('update:movement', index, 'label', $event)"
        />
        <BaseInput
          :label="`Item ${index + 1} name`"
          :model-value="movement.name"
          placeholder="e.g. Close Grip Bench Press"
          @update:model-value="emit('update:movement', index, 'name', $event)"
        />
      </div>

      <BaseTextarea
        label="Prescription"
        :model-value="movement.prescription ?? ''"
        placeholder="e.g. Accumulate 20 reps. *Use a 20kg and 15kg plate."
        :rows="openFormat ? 4 : 3"
        @update:model-value="emit('update:movement', index, 'prescription', $event)"
      />

      <ProgramItemPreview :movement="movement" />

      <button type="button" class="secondary" @click="toggleAdvanced(index)">
        {{ expandedAdvanced[index] ? 'Hide advanced fields' : 'Show advanced fields' }}
      </button>

      <div v-if="expandedAdvanced[index]" class="stack">
        <div class="row">
          <BaseNumberInput
            label="Reps"
            :min="1"
            :model-value="movement.reps"
            @update:model-value="emit('update:movement', index, 'reps', $event)"
          />
          <BaseNumberInput
            label="Load"
            :min="0"
            :model-value="movement.loadValue"
            @update:model-value="emit('update:movement', index, 'loadValue', $event)"
          />
          <BaseSelect
            label="Load unit"
            :model-value="movement.loadUnit ?? ''"
            :options="loadUnitOptions"
            @update:model-value="emit('update:movement', index, 'loadUnit', $event as LoadUnit)"
          />
        </div>
        <BaseTextarea
          label="Notes"
          :model-value="movement.notes ?? ''"
          placeholder="Equipment or coaching notes"
          :rows="2"
          @update:model-value="emit('update:movement', index, 'notes', $event)"
        />
      </div>

      <button type="button" class="secondary" @click="emit('remove', index)">Remove</button>
    </div>
  </div>
</template>
