<script setup lang="ts">
import { computed } from 'vue'
import BaseSelect from '@/shared/components/BaseSelect.vue'
import type { WODType } from '@/features/wod/model/wodTypes'
import { WOD_TYPES } from '@/features/wod/model/wodTypes'

defineProps<{
  modelValue: WODType
}>()

const emit = defineEmits<{
  'update:modelValue': [value: WODType]
}>()

const options = computed(() =>
  WOD_TYPES.map((type) => ({
    value: type,
    label: type,
  })),
)
</script>

<template>
  <BaseSelect
    label="WOD Type"
    :model-value="modelValue"
    :options="options"
    @update:model-value="emit('update:modelValue', $event as WODType)"
  />
</template>
