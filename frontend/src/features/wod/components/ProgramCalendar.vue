<script setup lang="ts">
import type { CalendarDaySummary } from '@/features/wod/model/wodTypes'
import {
  buildCalendarRange,
  buildCalendarWeeks,
  buildMonthLabels,
  calendarCellLabel,
} from '@/features/wod/utils/calendarUtils'
import { computed } from 'vue'

const props = defineProps<{
  days: CalendarDaySummary[]
  loading: boolean
  selectedDate: string | null
  showDrafts: boolean
}>()

const emit = defineEmits<{
  selectDate: [date: string | null]
}>()

const dayMap = computed(() => new Map(props.days.map((day) => [day.date, day])))
const range = computed(() => buildCalendarRange(new Date()))
const weeks = computed(() => buildCalendarWeeks(range.value.start, range.value.end, dayMap.value))
const monthLabels = computed(() => buildMonthLabels(weeks.value, range.value.start, range.value.end))
const weekdayLabels = ['Mon', '', 'Wed', '', 'Fri', '', '']

function cellClass(cell: ReturnType<typeof buildCalendarWeeks>[number]['days'][number]) {
  if (!cell.inRange) {
    return 'program-calendar__cell program-calendar__cell--empty'
  }
  const classes = ['program-calendar__cell']
  if (cell.publishedCount > 0) {
    classes.push('program-calendar__cell--has-plan')
  } else {
    classes.push('program-calendar__cell--empty-day')
  }
  if (props.selectedDate === cell.date) {
    classes.push('program-calendar__cell--selected')
  }
  if (props.showDrafts && cell.draftCount > 0) {
    classes.push('program-calendar__cell--has-draft')
  }
  return classes.join(' ')
}

function toggleDate(date: string, inRange: boolean) {
  if (!inRange) {
    return
  }
  emit('selectDate', props.selectedDate === date ? null : date)
}
</script>

<template>
  <section class="card program-calendar" aria-label="Program schedule calendar">
    <div class="program-calendar__header">
      <div>
        <p class="section-title program-calendar__eyebrow">Schedule overview</p>
        <h2 class="program-calendar__title">Program calendar</h2>
      </div>
      <button
        v-if="selectedDate"
        type="button"
        class="btn secondary program-calendar__clear"
        @click="emit('selectDate', null)"
      >
        Show all
      </button>
    </div>

    <p v-if="loading" class="loading-state">Loading calendar...</p>

    <div v-else class="program-calendar__scroll">
      <div class="program-calendar__grid">
        <div class="program-calendar__month-row">
          <span class="program-calendar__weekday-spacer" aria-hidden="true" />
          <div
            class="program-calendar__months"
            :style="{ gridTemplateColumns: `repeat(${weeks.length}, 14px)` }"
          >
            <span
              v-for="month in monthLabels"
              :key="`${month.label}-${month.column}`"
              class="program-calendar__month-label"
              :style="{ gridColumnStart: month.column + 1 }"
            >
              {{ month.label }}
            </span>
          </div>
        </div>

        <div class="program-calendar__body">
          <div class="program-calendar__weekdays" aria-hidden="true">
            <span v-for="(label, index) in weekdayLabels" :key="index" class="program-calendar__weekday">
              {{ label }}
            </span>
          </div>

          <div
            class="program-calendar__weeks"
            :style="{ gridTemplateColumns: `repeat(${weeks.length}, 14px)` }"
          >
            <div v-for="(week, weekIndex) in weeks" :key="weekIndex" class="program-calendar__week">
              <button
                v-for="cell in week.days"
                :key="cell.date"
                type="button"
                :class="cellClass(cell)"
                :aria-label="calendarCellLabel(cell)"
                :disabled="!cell.inRange"
                :title="cell.inRange ? calendarCellLabel(cell) : undefined"
                @click="toggleDate(cell.date, cell.inRange)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="program-calendar__legend">
      <span class="program-calendar__legend-swatch program-calendar__legend-swatch--empty-day" />
      <span class="program-calendar__legend-label">No program</span>
      <span class="program-calendar__legend-swatch program-calendar__legend-swatch--has-plan" />
      <span class="program-calendar__legend-label">Scheduled</span>
      <span v-if="showDrafts" class="program-calendar__legend-draft">Draft days marked with a ring</span>
    </div>
  </section>
</template>
