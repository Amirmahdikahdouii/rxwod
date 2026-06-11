<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { listWODs } from '@/features/wod/api/wodApi'
import { WOD_TYPE_BADGE_CLASS } from '@/features/wod/model/wodTheme'
import type { WODSummary } from '@/features/wod/model/wodTypes'

const items = ref<WODSummary[]>([])
const error = ref<string | null>(null)
const loading = ref(true)

onMounted(async () => {
  const response = await listWODs()
  loading.value = false
  if (!response.ok) {
    error.value = response.error
    return
  }
  items.value = response.value
})
</script>

<template>
  <div class="container stack-lg">
    <header class="page-header">
      <h1 class="page-title">Saved WODs</h1>
      <p class="page-subtitle">Your created workouts, ready to review and reuse.</p>
    </header>

    <p v-if="loading" class="loading-state">Loading workouts...</p>
    <div v-else-if="error" class="alert alert--error" role="alert">{{ error }}</div>

    <div v-else-if="items.length === 0" class="card empty-state">
      <h2 class="empty-state__title">No workouts yet</h2>
      <p class="empty-state__text">Create your first WOD to get started.</p>
      <RouterLink to="/" class="empty-state__link">Create WOD</RouterLink>
    </div>

    <div v-else class="wod-list">
      <article v-for="item in items" :key="item.id" class="card card--hover stack">
        <div class="row row--align-center row--between">
          <strong>{{ item.name }}</strong>
          <span class="badge" :class="WOD_TYPE_BADGE_CLASS[item.type]">{{ item.type }}</span>
        </div>
        <p class="wod-card__meta">Status: {{ item.status }} · Scoring: {{ item.scoringKind }}</p>
      </article>
    </div>
  </div>
</template>
