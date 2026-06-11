<script setup lang="ts">
import type { ScoringKind, WODType } from '@/features/wod/model/wodTypes'
import { SCORING_BY_TYPE } from '@/features/wod/model/wodTypes'
import { WOD_TYPE_BADGE_CLASS } from '@/features/wod/model/wodTheme'

defineProps<{
  type: WODType
}>()

function label(kind: ScoringKind): string {
  switch (kind) {
    case 'ROUNDS_REPS':
      return 'Scored by rounds + reps'
    case 'TIME_TO_COMPLETE':
      return 'Scored by completion time'
    case 'TOTAL_REPS':
      return 'Scored by total reps'
    case 'NONE':
      return 'No structured scoring — prescription-based stage'
  }
}
</script>

<template>
  <div class="score-preview">
    <span class="badge" :class="WOD_TYPE_BADGE_CLASS[type]">{{ type }}</span>
    <p class="score-preview__text">{{ label(SCORING_BY_TYPE[type]) }}</p>
  </div>
</template>
