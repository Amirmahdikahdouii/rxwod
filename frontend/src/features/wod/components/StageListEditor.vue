<script setup lang="ts">
import StageEditor from '@/features/wod/components/StageEditor.vue'
import type { MovementInput, StageFormState, StageKind, WODType } from '@/features/wod/model/wodTypes'

defineProps<{
  stages: StageFormState[]
}>()

const emit = defineEmits<{
  addStage: []
  removeStage: [index: number]
  moveStageUp: [index: number]
  moveStageDown: [index: number]
  updateStageKind: [index: number, kind: StageKind]
  updateStageType: [index: number, type: WODType]
  updateStageConfigField: [index: number, key: string, value: number | undefined]
  addMovement: [stageIndex: number]
  removeMovement: [stageIndex: number, movementIndex: number]
  updateMovement: [
    stageIndex: number,
    movementIndex: number,
    field: keyof MovementInput,
    value: string | number | undefined,
  ]
}>()
</script>

<template>
  <div class="stack-lg">
    <StageEditor
      v-for="(stage, index) in stages"
      :key="index"
      :index="index"
      :stage="stage"
      :can-remove="stages.length > 1"
      :can-move-up="index > 0"
      :can-move-down="index < stages.length - 1"
      @remove="emit('removeStage', index)"
      @move-up="emit('moveStageUp', index)"
      @move-down="emit('moveStageDown', index)"
      @update:kind="(kind) => emit('updateStageKind', index, kind)"
      @update:type="(type) => emit('updateStageType', index, type)"
      @update:config-field="(key, value) => emit('updateStageConfigField', index, key, value)"
      @add-movement="emit('addMovement', index)"
      @remove-movement="(movementIndex) => emit('removeMovement', index, movementIndex)"
      @update:movement="(movementIndex, field, value) => emit('updateMovement', index, movementIndex, field, value)"
    />

    <button type="button" class="secondary" @click="emit('addStage')">Add Stage</button>
  </div>
</template>
