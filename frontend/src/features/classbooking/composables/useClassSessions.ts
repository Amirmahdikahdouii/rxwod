import { ref } from 'vue'
import {
  bookSession,
  cancelBooking,
  listSessions,
} from '@/features/classbooking/api/classBookingApi'
import type { ClassBookingResponse, ClassSessionResponse } from '@/features/classbooking/model/classBookingTypes'
import { err, ok, type Result } from '@/shared/utils/result'

const CLASS_FULL_ERROR = 'class session has reached capacity'

function dayBounds(date: string): { from: string; to: string } {
  const start = new Date(`${date}T00:00:00`)
  const end = new Date(`${date}T23:59:59.999`)
  return { from: start.toISOString(), to: end.toISOString() }
}

export function useClassSessions() {
  const sessions = ref<ClassSessionResponse[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const actionError = ref<string | null>(null)
  const actionSessionId = ref<string | null>(null)
  const classFullMessage = ref<string | null>(null)

  function patchSession(
    sessionId: string,
    updater: (session: ClassSessionResponse) => ClassSessionResponse,
  ) {
    const index = sessions.value.findIndex((session) => session.id === sessionId)
    if (index === -1) {
      return
    }
    sessions.value[index] = updater(sessions.value[index])
  }

  async function fetchSessions(date: string): Promise<Result<ClassSessionResponse[]>> {
    loading.value = true
    error.value = null

    const { from, to } = dayBounds(date)
    const response = await listSessions(from, to)
    loading.value = false

    if (!response.ok) {
      error.value = response.error
      return err(response.error)
    }

    sessions.value = response.value
    return ok(response.value)
  }

  async function book(sessionId: string): Promise<Result<ClassBookingResponse>> {
    actionSessionId.value = sessionId
    actionError.value = null
    classFullMessage.value = null

    const response = await bookSession(sessionId)
    actionSessionId.value = null

    if (!response.ok) {
      if (response.error === CLASS_FULL_ERROR) {
        classFullMessage.value = response.error
      } else {
        actionError.value = response.error
      }
      return err(response.error)
    }

    patchSession(sessionId, (session) => {
      const wasBooked = session.myBookingStatus === 'BOOKED'
      return {
        ...session,
        myBookingStatus: 'BOOKED',
        bookedCount: wasBooked ? session.bookedCount : session.bookedCount + 1,
      }
    })

    return ok(response.value)
  }

  async function cancel(sessionId: string): Promise<Result<ClassBookingResponse>> {
    actionSessionId.value = sessionId
    actionError.value = null
    classFullMessage.value = null

    const response = await cancelBooking(sessionId)
    actionSessionId.value = null

    if (!response.ok) {
      actionError.value = response.error
      return err(response.error)
    }

    patchSession(sessionId, (session) => {
      const wasBooked = session.myBookingStatus === 'BOOKED'
      return {
        ...session,
        myBookingStatus: 'CANCELLED',
        bookedCount: wasBooked ? Math.max(0, session.bookedCount - 1) : session.bookedCount,
      }
    })

    return ok(response.value)
  }

  return {
    sessions,
    loading,
    error,
    actionError,
    actionSessionId,
    classFullMessage,
    fetchSessions,
    book,
    cancel,
  }
}
