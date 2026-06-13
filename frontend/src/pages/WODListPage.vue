<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { listWODs } from '@/features/wod/api/wodApi'
import { STAGE_KIND_BADGE_CLASS, WOD_TYPE_BADGE_CLASS } from '@/features/wod/model/wodTheme'
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
      <h1 class="page-title">Saved Programs</h1>
      <p class="page-subtitle">Your multi-stage workouts, ready to review and reuse.</p>
    </header>

    <p v-if="loading" class="loading-state">Loading programs...</p>
    <div v-else-if="error" class="alert alert--error" role="alert">{{ error }}</div>

    <div v-else-if="items.length === 0" class="card empty-state">
      <h2 class="empty-state__title">No programs yet</h2>
      <p class="empty-state__text">Create your first multi-stage WOD to get started.</p>
      <RouterLink to="/" class="empty-state__link">Create Program</RouterLink>
    </div>

    <div v-else class="wod-list">
      <article v-for="item in items" :key="item.id" class="card card--hover stack">
        <div class="row row--align-center row--between">
          <strong>{{ item.name }}</strong>
          <div class="wod-card__actions">
            <span class="wod-card__meta">{{ item.stageCount }} stage(s)</span>
            <RouterLink :to="`/wods/${item.id}/edit`" class="wod-card__link">Edit</RouterLink>
          </div>
        </div>
        <p class="wod-card__meta">Status: {{ item.status }}</p>
        <div class="badge-row">
          <span v-for="stage in item.stages" :key="`${item.id}-${stage.position}`" class="stage-chip">
            <span class="badge" :class="STAGE_KIND_BADGE_CLASS[stage.kind]">{{ stage.kind }}</span>
            <span class="badge" :class="WOD_TYPE_BADGE_CLASS[stage.type]">{{ stage.type }}</span>
          </span>
        </div>
      </article>
    </div>
  </div>
</template>
