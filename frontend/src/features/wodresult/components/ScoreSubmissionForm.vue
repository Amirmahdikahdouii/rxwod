<script setup lang="ts">
import { computed, ref } from 'vue'
import type { ScoringKind } from '@/features/wod/model/wodTypes'
import { useWODResults } from '@/features/wodresult/composables/useWODResults'
import { encodeRoundsReps, encodeTimeToComplete, scoringKindLabel } from '@/features/wodresult/utils/scoreFormat'
import BaseNumberInput from '@/shared/components/BaseNumberInput.vue'
import BaseTextarea from '@/shared/components/BaseTextarea.vue'
import Toast from '@/shared/components/Toast.vue'

const props = defineProps<{
  wodId: string
  scoringKind: ScoringKind
  stageLabel: string
}>()

const emit = defineEmits<{
  submitted: []
}>()

const { submitting, submitError, submit } = useWODResults()

const expanded = ref(true)
const showToast = ref(false)

const minutes = ref<number>()
const seconds = ref<number>()
const rounds = ref<number>()
const reps = ref<number>()
const isRx = ref(true)
const notes = ref('')

const scoreValue = computed(() => {
  switch (props.scoringKind) {
    case 'TIME_TO_COMPLETE':
      return encodeTimeToComplete({ minutes: minutes.value ?? 0, seconds: seconds.value ?? 0 })
    case 'ROUNDS_REPS':
      return encodeRoundsReps({ rounds: rounds.value ?? 0, reps: reps.value ?? 0 })
    case 'TOTAL_REPS':
      return reps.value ?? 0
    case 'NONE':
    default:
      return 0
  }
})

const canSubmit = computed(() => {
  switch (props.scoringKind) {
    case 'TIME_TO_COMPLETE':
      return (minutes.value ?? 0) > 0 || (seconds.value ?? 0) > 0
    case 'ROUNDS_REPS':
      return (rounds.value ?? 0) > 0 || (reps.value ?? 0) > 0
    case 'TOTAL_REPS':
      return (reps.value ?? 0) > 0
    case 'NONE':
    default:
      return true
  }
})

async function handleSubmit() {
  const response = await submit(props.wodId, {
    scoreValue: scoreValue.value,
    isRx: isRx.value,
    notes: notes.value.trim(),
  })

  if (!response.ok) {
    return
  }

  expanded.value = false
  showToast.value = true
  emit('submitted')
}

function editResult() {
  expanded.value = true
  showToast.value = false
}
</script>

<template>
  <div class="score-submission stack">
    <Toast
      v-if="showToast"
      message="Result saved successfully!"
      variant="success"
      @dismiss="showToast = false"
    />

    <form v-if="expanded" class="stack" @submit.prevent="handleSubmit">
      <div class="score-context">
        <h3 class="score-context__title">Log Your Score for {{ stageLabel }}</h3>
        <p class="score-context__subtitle">Scoring Type: {{ scoringKindLabel(scoringKind) }}</p>
      </div>

      <div v-if="scoringKind === 'TIME_TO_COMPLETE'" class="row">
        <BaseNumberInput v-model="minutes" label="Min" placeholder="MM" :min="0" />
        <BaseNumberInput v-model="seconds" label="Sec" placeholder="SS" :min="0" />
      </div>

      <div v-else-if="scoringKind === 'ROUNDS_REPS'" class="row">
        <BaseNumberInput v-model="rounds" label="Rounds" placeholder="e.g. 5" :min="0" />
        <BaseNumberInput v-model="reps" label="Reps" placeholder="e.g. 12" :min="0" />
      </div>

      <div v-else-if="scoringKind === 'TOTAL_REPS'" class="row">
        <BaseNumberInput v-model="reps" label="Total Reps" placeholder="Enter your total reps" :min="0" />
      </div>

      <p v-else class="helper-text">
        This WOD has no structured scoring. You can still log your attendance and notes below.
      </p>

      <div class="field">
        <label>Result type</label>
        <div class="segmented-control" role="group" aria-label="Rx or Scaled">
          <button
            type="button"
            class="segmented-control__option"
            :class="{ 'segmented-control__option--active': isRx }"
            @click="isRx = true"
          >
            Rx
          </button>
          <button
            type="button"
            class="segmented-control__option"
            :class="{ 'segmented-control__option--active': !isRx }"
            @click="isRx = false"
          >
            Scaled
          </button>
        </div>
      </div>

      <BaseTextarea
        v-model="notes"
        label="Notes (optional)"
        placeholder="How did it feel? Any modifications?"
        :rows="2"
      />

      <div v-if="submitError" class="alert alert--error" role="alert">{{ submitError }}</div>

      <button type="submit" class="btn-full" :disabled="submitting || !canSubmit">
        {{ submitting ? 'Saving...' : 'Log Result' }}
      </button>
    </form>

    <div v-else class="confirmation-panel">
      <p class="confirmation-panel__text">Your result has been logged.</p>
      <div class="confirmation-panel__actions">
        <button type="button" class="secondary" @click="editResult">Edit result</button>
      </div>
    </div>
  </div>
</template>
