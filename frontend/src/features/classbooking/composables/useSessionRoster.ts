import { ref } from 'vue'
import { getRoster, overbookSession } from '@/features/classbooking/api/classBookingApi'
import type { ClassBookingResponse } from '@/features/classbooking/model/classBookingTypes'
import { err, ok, type Result } from '@/shared/utils/result'

export function useSessionRoster() {
  const roster = ref<ClassBookingResponse[]>([])
  const loading = ref(false)
  const overbooking = ref(false)
  const overbookError = ref<string | null>(null)

  async function fetchRoster(sessionId: string): Promise<Result<ClassBookingResponse[]>> {
    loading.value = true
    overbookError.value = null

    const response = await getRoster(sessionId)
    loading.value = false

    if (!response.ok) {
      overbookError.value = response.error
      return err(response.error)
    }

    roster.value = response.value
    return ok(response.value)
  }

  async function overbook(
    sessionId: string,
    athleteUserId: string,
  ): Promise<Result<ClassBookingResponse>> {
    overbooking.value = true
    overbookError.value = null

    const response = await overbookSession(sessionId, athleteUserId)
    overbooking.value = false

    if (!response.ok) {
      overbookError.value = response.error
      return err(response.error)
    }

    await fetchRoster(sessionId)
    return ok(response.value)
  }

  return {
    roster,
    loading,
    overbooking,
    overbookError,
    fetchRoster,
    overbook,
  }
}
