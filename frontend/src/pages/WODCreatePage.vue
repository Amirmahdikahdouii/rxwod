<script setup lang="ts">
import ProgramOutlinePanel from '@/features/wod/components/ProgramOutlinePanel.vue'
import StageListEditor from '@/features/wod/components/StageListEditor.vue'
import { useSession } from '@/features/auth/composables/useSession'
import {
  clearStoredDraft,
  draftStorageKey,
  draftsMatch,
  formatDraftAge,
  hasMeaningfulDraft,
  loadStoredDraft,
  saveStoredDraft,
  type StoredWODDraft,
} from '@/features/wod/composables/useWODDraftStorage'
import { useWODForm } from '@/features/wod/composables/useWODForm'
import type { WODDetail } from '@/features/wod/model/wodTypes'
import { canCreateWOD, canEditWOD, canPublishWOD, ROLE_LABELS } from '@/features/workspace/model/workspaceTypes'
import BaseInput from '@/shared/components/BaseInput.vue'
import BaseTextarea from '@/shared/components/BaseTextarea.vue'
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

const route = useRoute()
const session = useSession()

function routeWODId() {
  const id = route.params.id
  return Array.isArray(id) ? id[0] : id
}

const initialWODId = routeWODId()
const loadedDetail = ref<WODDetail | null>(null)
const pendingDraft = ref<StoredWODDraft | null>(null)
const showRecoveryPrompt = ref(false)
let autosaveTimer: ReturnType<typeof setTimeout> | null = null

const {
  mode,
  wodId,
  name,
  description,
  scheduledDate,
  stages,
  loading,
  savingAction,
  initialLoading,
  error,
  result,
  lastAction,
  isDirty,
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
  applyDraftState,
  initEdit,
  saveDraft,
  publishProgram,
} = useWODForm(initialWODId ? 'edit' : 'create')

const storageKey = computed(() => {
  const workspaceId = session.activeWorkspaceId.value
  if (!workspaceId) {
    return null
  }
  return draftStorageKey(workspaceId, initialWODId ?? wodId.value)
})

const successSummary = computed(() => {
  if (!result.value) {
    return ''
  }
  const stageCount = 'stageCount' in result.value ? result.value.stageCount : result.value.stages.length
  const stageLabels = result.value.stages.map((stage) => `${stage.kind}/${stage.type}`).join(' -> ')
  if (lastAction.value === 'publish') {
    const dateLabel = scheduledDate.value || result.value.scheduledDate || 'the selected date'
    return `${result.value.name} published for ${dateLabel} with ${stageCount} stage(s): ${stageLabels}. Athletes can now see it.`
  }
  return `${result.value.name} saved as draft with ${stageCount} stage(s): ${stageLabels}. Only coaches and owners can see it.`
})

const stageCountLabel = computed(() => `${stages.value.length} stage${stages.value.length === 1 ? '' : 's'}`)
const pageTitle = computed(() => (mode.value === 'edit' ? 'Edit WOD Program' : 'Create WOD Program'))
const pageSubtitle = computed(() =>
  mode.value === 'edit'
    ? 'Update your class plan stages, scoring, and prescriptions.'
    : 'Build your class plan with instructions and free-text prescriptions.',
)
const isEditMode = computed(() => Boolean(initialWODId))
const canUseForm = computed(() => {
  if (isEditMode.value) {
    if (!loadedDetail.value) {
      return false
    }
    return canEditWOD(session.activeWorkspaceRole.value, loadedDetail.value, session.currentUser.value?.id)
  }
  return canCreateWOD(session.activeWorkspaceRole.value)
})
const canPublish = computed(() => canPublishWOD(session.activeWorkspaceRole.value))
const roleLabel = computed(() =>
  session.activeWorkspaceRole.value ? ROLE_LABELS[session.activeWorkspaceRole.value] : 'Member',
)
const recoveryMessage = computed(() =>
  pendingDraft.value
    ? `You have an unfinished program from ${formatDraftAge(pendingDraft.value.savedAt)}. Continue where you left off?`
    : '',
)

function currentDraftSnapshot(): StoredWODDraft {
  return {
    name: name.value,
    description: description.value,
    scheduledDate: scheduledDate.value,
    stages: stages.value,
    mode: mode.value,
    wodId: wodId.value,
    savedAt: new Date().toISOString(),
  }
}

function queueAutosave() {
  if (!canUseForm.value || !storageKey.value || !isDirty.value) {
    return
  }
  if (autosaveTimer) {
    clearTimeout(autosaveTimer)
  }
  autosaveTimer = setTimeout(() => {
    if (storageKey.value) {
      saveStoredDraft(storageKey.value, currentDraftSnapshot())
    }
  }, 2000)
}

function clearDraftStorage() {
  if (storageKey.value) {
    clearStoredDraft(storageKey.value)
  }
}

function maybeOfferRecovery() {
  if (!storageKey.value) {
    return
  }
  const stored = loadStoredDraft(storageKey.value)
  if (!stored || !hasMeaningfulDraft(stored)) {
    return
  }
  if (draftsMatch(stored, currentDraftSnapshot())) {
    return
  }
  pendingDraft.value = stored
  showRecoveryPrompt.value = true
}

function continueDraftRecovery() {
  if (!pendingDraft.value) {
    return
  }
  applyDraftState(pendingDraft.value)
  showRecoveryPrompt.value = false
  pendingDraft.value = null
}

function dismissDraftRecovery() {
  clearDraftStorage()
  showRecoveryPrompt.value = false
  pendingDraft.value = null
}

async function handleSaveDraft() {
  const ok = await saveDraft()
  if (ok) {
    clearDraftStorage()
  }
}

async function handlePublish() {
  const ok = await publishProgram()
  if (ok) {
    clearDraftStorage()
  }
}

watch([name, description, scheduledDate, stages], queueAutosave, { deep: true })

watch(result, (value) => {
  if (value && (lastAction.value === 'draft' || lastAction.value === 'publish')) {
    clearDraftStorage()
  }
})

onMounted(async () => {
  if (initialWODId) {
    loadedDetail.value = await initEdit(initialWODId)
  }
  maybeOfferRecovery()
})
</script>

<template>
  <div class="container stack-lg">
    <header class="page-header">
      <p class="eyebrow">{{ session.activeWorkspace.value?.name ?? 'Active workspace' }}</p>
      <h1 class="page-title">{{ pageTitle }}</h1>
      <p class="page-subtitle">
        {{ pageSubtitle }}
        <span class="inline-role">{{ roleLabel }}</span>
      </p>
    </header>

    <div v-if="showRecoveryPrompt" class="confirmation-panel" role="alertdialog" aria-live="polite">
      <p class="confirmation-panel__text">{{ recoveryMessage }}</p>
      <div class="confirmation-panel__actions">
        <button type="button" class="btn" @click="continueDraftRecovery">Continue editing</button>
        <button type="button" class="btn secondary" @click="dismissDraftRecovery">Start fresh</button>
      </div>
    </div>

    <div v-if="!canUseForm && !initialLoading" class="card empty-state">
      <h2 class="empty-state__title">This workspace is read-only for your role</h2>
      <p class="empty-state__text">
        {{ roleLabel }} access
        {{ isEditMode ? 'cannot edit this saved program.' : 'cannot create programs.' }}
        You can still review saved plans for this gym.
      </p>
      <RouterLink to="/wods" class="empty-state__link">View Programs</RouterLink>
    </div>

    <p v-else-if="initialLoading" class="loading-state">Loading program...</p>

    <div v-else class="create-layout">
      <div class="create-layout__main">
        <form id="wod-create-form" class="card stack-lg" @submit.prevent="handlePublish">
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
            <div class="date-field">
              <label class="date-field__label" for="program-date">Program date</label>
              <input
                id="program-date"
                v-model="scheduledDate"
                class="date-field__input"
                type="date"
              />
              <p class="date-field__hint">The date this class plan is scheduled for.</p>
            </div>
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
            <div class="program-outline__actions-row">
              <button
                type="button"
                class="btn secondary btn-full"
                :disabled="loading"
                @click="handleSaveDraft"
              >
                {{ savingAction === 'draft' ? 'Saving draft...' : 'Save as Draft' }}
              </button>
              <button
                v-if="canPublish"
                type="submit"
                class="btn-full"
                :disabled="loading"
              >
                {{ savingAction === 'publish' ? 'Publishing...' : 'Publish Program' }}
              </button>
            </div>
          </div>
        </form>
      </div>

      <div class="create-layout__aside">
        <ProgramOutlinePanel
          :name="name"
          :description="description"
          :scheduled-date="scheduledDate"
          :stages="stages"
          :loading="loading"
          :saving-action="savingAction"
          :error="error"
          :success-summary="successSummary"
          :can-publish="canPublish"
          @save-draft="handleSaveDraft"
          @publish="handlePublish"
        />
      </div>
    </div>
  </div>
</template>
