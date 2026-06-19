<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useSession } from '@/features/auth/composables/useSession'
import ProgramCalendar from '@/features/wod/components/ProgramCalendar.vue'
import { getWODCalendar, listWODs, publishWOD } from '@/features/wod/api/wodApi'
import { buildCalendarRange } from '@/features/wod/utils/calendarUtils'
import { STAGE_KIND_BADGE_CLASS, WOD_TYPE_BADGE_CLASS } from '@/features/wod/model/wodTheme'
import type { CalendarDaySummary, WODSummary } from '@/features/wod/model/wodTypes'
import {
  canCreateWOD,
  canEditWOD,
  canPublishWOD,
  ROLE_LABELS,
} from '@/features/workspace/model/workspaceTypes'

const session = useSession()
const items = ref<WODSummary[]>([])
const calendarDays = ref<CalendarDaySummary[]>([])
const error = ref<string | null>(null)
const calendarError = ref<string | null>(null)
const loading = ref(true)
const calendarLoading = ref(true)
const selectedDate = ref<string | null>(null)
const publishingId = ref<string | null>(null)

const showDrafts = computed(
  () => session.activeWorkspaceRole.value === 'owner' || session.activeWorkspaceRole.value === 'coach',
)

const filteredItems = computed(() => {
  if (!selectedDate.value) {
    return items.value
  }
  return items.value.filter((item) => item.scheduledDate === selectedDate.value)
})

function formatScheduledDate(value?: string) {
  if (!value) {
    return 'No date scheduled'
  }
  const parsed = new Date(`${value}T00:00:00`)
  return parsed.toLocaleDateString(undefined, {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

function statusClass(status: WODSummary['status']) {
  return status === 'PUBLISHED' ? 'status-pill status-pill--published' : 'status-pill status-pill--draft'
}

function canEditItem(item: WODSummary) {
  return canEditWOD(session.activeWorkspaceRole.value, item, session.currentUser.value?.id)
}

async function loadPrograms() {
  loading.value = true
  const response = await listWODs()
  loading.value = false
  if (!response.ok) {
    error.value = response.error
    return
  }
  items.value = response.value
}

async function loadCalendar() {
  calendarLoading.value = true
  const range = buildCalendarRange(new Date())
  const response = await getWODCalendar(range.from, range.to)
  calendarLoading.value = false
  if (!response.ok) {
    calendarError.value = response.error
    return
  }
  calendarDays.value = response.value
}

async function handlePublish(item: WODSummary) {
  publishingId.value = item.id
  const response = await publishWOD(item.id)
  publishingId.value = null
  if (!response.ok) {
    error.value = response.error
    return
  }
  await Promise.all([loadPrograms(), loadCalendar()])
}

onMounted(async () => {
  await Promise.all([loadPrograms(), loadCalendar()])
})
</script>

<template>
  <div class="container stack-lg">
    <header class="page-header">
      <p class="eyebrow">Workspace Plans</p>
      <h1 class="page-title">Saved Programs</h1>
      <p class="page-subtitle">
        Programs for {{ session.activeWorkspace.value?.name ?? 'your active gym' }}.
        <span v-if="session.activeWorkspaceRole.value" class="inline-role">
          {{ ROLE_LABELS[session.activeWorkspaceRole.value] }}
        </span>
      </p>
    </header>

    <div v-if="calendarError" class="alert alert--error" role="alert">{{ calendarError }}</div>
    <ProgramCalendar
      :days="calendarDays"
      :loading="calendarLoading"
      :selected-date="selectedDate"
      :show-drafts="showDrafts"
      @select-date="selectedDate = $event"
    />

    <p v-if="loading" class="loading-state">Loading programs...</p>
    <div v-else-if="error" class="alert alert--error" role="alert">{{ error }}</div>

    <div v-else-if="items.length === 0" class="card empty-state">
      <h2 class="empty-state__title">No programs yet</h2>
      <p class="empty-state__text">Create your first multi-stage WOD to get started.</p>
      <RouterLink v-if="canCreateWOD(session.activeWorkspaceRole.value)" to="/" class="empty-state__link">
        Create Program
      </RouterLink>
    </div>

    <div v-else-if="filteredItems.length === 0" class="card empty-state">
      <h2 class="empty-state__title">No programs on this date</h2>
      <p class="empty-state__text">Try another day on the calendar or show all programs.</p>
      <button type="button" class="empty-state__link btn secondary" @click="selectedDate = null">Show all</button>
    </div>

    <div v-else class="wod-list">
      <article v-for="item in filteredItems" :key="item.id" class="card card--hover stack">
        <div class="row row--align-center row--between">
          <div class="stack">
            <strong>{{ item.name }}</strong>
            <p class="wod-card__date">{{ formatScheduledDate(item.scheduledDate) }}</p>
          </div>
          <div class="wod-card__actions">
            <span :class="statusClass(item.status)">{{ item.status.toLowerCase() }}</span>
            <RouterLink
              v-if="canEditItem(item)"
              :to="`/wods/${item.id}/edit`"
              class="wod-card__link"
            >
              Edit
            </RouterLink>
            <button
              v-if="canEditItem(item) && item.status === 'DRAFT' && canPublishWOD(session.activeWorkspaceRole.value)"
              type="button"
              class="wod-card__link btn secondary compact-button"
              :disabled="publishingId === item.id"
              @click="handlePublish(item)"
            >
              {{ publishingId === item.id ? 'Publishing...' : 'Publish' }}
            </button>
            <span v-if="!canEditItem(item)" class="wod-card__meta">Read only</span>
          </div>
        </div>
        <p class="wod-card__meta">{{ item.stageCount }} stage(s)</p>
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
