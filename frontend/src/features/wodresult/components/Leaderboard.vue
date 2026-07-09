<script setup lang="ts">
import type { ScoringKind } from '@/features/wod/model/wodTypes'
import type { LeaderboardEntry } from '@/features/wodresult/model/wodResultTypes'
import { decodeScore } from '@/features/wodresult/utils/scoreFormat'

defineProps<{
  entries: LeaderboardEntry[]
  scoringKind: ScoringKind
  stageLabel: string
  loading: boolean
  error: string | null
}>()
</script>

<template>
  <div class="leaderboard">
    <h3 class="score-context__title">Daily Leaderboard - {{ stageLabel }}</h3>

    <p v-if="loading" class="loading-state">Loading leaderboard...</p>
    <div v-else-if="error" class="alert alert--error" role="alert">{{ error }}</div>
    <p v-else-if="entries.length === 0" class="helper-text">No results logged yet. Be the first!</p>

    <ol v-else class="leaderboard__list">
      <li v-for="entry in entries" :key="entry.gymMembershipId" class="card leaderboard__row">
        <span class="leaderboard__rank">{{ entry.rank }}</span>
        <div class="leaderboard__info">
          <div class="leaderboard__name-row">
            <span class="leaderboard__name">{{ entry.displayName }}</span>
            <span class="badge" :class="entry.isRx ? 'badge-rx' : 'badge-scaled'">
              {{ entry.isRx ? 'Rx' : 'Scaled' }}
            </span>
          </div>
          <p v-if="entry.notes" class="leaderboard__notes">{{ entry.notes }}</p>
        </div>
        <span class="leaderboard__score">{{ decodeScore(scoringKind, entry.scoreValue) }}</span>
      </li>
    </ol>
  </div>
</template>
