import { ref } from 'vue'
import { getLeaderboard, submitResult } from '@/features/wodresult/api/wodResultApi'
import type {
  LeaderboardEntry,
  SubmitResultPayload,
  WODResultResponse,
} from '@/features/wodresult/model/wodResultTypes'
import { err, ok, type Result } from '@/shared/utils/result'

export function useWODResults() {
  const submitting = ref(false)
  const submitError = ref<string | null>(null)
  const loadingLeaderboard = ref(false)
  const leaderboardError = ref<string | null>(null)
  const leaderboard = ref<LeaderboardEntry[]>([])

  async function submit(
    wodId: string,
    payload: SubmitResultPayload,
  ): Promise<Result<WODResultResponse>> {
    submitting.value = true
    submitError.value = null

    const response = await submitResult(wodId, payload)
    submitting.value = false

    if (!response.ok) {
      submitError.value = response.error
      return err(response.error)
    }

    return ok(response.value)
  }

  async function fetchLeaderboard(wodId: string): Promise<Result<LeaderboardEntry[]>> {
    loadingLeaderboard.value = true
    leaderboardError.value = null

    const response = await getLeaderboard(wodId)
    loadingLeaderboard.value = false

    if (!response.ok) {
      leaderboardError.value = response.error
      return err(response.error)
    }

    leaderboard.value = response.value.entries
    return ok(response.value.entries)
  }

  return {
    submitting,
    submitError,
    loadingLeaderboard,
    leaderboardError,
    leaderboard,
    submit,
    fetchLeaderboard,
  }
}
