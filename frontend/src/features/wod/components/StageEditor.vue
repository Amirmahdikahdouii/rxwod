<script setup lang="ts">
import { computed } from 'vue'
import ProgramItemEditor from '@/features/wod/components/ProgramItemEditor.vue'
import ScorePreview from '@/features/wod/components/ScorePreview.vue'
import StageKindSelector from '@/features/wod/components/StageKindSelector.vue'
import WODDynamicConfigForm from '@/features/wod/components/WODDynamicConfigForm.vue'
import WODTypeSelector from '@/features/wod/components/WODTypeSelector.vue'
import { STAGE_KIND_BADGE_CLASS } from '@/features/wod/model/wodTheme'
import type { MovementInput, StageFormState, StageFormat, StageKind, WODType } from '@/features/wod/model/wodTypes'
import { isOpenFormat } from '@/features/wod/model/wodTypes'
import BaseSelect from '@/shared/components/BaseSelect.vue'
import BaseTextarea from '@/shared/components/BaseTextarea.vue'

const props = defineProps<{
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
  'update:format': [value: StageFormat]
  'update:type': [value: WODType]
  'update:instructions': [value: string]
  'update:configField': [key: string, value: number | undefined]
  addMovement: []
  removeMovement: [movementIndex: number]
  'update:movement': [movementIndex: number, field: keyof MovementInput, value: string | number | undefined]
}>()

const formatOptions = [
  { value: 'OPEN', label: 'Open (prescription)' },
  { value: 'STRUCTURED', label: 'Structured (scored WOD)' },
]

const stageFormat = computed<StageFormat>(() => (isOpenFormat(props.stage.type) ? 'OPEN' : 'STRUCTURED'))
const openFormat = computed(() => stageFormat.value === 'OPEN')

function onFormatChange(value: string) {
  emit('update:format', value as StageFormat)
}
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
      <BaseSelect
        label="Format"
        :model-value="stageFormat"
        :options="formatOptions"
        @update:model-value="onFormatChange"
      />
    </div>

    <BaseTextarea
      label="Stage instructions"
      :model-value="stage.instructions"
      placeholder="e.g. Complete in 20 minutes."
      :rows="2"
      @update:model-value="emit('update:instructions', $event)"
    />

    <template v-if="openFormat">
      <p class="page-subtitle" style="margin: 0">
        Open format uses free-text prescriptions. Add items with labels, names, and coaching details.
      </p>
    </template>

    <template v-else>
      <WODTypeSelector :model-value="stage.type" @update:model-value="emit('update:type', $event)" />
      <ScorePreview :type="stage.type" />
      <WODDynamicConfigForm
        :type="stage.type"
        :config="stage.config"
        @update:field="(key, value) => emit('update:configField', key, value)"
      />
    </template>

    <ProgramItemEditor
      :movements="stage.movements"
      :open-format="openFormat"
      @add="emit('addMovement')"
      @remove="(movementIndex) => emit('removeMovement', movementIndex)"
      @update:movement="(movementIndex, field, value) => emit('update:movement', movementIndex, field, value)"
    />
  </article>
</template>
