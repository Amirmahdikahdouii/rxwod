<script setup lang="ts">
import { computed } from 'vue'
import BaseSelect from '@/shared/components/BaseSelect.vue'
import type { StageKind } from '@/features/wod/model/wodTypes'
import { STAGE_KINDS } from '@/features/wod/model/wodTypes'

defineProps<{
  modelValue: StageKind
}>()

const emit = defineEmits<{
  'update:modelValue': [value: StageKind]
}>()

const options = computed(() =>
  STAGE_KINDS.map((kind) => ({
    value: kind,
    label: kind.charAt(0) + kind.slice(1).toLowerCase(),
  })),
)
</script>

<template>
  <BaseSelect
    label="Stage Kind"
    :model-value="modelValue"
    :options="options"
    @update:model-value="emit('update:modelValue', $event as StageKind)"
  />
</template>
