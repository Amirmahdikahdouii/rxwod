<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { useSession } from '@/features/auth/composables/useSession'
import { useClassSessions } from '@/features/classbooking/composables/useClassSessions'
import type { ClassSessionResponse } from '@/features/classbooking/model/classBookingTypes'
import {
  fetchMembershipDirectory,
  type MembershipDirectory,
} from '@/features/classbooking/utils/membershipDirectory'
import { canViewWOD } from '@/features/workspace/model/workspaceTypes'
import Toast from '@/shared/components/Toast.vue'

const session = useSession()
const {
  sessions,
  loading,
  error,
  actionError,
  actionSessionId,
  classFullMessage,
  fetchSessions,
  book,
  cancel,
} = useClassSessions()

const selectedDate = ref(todayDateString())
const directory = ref<MembershipDirectory>(new Map())
const directoryError = ref<string | null>(null)

const sortedSessions = computed(() =>
  [...sessions.value].sort(
    (a, b) => new Date(a.startTime).getTime() - new Date(b.startTime).getTime(),
  ),
)

function todayDateString(): string {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function formatTimeRange(item: ClassSessionResponse): string {
  const start = new Date(item.startTime).toLocaleTimeString(undefined, {
    hour: 'numeric',
    minute: '2-digit',
  })
  const end = new Date(item.endTime).toLocaleTimeString(undefined, {
    hour: 'numeric',
    minute: '2-digit',
  })
  return `${start} – ${end}`
}

function coachName(coachMembershipId: string): string {
  return directory.value.get(coachMembershipId)?.displayName ?? 'Unknown coach'
}

function isFull(item: ClassSessionResponse): boolean {
  return item.bookedCount >= item.capacity && item.myBookingStatus !== 'BOOKED'
}

function canLinkWod(): boolean {
  return canViewWOD(session.activeWorkspaceRole.value, { status: 'PUBLISHED' })
}

async function loadDirectory() {
  const gymId = session.activeWorkspaceId.value
  if (!gymId) {
    return
  }

  const response = await fetchMembershipDirectory(gymId)
  if (!response.ok) {
    directoryError.value = response.error
    return
  }

  directory.value = response.value
}

async function loadSessions() {
  await fetchSessions(selectedDate.value)
}

async function handleBook(sessionId: string) {
  await book(sessionId)
}

async function handleCancel(sessionId: string) {
  await cancel(sessionId)
}

function dismissClassFullMessage() {
  classFullMessage.value = null
}

onMounted(async () => {
  await Promise.all([loadDirectory(), loadSessions()])
})

watch(
  () => session.activeWorkspaceId.value,
  async () => {
    directoryError.value = null
    await Promise.all([loadDirectory(), loadSessions()])
  },
)

watch(selectedDate, () => {
  void loadSessions()
})
</script>

<template>
  <div class="class-schedule stack-lg">
    <Toast
      v-if="classFullMessage"
      :message="classFullMessage"
      variant="error"
      @dismiss="dismissClassFullMessage"
    />

    <div class="field">
      <label for="schedule-date">Date</label>
      <input id="schedule-date" v-model="selectedDate" type="date" />
    </div>

    <div v-if="directoryError" class="alert alert--error" role="alert">{{ directoryError }}</div>

    <p v-if="loading" class="loading-state">Loading sessions...</p>
    <div v-else-if="error" class="alert alert--error" role="alert">{{ error }}</div>

    <div v-else-if="sortedSessions.length === 0" class="card empty-state">
      <h2 class="empty-state__title">No classes scheduled</h2>
      <p class="empty-state__text">There are no sessions on this date.</p>
    </div>

    <div v-else class="wod-list">
      <article v-for="item in sortedSessions" :key="item.id" class="card card--hover stack">
        <div class="row row--align-center row--between">
          <div class="stack">
            <strong>{{ formatTimeRange(item) }}</strong>
            <p class="wod-card__meta">{{ coachName(item.coachId) }}</p>
            <p class="wod-card__meta">{{ item.bookedCount }}/{{ item.capacity }} Booked</p>
          </div>
          <div class="wod-card__actions">
            <RouterLink
              v-if="item.wodId && canLinkWod()"
              :to="`/wods/${item.wodId}/edit`"
              class="badge badge-open"
            >
              Program
            </RouterLink>
            <span v-else-if="item.wodId" class="badge badge-open">Program</span>
            <button
              v-if="item.myBookingStatus === 'BOOKED'"
              type="button"
              class="btn secondary compact-button"
              :disabled="actionSessionId === item.id"
              @click="handleCancel(item.id)"
            >
              {{ actionSessionId === item.id ? 'Cancelling...' : 'Cancel' }}
            </button>
            <button
              v-else
              type="button"
              class="btn compact-button"
              :disabled="isFull(item) || actionSessionId === item.id"
              @click="handleBook(item.id)"
            >
              {{
                actionSessionId === item.id
                  ? 'Booking...'
                  : isFull(item)
                    ? 'Full'
                    : 'Book'
              }}
            </button>
          </div>
        </div>
        <p v-if="actionError" class="alert alert--error" role="alert">
          {{ actionError }}
        </p>
      </article>
    </div>
  </div>
</template>
