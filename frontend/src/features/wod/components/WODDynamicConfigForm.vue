<script setup lang="ts">
import { computed } from 'vue'
import BaseNumberInput from '@/shared/components/BaseNumberInput.vue'
import { WOD_CONFIG_SCHEMAS } from '@/features/wod/model/wodSchemas'
import type { WODFormConfig, WODType } from '@/features/wod/model/wodTypes'

const props = defineProps<{
  type: WODType
  config: WODFormConfig
}>()

const emit = defineEmits<{
  'update:field': [key: string, value: number | undefined]
}>()

const fields = computed(() => WOD_CONFIG_SCHEMAS[props.type])

function fieldValue(key: string): number | undefined {
  const value = Object.entries(props.config).find(([entryKey]) => entryKey === key)?.[1]
  return typeof value === 'number' ? value : undefined
}
</script>

<template>
  <div class="stack">
    <h3 class="section-title">Configuration</h3>
    <div class="row">
      <BaseNumberInput
        v-for="field in fields"
        :key="field.key"
        :label="field.label"
        :min="field.min"
        :model-value="fieldValue(field.key)"
        @update:model-value="emit('update:field', field.key, $event)"
      />
    </div>
  </div>
</template>
