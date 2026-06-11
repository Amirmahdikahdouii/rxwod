<script setup lang="ts">
import MovementListEditor from '@/features/wod/components/MovementListEditor.vue'
import ScorePreview from '@/features/wod/components/ScorePreview.vue'
import StageKindSelector from '@/features/wod/components/StageKindSelector.vue'
import WODDynamicConfigForm from '@/features/wod/components/WODDynamicConfigForm.vue'
import WODTypeSelector from '@/features/wod/components/WODTypeSelector.vue'
import { STAGE_KIND_BADGE_CLASS } from '@/features/wod/model/wodTheme'
import type { MovementInput, StageFormState, StageKind, WODType } from '@/features/wod/model/wodTypes'

defineProps<{
  index: number
  stage: StageFormState
  canRemove: boolean
  canMoveUp: boolean
  canMoveDown: boolean
}>()

const emit = defineEmits<{
  remove: []
  moveUp: []
  moveDown: []
  'update:kind': [value: StageKind]
  'update:type': [value: WODType]
  'update:configField': [key: string, value: number | undefined]
  addMovement: []
  removeMovement: [movementIndex: number]
  'update:movement': [movementIndex: number, field: keyof MovementInput, value: string | number | undefined]
}>()
</script>

<template>
  <article class="stage-editor stack">
    <header class="stage-editor__header">
      <div class="row row--align-center" style="gap: 0.75rem">
        <h3 class="stage-editor__title">Stage {{ index + 1 }}</h3>
        <span class="badge" :class="STAGE_KIND_BADGE_CLASS[stage.kind]">{{ stage.kind }}</span>
      </div>
      <div class="stage-editor__actions">
        <button type="button" class="secondary" :disabled="!canMoveUp" @click="emit('moveUp')">
          Up
        </button>
        <button type="button" class="secondary" :disabled="!canMoveDown" @click="emit('moveDown')">
          Down
        </button>
        <button type="button" class="secondary" :disabled="!canRemove" @click="emit('remove')">
          Remove
        </button>
      </div>
    </header>

    <div class="type-config-grid">
      <StageKindSelector :model-value="stage.kind" @update:model-value="emit('update:kind', $event)" />
      <WODTypeSelector :model-value="stage.type" @update:model-value="emit('update:type', $event)" />
    </div>

    <ScorePreview :type="stage.type" />

    <WODDynamicConfigForm
      :type="stage.type"
      :config="stage.config"
      @update:field="(key, value) => emit('update:configField', key, value)"
    />

    <MovementListEditor
      :movements="stage.movements"
      @add="emit('addMovement')"
      @remove="(movementIndex) => emit('removeMovement', movementIndex)"
      @update:movement="(movementIndex, field, value) => emit('update:movement', movementIndex, field, value)"
    />
  </article>
</template>
