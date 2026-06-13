<script setup lang="ts">
import { STAGE_KIND_BADGE_CLASS } from '@/features/wod/model/wodTheme'
import type { StageKind } from '@/features/wod/model/wodTypes'
import { STAGE_KINDS } from '@/features/wod/model/wodTypes'

const emit = defineEmits<{
  add: [kind: StageKind]
}>()

const stageDescriptions: Record<StageKind, string> = {
  WARMUP: 'Prep movement patterns',
  STRENGTH: 'Build heavy work sets',
  CORE: 'Add trunk accessory work',
  METCON: 'Scored conditioning',
  COOLDOWN: 'Finish with recovery',
}

function stageKindLabel(kind: StageKind) {
  return kind.charAt(0) + kind.slice(1).toLowerCase()
}
</script>

<template>
  <section class="add-stage-panel" aria-label="Add a stage">
    <div class="add-stage-panel__header">
      <h3 class="add-stage-panel__title">Add a stage</h3>
      <p class="add-stage-panel__text">Choose a block type to append to your program.</p>
    </div>

    <div class="add-stage-panel__grid">
      <button
        v-for="kind in STAGE_KINDS"
        :key="kind"
        type="button"
        class="add-stage-card"
        @click="emit('add', kind)"
      >
        <span class="badge" :class="STAGE_KIND_BADGE_CLASS[kind]">{{ stageKindLabel(kind) }}</span>
        <span class="add-stage-card__label">{{ stageKindLabel(kind) }}</span>
        <span class="add-stage-card__description">{{ stageDescriptions[kind] }}</span>
      </button>
    </div>
  </section>
</template>
