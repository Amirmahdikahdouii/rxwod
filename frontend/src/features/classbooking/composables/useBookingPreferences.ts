import { ref } from 'vue'
import {
  getMyDefaultSessions,
  removeDefaultSession as removeDefaultSessionRequest,
  setDefaultSession,
} from '@/features/classbooking/api/classBookingApi'
import type {
  DefaultSessionResponse,
  SetDefaultSessionPayload,
} from '@/features/classbooking/model/classBookingTypes'
import { err, ok, type Result } from '@/shared/utils/result'

export function useBookingPreferences() {
  const defaultSessions = ref<DefaultSessionResponse[]>([])
  const loading = ref(false)
  const submitting = ref(false)
  const error = ref<string | null>(null)

  async function fetchMine(): Promise<Result<DefaultSessionResponse[]>> {
    loading.value = true
    error.value = null

    const response = await getMyDefaultSessions()
    loading.value = false

    if (!response.ok) {
      error.value = response.error
      return err(response.error)
    }

    defaultSessions.value = response.value
    return ok(response.value)
  }

  async function addDefaultSession(
    payload: SetDefaultSessionPayload,
  ): Promise<Result<DefaultSessionResponse>> {
    submitting.value = true
    error.value = null

    const response = await setDefaultSession(payload)
    submitting.value = false

    if (!response.ok) {
      error.value = response.error
      return err(response.error)
    }

    await fetchMine()
    return ok(response.value)
  }

  async function removeDefaultSession(id: string): Promise<Result<void>> {
    error.value = null

    const response = await removeDefaultSessionRequest(id)
    if (!response.ok) {
      error.value = response.error
      return err(response.error)
    }

    defaultSessions.value = defaultSessions.value.filter((session) => session.id !== id)
    return ok(undefined)
  }

  return {
    defaultSessions,
    loading,
    submitting,
    error,
    fetchMine,
    addDefaultSession,
    removeDefaultSession,
  }
}
