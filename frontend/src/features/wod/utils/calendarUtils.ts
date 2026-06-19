import type { CalendarDaySummary } from '@/features/wod/model/wodTypes'

export interface CalendarCell {
  date: string
  inRange: boolean
  publishedCount: number
  draftCount: number
  plans: CalendarDaySummary['plans']
}

export interface CalendarWeek {
  days: CalendarCell[]
}

export interface CalendarMonthLabel {
  label: string
  column: number
}

export function formatDateKey(date: Date): string {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

export function heatmapLevel(publishedCount: number): number {
  if (publishedCount <= 0) {
    return 0
  }
  if (publishedCount === 1) {
    return 1
  }
  if (publishedCount === 2) {
    return 2
  }
  if (publishedCount === 3) {
    return 3
  }
  return 4
}

export function buildCalendarRange(endDate: Date, months = 12): { from: string; to: string; start: Date; end: Date } {
  const end = new Date(endDate.getFullYear(), endDate.getMonth(), endDate.getDate())
  const start = new Date(end.getFullYear(), end.getMonth() - (months - 1), 1)
  return {
    from: formatDateKey(start),
    to: formatDateKey(end),
    start,
    end,
  }
}

function startOfWeekMonday(date: Date): Date {
  const copy = new Date(date.getFullYear(), date.getMonth(), date.getDate())
  const day = copy.getDay()
  const diff = day === 0 ? -6 : 1 - day
  copy.setDate(copy.getDate() + diff)
  return copy
}

export function buildCalendarWeeks(
  start: Date,
  end: Date,
  dayMap: Map<string, CalendarDaySummary>,
): CalendarWeek[] {
  const gridStart = startOfWeekMonday(start)
  const gridEnd = startOfWeekMonday(end)
  gridEnd.setDate(gridEnd.getDate() + 6)

  const weeks: CalendarWeek[] = []
  const cursor = new Date(gridStart)

  while (cursor <= gridEnd) {
    const week: CalendarWeek = { days: [] }
    for (let i = 0; i < 7; i += 1) {
      const dateKey = formatDateKey(cursor)
      const inRange = cursor >= start && cursor <= end
      const summary = dayMap.get(dateKey)
      week.days.push({
        date: dateKey,
        inRange,
        publishedCount: summary?.publishedCount ?? 0,
        draftCount: summary?.draftCount ?? 0,
        plans: summary?.plans ?? [],
      })
      cursor.setDate(cursor.getDate() + 1)
    }
    weeks.push(week)
  }

  return weeks
}

export function buildMonthLabels(weeks: CalendarWeek[], start: Date, end: Date): CalendarMonthLabel[] {
  const labels: CalendarMonthLabel[] = []
  let lastMonth = -1

  weeks.forEach((week, column) => {
    const firstInRange = week.days.find((day) => {
      const parsed = new Date(`${day.date}T00:00:00`)
      return parsed >= start && parsed <= end
    })
    if (!firstInRange) {
      return
    }
    const parsed = new Date(`${firstInRange.date}T00:00:00`)
    if (parsed.getMonth() === lastMonth) {
      return
    }
    lastMonth = parsed.getMonth()
    labels.push({
      label: parsed.toLocaleString(undefined, { month: 'short' }),
      column,
    })
  })

  return labels
}

export function calendarCellLabel(cell: CalendarCell): string {
  const parsed = new Date(`${cell.date}T00:00:00`)
  const formatted = parsed.toLocaleDateString(undefined, {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  })
  const total = cell.publishedCount + cell.draftCount
  if (total === 0) {
    return `${formatted}. No programs scheduled.`
  }
  const planNames = cell.plans.map((plan) => `${plan.name} (${plan.status.toLowerCase()})`).join(', ')
  return `${formatted}. ${total} program${total === 1 ? '' : 's'}: ${planNames}.`
}
