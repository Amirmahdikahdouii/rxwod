<script setup lang="ts">
import ProgramOutlinePanel from '@/features/wod/components/ProgramOutlinePanel.vue'
import StageListEditor from '@/features/wod/components/StageListEditor.vue'
import { useWODForm } from '@/features/wod/composables/useWODForm'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseTextarea from '@/shared/components/BaseTextarea.vue'
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

function routeWODId() {
  const id = route.params.id
  return Array.isArray(id) ? id[0] : id
}

const initialWODId = routeWODId()

const {
  mode,
  name,
  description,
  stages,
  loading,
  initialLoading,
  error,
  result,
  addStage,
  removeStage,
  moveStageUp,
  moveStageDown,
  updateStageKind,
  updateStageFormat,
  updateStageInstructions,
  updateStageType,
  updateStageConfigField,
  addMovement,
  removeMovement,
  updateMovement,
  initEdit,
  submit,
} = useWODForm(initialWODId ? 'edit' : 'create')

const successSummary = computed(() => {
  if (!result.value) {
    return ''
  }
  const stageCount = 'stageCount' in result.value ? result.value.stageCount : result.value.stages.length
  const stageLabels = result.value.stages
    .map((stage) => `${stage.kind}/${stage.type}`)
    .join(' -> ')
  const action = mode.value === 'edit' ? 'updated' : 'saved'
  return `${result.value.name} ${action} with ${stageCount} stage(s): ${stageLabels}.`
})

const stageCountLabel = computed(() => `${stages.value.length} stage${stages.value.length === 1 ? '' : 's'}`)
const pageTitle = computed(() => (mode.value === 'edit' ? 'Edit WOD Program' : 'Create WOD Program'))
const pageSubtitle = computed(() =>
  mode.value === 'edit'
    ? 'Update your class plan stages, scoring, and prescriptions.'
    : 'Build your class plan with instructions and free-text prescriptions.',
)
const submitLabel = computed(() => (mode.value === 'edit' ? 'Save Changes' : 'Save Program'))

onMounted(async () => {
  if (initialWODId) {
    await initEdit(initialWODId)
  }
})
</script>

<template>
  <div class="container stack-lg">
    <header class="page-header">
      <h1 class="page-title">{{ pageTitle }}</h1>
      <p class="page-subtitle">{{ pageSubtitle }}</p>
    </header>

    <p v-if="initialLoading" class="loading-state">Loading program...</p>

    <div v-else class="create-layout">
      <div class="create-layout__main">
        <form id="wod-create-form" class="card stack-lg" @submit.prevent="submit">
          <section class="card-section stack">
            <div class="section-heading-row">
              <h2 class="section-title">Basics</h2>
              <span class="count-chip">{{ stageCountLabel }}</span>
            </div>
            <BaseInput v-model="name" label="Name" placeholder="Monday Session" />
            <BaseTextarea
              v-model="description"
              label="Description"
              placeholder="Full class plan notes"
              :rows="3"
            />
          </section>

          <section class="card-section stack">
            <div>
              <h2 class="section-title">Program Stages</h2>
              <p class="page-subtitle page-subtitle--flush">
                Add ordered stages with instructions and prescriptions, or structured metcon scoring.
              </p>
            </div>
            <StageListEditor
              :stages="stages"
              @add-stage="addStage"
              @remove-stage="removeStage"
              @move-stage-up="moveStageUp"
              @move-stage-down="moveStageDown"
              @update-stage-kind="updateStageKind"
              @update-stage-format="updateStageFormat"
              @update-stage-instructions="updateStageInstructions"
              @update-stage-type="updateStageType"
              @update-stage-config-field="updateStageConfigField"
              @add-movement="addMovement"
              @remove-movement="removeMovement"
              @update-movement="updateMovement"
            />
          </section>

          <div class="program-outline__mobile-actions stack">
            <div v-if="error" class="alert alert--error" role="alert">{{ error }}</div>
            <div v-if="result" class="alert alert--success" role="status">{{ successSummary }}</div>
            <button type="submit" class="btn-full" :disabled="loading">
              {{ loading ? 'Saving...' : submitLabel }}
            </button>
          </div>
        </form>
      </div>

      <div class="create-layout__aside">
        <ProgramOutlinePanel
          :name="name"
          :description="description"
          :stages="stages"
          :loading="loading"
          :error="error"
          :success-summary="successSummary"
          :submit-label="submitLabel"
        />
      </div>
    </div>
  </div>
</template>
