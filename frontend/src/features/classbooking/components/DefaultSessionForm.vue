<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useBookingPreferences } from '@/features/classbooking/composables/useBookingPreferences'
import { DAY_OF_WEEK_LABELS } from '@/features/classbooking/model/classBookingTypes'
import BaseSelect from '@/shared/components/BaseSelect.vue'

const { defaultSessions, loading, submitting, error, fetchMine, addDefaultSession, removeDefaultSession } =
  useBookingPreferences()

const dayOfWeek = ref('1')
const timeSlot = ref('')
const removingId = ref<string | null>(null)

const dayOptions = computed(() =>
  DAY_OF_WEEK_LABELS.map((label, index) => ({
    value: String(index),
    label,
  })),
)

function normalizeTimeSlot(value: string): string {
  return value.slice(0, 5)
}

async function handleSubmit() {
  if (!timeSlot.value) {
    return
  }

  const response = await addDefaultSession({
    dayOfWeek: Number(dayOfWeek.value),
    timeSlot: normalizeTimeSlot(timeSlot.value),
  })

  if (response.ok) {
    timeSlot.value = ''
  }
}

async function handleRemove(id: string) {
  removingId.value = id
  await removeDefaultSession(id)
  removingId.value = null
}

onMounted(() => {
  void fetchMine()
})
</script>

<template>
  <article class="card stack">
    <h2 class="section-title">Default Class Times</h2>
    <p class="page-subtitle page-subtitle--flush">
      Set your preferred class days and times for automatic booking when sessions are created.
    </p>

    <div v-if="loading" class="helper-text">Loading preferences...</div>
    <div v-else-if="error && defaultSessions.length === 0" class="alert alert--error" role="alert">
      {{ error }}
    </div>

    <ul v-if="defaultSessions.length > 0" class="default-session-list stack">
      <li
        v-for="session in defaultSessions"
        :key="session.id"
        class="default-session-list__item"
      >
        <span>{{ DAY_OF_WEEK_LABELS[session.dayOfWeek] }} at {{ session.timeSlot }}</span>
        <button
          type="button"
          class="secondary"
          :disabled="removingId === session.id"
          @click="handleRemove(session.id)"
        >
          {{ removingId === session.id ? 'Removing...' : 'Remove' }}
        </button>
      </li>
    </ul>

    <form class="stack" @submit.prevent="handleSubmit">
      <BaseSelect v-model="dayOfWeek" label="Day of week" :options="dayOptions" />
      <div class="field">
        <label for="default-session-time">Time slot</label>
        <input
          id="default-session-time"
          v-model="timeSlot"
          type="time"
          required
        />
      </div>
      <div v-if="error" class="alert alert--error" role="alert">{{ error }}</div>
      <button type="submit" class="btn-full" :disabled="submitting || !timeSlot">
        {{ submitting ? 'Saving...' : 'Add default time' }}
      </button>
    </form>
  </article>
</template>

<style scoped>
.default-session-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.default-session-list__item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}
</style>
