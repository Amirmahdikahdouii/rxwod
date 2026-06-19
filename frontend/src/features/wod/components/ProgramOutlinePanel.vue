<script setup lang="ts">
import { STAGE_KIND_BADGE_CLASS, WOD_TYPE_BADGE_CLASS } from '@/features/wod/model/wodTheme'
import type { StageFormState, StageKind } from '@/features/wod/model/wodTypes'
import { stageFormat } from '@/features/wod/model/wodTypes'
import { computed } from 'vue'

const props = defineProps<{
  name: string
  description: string
  scheduledDate: string
  stages: StageFormState[]
  loading: boolean
  savingAction: 'draft' | 'publish' | null
  error: string | null
  successSummary: string
  canPublish: boolean
}>()

const emit = defineEmits<{
  saveDraft: []
  publish: []
}>()

const scheduledDateLabel = computed(() => {
  if (!props.scheduledDate) {
    return 'No program date selected yet.'
  }
  const parsed = new Date(`${props.scheduledDate}T00:00:00`)
  return parsed.toLocaleDateString(undefined, {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
})

function stageKindLabel(kind: StageKind) {
  return kind.charAt(0) + kind.slice(1).toLowerCase()
}

function scrollToStage(index: number) {
  document.getElementById(`stage-${index}`)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}
</script>

<template>
  <aside class="card program-outline stack">
    <div>
      <p class="section-title program-outline__eyebrow">Program outline</p>
      <h2 class="program-outline__title">{{ props.name.trim() || 'Untitled program' }}</h2>
      <p class="program-outline__description">
        {{ props.description.trim() || 'Add a short class plan summary.' }}
      </p>
      <p class="program-outline__date">{{ scheduledDateLabel }}</p>
    </div>

    <div class="program-outline__list" aria-label="Program stages">
      <button
        v-for="(stage, index) in props.stages"
        :key="index"
        type="button"
        class="program-outline__item"
        @click="scrollToStage(index)"
      >
        <span class="program-outline__number">{{ index + 1 }}</span>
        <span>
          <span class="badge-row">
            <span class="badge" :class="STAGE_KIND_BADGE_CLASS[stage.kind]">{{ stageKindLabel(stage.kind) }}</span>
            <span class="badge" :class="WOD_TYPE_BADGE_CLASS[stage.type]">{{ stage.type }}</span>
          </span>
          <span class="program-outline__meta">
            {{ stage.movements.length }} item(s) - {{ stageFormat(stage.type).toLowerCase() }}
          </span>
        </span>
      </button>
    </div>

    <div class="program-outline__actions program-outline__desktop-actions stack">
      <div v-if="props.error" class="alert alert--error" role="alert">{{ props.error }}</div>
      <div v-if="props.successSummary" class="alert alert--success" role="status">{{ props.successSummary }}</div>
      <div class="program-outline__actions-row">
        <button
          type="button"
          class="btn secondary btn-full"
          :disabled="props.loading"
          @click="emit('saveDraft')"
        >
          {{ props.savingAction === 'draft' ? 'Saving draft...' : 'Save as Draft' }}
        </button>
        <button
          v-if="props.canPublish"
          type="button"
          class="btn-full"
          :disabled="props.loading"
          @click="emit('publish')"
        >
          {{ props.savingAction === 'publish' ? 'Publishing...' : 'Publish Program' }}
        </button>
      </div>
    </div>
  </aside>
</template>
