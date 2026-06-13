<script setup lang="ts">
import { nextTick } from 'vue'
import AddStagePanel from '@/features/wod/components/AddStagePanel.vue'
import StageEditor from '@/features/wod/components/StageEditor.vue'
import type { MovementInput, StageFormState, StageFormat, StageKind, WODType } from '@/features/wod/model/wodTypes'

const props = defineProps<{
  stages: StageFormState[]
}>()

const emit = defineEmits<{
  addStage: [kind: StageKind]
  removeStage: [index: number]
  moveStageUp: [index: number]
  moveStageDown: [index: number]
  updateStageKind: [index: number, kind: StageKind]
  updateStageFormat: [index: number, format: StageFormat]
  updateStageInstructions: [index: number, value: string]
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

async function addStage(kind: StageKind) {
  const nextIndex = props.stages.length
  emit('addStage', kind)
  await nextTick()
  document.getElementById(`stage-${nextIndex}`)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}
</script>

<template>
  <div class="stage-timeline">
    <div v-for="(stage, index) in stages" :key="index" class="stage-timeline__item">
      <span class="stage-timeline__node">{{ index + 1 }}</span>
      <StageEditor
        :index="index"
        :stage="stage"
        :can-remove="stages.length > 1"
        :can-move-up="index > 0"
        :can-move-down="index < stages.length - 1"
        @remove="emit('removeStage', index)"
        @move-up="emit('moveStageUp', index)"
        @move-down="emit('moveStageDown', index)"
        @update:kind="(kind) => emit('updateStageKind', index, kind)"
        @update:format="(format) => emit('updateStageFormat', index, format)"
        @update:instructions="(value) => emit('updateStageInstructions', index, value)"
        @update:type="(type) => emit('updateStageType', index, type)"
        @update:config-field="(key, value) => emit('updateStageConfigField', index, key, value)"
        @add-movement="emit('addMovement', index)"
        @remove-movement="(movementIndex) => emit('removeMovement', index, movementIndex)"
        @update:movement="(movementIndex, field, value) => emit('updateMovement', index, movementIndex, field, value)"
      />
    </div>

    <div class="stage-timeline__add">
      <AddStagePanel @add="addStage" />
    </div>
  </div>
</template>
