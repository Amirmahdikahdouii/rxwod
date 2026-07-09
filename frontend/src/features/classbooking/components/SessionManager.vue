<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useSession } from '@/features/auth/composables/useSession'
import { createSession } from '@/features/classbooking/api/classBookingApi'
import { useClassSessions } from '@/features/classbooking/composables/useClassSessions'
import { useSessionRoster } from '@/features/classbooking/composables/useSessionRoster'
import type { ClassSessionResponse } from '@/features/classbooking/model/classBookingTypes'
import {
  fetchMembershipDirectory,
  type MembershipDirectory,
} from '@/features/classbooking/utils/membershipDirectory'
import { listMembers } from '@/features/workspace/api/workspaceApi'
import type { MemberResponse } from '@/features/workspace/model/workspaceTypes'
import { listWODs } from '@/features/wod/api/wodApi'
import type { WODSummary } from '@/features/wod/model/wodTypes'
import BaseNumberInput from '@/shared/components/BaseNumberInput.vue'
import BaseSelect from '@/shared/components/BaseSelect.vue'

const session = useSession()
const {
  sessions,
  loading,
  error,
  fetchSessions,
} = useClassSessions()
const { roster, loading: rosterLoading, overbooking, overbookError, fetchRoster, overbook } =
  useSessionRoster()

const selectedDate = ref(todayDateString())
const directory = ref<MembershipDirectory>(new Map())
const directoryError = ref<string | null>(null)
const athleteMembers = ref<MemberResponse[]>([])
const publishedWods = ref<WODSummary[]>([])

const expandedSessionId = ref<string | null>(null)
const selectedAthleteUserId = ref('')

const startTime = ref('')
const endTime = ref('')
const capacity = ref<number | undefined>(20)
const selectedWodId = ref('')
const creating = ref(false)
const createError = ref<string | null>(null)
const createSuccess = ref<string | null>(null)

const sortedSessions = computed(() =>
  [...sessions.value].sort(
    (a, b) => new Date(a.startTime).getTime() - new Date(b.startTime).getTime(),
  ),
)

const wodOptions = computed(() => [
  { value: '', label: 'No WOD' },
  ...publishedWods.value.map((wod) => ({ value: wod.id, label: wod.name })),
])

const bookedMembershipIds = computed(() =>
  roster.value.filter((booking) => booking.status !== 'CANCELLED').map((booking) => booking.gymMembershipId),
)

const overbookAthleteOptions = computed(() =>
  athleteMembers.value
    .filter((member) => !bookedMembershipIds.value.includes(member.membershipId))
    .map((member) => ({ value: member.userId, label: member.displayName })),
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

function memberName(membershipId: string): string {
  return directory.value.get(membershipId)?.displayName ?? 'Unknown member'
}

function bookingStatusClass(status: string): string {
  if (status === 'BOOKED') {
    return 'status-pill status-pill--published'
  }
  if (status === 'ATTENDED') {
    return 'status-pill'
  }
  return 'status-pill status-pill--archived'
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

async function loadAthleteMembers() {
  const gymId = session.activeWorkspaceId.value
  if (!gymId) {
    return
  }

  const members: MemberResponse[] = []
  let page = 1

  while (true) {
    const response = await listMembers(gymId, { page, limit: 100 })
    if (!response.ok) {
      directoryError.value = response.error
      return
    }

    members.push(...response.value.data.filter((member) => member.role === 'athlete'))
    if (page >= response.value.meta.totalPages) {
      break
    }
    page++
  }

  athleteMembers.value = members
}

async function loadPublishedWods() {
  const response = await listWODs({ status: 'PUBLISHED', page: 1, limit: 100 })
  if (!response.ok) {
    createError.value = response.error
    return
  }

  publishedWods.value = response.value.data
}

async function loadSessions() {
  await fetchSessions(selectedDate.value)
}

async function handleCreate() {
  if (!startTime.value || !endTime.value || capacity.value === undefined) {
    return
  }

  creating.value = true
  createError.value = null
  createSuccess.value = null

  const payload = {
    startTime: new Date(startTime.value).toISOString(),
    endTime: new Date(endTime.value).toISOString(),
    capacity: capacity.value,
    ...(selectedWodId.value ? { wodId: selectedWodId.value } : {}),
  }

  const response = await createSession(payload)
  creating.value = false

  if (!response.ok) {
    createError.value = response.error
    return
  }

  createSuccess.value =
    response.value.autoBookedCount > 0
      ? `Session created. ${response.value.autoBookedCount} athlete(s) auto-booked from default preferences.`
      : 'Session created successfully.'

  startTime.value = ''
  endTime.value = ''
  selectedWodId.value = ''
  await loadSessions()
}

async function toggleSession(sessionId: string) {
  if (expandedSessionId.value === sessionId) {
    expandedSessionId.value = null
    selectedAthleteUserId.value = ''
    return
  }

  expandedSessionId.value = sessionId
  selectedAthleteUserId.value = ''
  await fetchRoster(sessionId)
}

async function handleOverbook(sessionId: string) {
  if (!selectedAthleteUserId.value) {
    return
  }

  const response = await overbook(sessionId, selectedAthleteUserId.value)
  if (response.ok) {
    selectedAthleteUserId.value = ''
    await loadSessions()
  }
}

onMounted(async () => {
  await Promise.all([loadDirectory(), loadAthleteMembers(), loadPublishedWods(), loadSessions()])
})

watch(
  () => session.activeWorkspaceId.value,
  async () => {
    directoryError.value = null
    expandedSessionId.value = null
    await Promise.all([loadDirectory(), loadAthleteMembers(), loadPublishedWods(), loadSessions()])
  },
)

watch(selectedDate, () => {
  expandedSessionId.value = null
  void loadSessions()
})
</script>

<template>
  <div class="session-manager stack-lg">
    <article class="card stack session-manager__create">
      <h2 class="section-title">Create Session</h2>
      <form class="stack" @submit.prevent="handleCreate">
        <div class="field">
          <label for="session-start">Start</label>
          <input id="session-start" v-model="startTime" type="datetime-local" required />
        </div>
        <div class="field">
          <label for="session-end">End</label>
          <input id="session-end" v-model="endTime" type="datetime-local" required />
        </div>
        <BaseNumberInput v-model="capacity" label="Capacity" :min="1" />
        <BaseSelect v-model="selectedWodId" label="WOD (optional)" :options="wodOptions" />
        <div v-if="createError" class="alert alert--error" role="alert">{{ createError }}</div>
        <div v-if="createSuccess" class="alert alert--success" role="status">{{ createSuccess }}</div>
        <button type="submit" class="btn-full" :disabled="creating">
          {{ creating ? 'Creating...' : 'Create session' }}
        </button>
      </form>
    </article>

    <div class="field session-manager__filter">
      <label for="manage-schedule-date">Date</label>
      <input id="manage-schedule-date" v-model="selectedDate" type="date" />
    </div>

    <div v-if="directoryError" class="alert alert--error" role="alert">{{ directoryError }}</div>

    <p v-if="loading" class="loading-state">Loading sessions...</p>
    <div v-else-if="error" class="alert alert--error" role="alert">{{ error }}</div>

    <div v-else-if="sortedSessions.length === 0" class="card empty-state">
      <h2 class="empty-state__title">No sessions on this date</h2>
      <p class="empty-state__text">Create a session above or choose another date.</p>
    </div>

    <div v-else class="wod-list">
      <article v-for="item in sortedSessions" :key="item.id" class="card stack">
        <button
          type="button"
          class="session-manager__session-toggle"
          :aria-expanded="expandedSessionId === item.id"
          @click="toggleSession(item.id)"
        >
          <div class="row row--align-center row--between">
            <div class="stack session-manager__session-summary">
              <strong class="session-manager__time">{{ formatTimeRange(item) }}</strong>
              <p class="wod-card__meta">{{ item.bookedCount }}/{{ item.capacity }} Booked</p>
            </div>
            <span class="wod-card__link">{{ expandedSessionId === item.id ? 'Hide roster' : 'View roster' }}</span>
          </div>
        </button>

        <div v-if="expandedSessionId === item.id" class="session-manager__roster stack">
          <p v-if="rosterLoading" class="loading-state">Loading roster...</p>
          <div v-else-if="overbookError && roster.length === 0" class="alert alert--error" role="alert">
            {{ overbookError }}
          </div>
          <ul v-else-if="roster.length > 0" class="session-manager__roster-list stack">
            <li
              v-for="booking in roster"
              :key="booking.id"
              class="session-manager__roster-item"
            >
              <span class="session-manager__roster-name">{{ memberName(booking.gymMembershipId) }}</span>
              <span :class="bookingStatusClass(booking.status)">{{ booking.status.toLowerCase() }}</span>
            </li>
          </ul>
          <p v-else class="helper-text">No bookings yet.</p>

          <div class="session-manager__overbook stack">
            <h3 class="section-title section-title--small">Overbook athlete</h3>
            <BaseSelect
              v-model="selectedAthleteUserId"
              label="Athlete"
              :options="overbookAthleteOptions"
              :disabled="overbookAthleteOptions.length === 0"
            />
            <div v-if="overbookError" class="alert alert--error" role="alert">{{ overbookError }}</div>
            <button
              type="button"
              class="btn secondary"
              :disabled="!selectedAthleteUserId || overbooking || overbookAthleteOptions.length === 0"
              @click="handleOverbook(item.id)"
            >
              {{ overbooking ? 'Overbooking...' : 'Overbook' }}
            </button>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>
