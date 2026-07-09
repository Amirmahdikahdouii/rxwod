import { httpClient } from '@/shared/api/httpClient'
import type {
  LeaderboardResponse,
  SubmitResultPayload,
  WODResultResponse,
} from '@/features/wodresult/model/wodResultTypes'
import type { Result } from '@/shared/utils/result'

export function submitResult(
  wodId: string,
  payload: SubmitResultPayload,
): Promise<Result<WODResultResponse>> {
  return httpClient.post<WODResultResponse>(`/api/v1/wods/${wodId}/results`, payload, {
    auth: true,
    workspace: true,
  })
}

export function getLeaderboard(wodId: string): Promise<Result<LeaderboardResponse>> {
  return httpClient.get<LeaderboardResponse>(`/api/v1/wods/${wodId}/leaderboard`, {
    auth: true,
    workspace: true,
  })
}
