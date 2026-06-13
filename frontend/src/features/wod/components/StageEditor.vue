<script setup lang="ts">
import { computed, ref } from 'vue'
import ProgramItemEditor from '@/features/wod/components/ProgramItemEditor.vue'
import ScorePreview from '@/features/wod/components/ScorePreview.vue'
import StageKindSelector from '@/features/wod/components/StageKindSelector.vue'
import WODDynamicConfigForm from '@/features/wod/components/WODDynamicConfigForm.vue'
import WODTypeSelector from '@/features/wod/components/WODTypeSelector.vue'
import { STAGE_KIND_BADGE_CLASS } from '@/features/wod/model/wodTheme'
import type { MovementInput, StageFormState, StageFormat, StageKind, WODType } from '@/features/wod/model/wodTypes'
import { isOpenFormat } from '@/features/wod/model/wodTypes'
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

const stageFormat = computed<StageFormat>(() => (isOpenFormat(props.stage.type) ? 'OPEN' : 'STRUCTURED'))
const openFormat = computed(() => stageFormat.value === 'OPEN')
const expanded = ref(true)
const stageId = computed(() => `stage-${props.index}`)
const bodyId = computed(() => `stage-${props.index}-body`)

function stageKindLabel(kind: StageKind) {
  return kind.charAt(0) + kind.slice(1).toLowerCase()
}

function setFormat(format: StageFormat) {
  emit('update:format', format)
}
</script>

<template>
  <article :id="stageId" class="stage-editor stack">
    <header class="stage-editor__header">
      <button
        type="button"
        class="stage-editor__summary"
        :aria-expanded="expanded"
        :aria-controls="bodyId"
        @click="expanded = !expanded"
      >
        <span class="stage-editor__summary-icon">{{ expanded ? '-' : '+' }}</span>
        <span>
          <p class="stage-editor__kicker">Stage {{ index + 1 }}</p>
          <span class="stage-editor__identity">
            <h3 class="stage-editor__title">{{ stageKindLabel(stage.kind) }}</h3>
            <span class="badge" :class="STAGE_KIND_BADGE_CLASS[stage.kind]">{{ stageKindLabel(stage.kind) }}</span>
          </span>
        </span>
      </button>
      <div class="stage-editor__actions">
        <button type="button" class="icon-btn" :disabled="!canMoveUp" title="Move stage up" @click="emit('moveUp')">
          ^
        </button>
        <button type="button" class="icon-btn" :disabled="!canMoveDown" title="Move stage down" @click="emit('moveDown')">
          v
        </button>
        <button type="button" class="icon-btn" :disabled="!canRemove" title="Remove stage" @click="emit('remove')">
          x
        </button>
      </div>
    </header>

    <div v-show="expanded" :id="bodyId" class="stage-editor__body stack">
      <section>
        <p class="subsection-label">Configuration</p>
        <div class="type-config-grid">
          <StageKindSelector :model-value="stage.kind" @update:model-value="emit('update:kind', $event)" />
          <div class="field">
            <label>Format</label>
            <div class="segmented-control" role="group" aria-label="Stage format">
              <button
                type="button"
                class="segmented-control__option"
                :class="{ 'segmented-control__option--active': stageFormat === 'OPEN' }"
                @click="setFormat('OPEN')"
              >
                Open
              </button>
              <button
                type="button"
                class="segmented-control__option"
                :class="{ 'segmented-control__option--active': stageFormat === 'STRUCTURED' }"
                @click="setFormat('STRUCTURED')"
              >
                Structured
              </button>
            </div>
          </div>
        </div>
      </section>

      <section>
        <p class="subsection-label">Instructions</p>
        <BaseTextarea
          label="Stage instructions"
          :model-value="stage.instructions"
          placeholder="e.g. Complete in 20 minutes."
          :rows="2"
          @update:model-value="emit('update:instructions', $event)"
        />
      </section>

      <template v-if="openFormat">
        <p class="page-subtitle page-subtitle--flush">
          Open format uses free-text prescriptions. Add items with labels, names, and coaching details.
        </p>
      </template>

      <section v-else class="stack">
        <p class="subsection-label">Scoring</p>
        <WODTypeSelector :model-value="stage.type" @update:model-value="emit('update:type', $event)" />
        <ScorePreview :type="stage.type" />
        <WODDynamicConfigForm
          :type="stage.type"
          :config="stage.config"
          @update:field="(key, value) => emit('update:configField', key, value)"
        />
      </section>

      <ProgramItemEditor
        :movements="stage.movements"
        :open-format="openFormat"
        @add="emit('addMovement')"
        @remove="(movementIndex) => emit('removeMovement', movementIndex)"
        @update:movement="(movementIndex, field, value) => emit('update:movement', movementIndex, field, value)"
      />
    </div>
  </article>
</template>
