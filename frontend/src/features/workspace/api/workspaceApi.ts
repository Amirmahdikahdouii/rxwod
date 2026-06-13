import { httpClient } from '@/shared/api/httpClient'
import type { UserWorkspace } from '@/features/auth/model/authTypes'
import type {
  CreateGymPayload,
  GymResponse,
  InvitationResponse,
  InviteMemberPayload,
  MemberResponse,
} from '@/features/workspace/model/workspaceTypes'
import type { Result } from '@/shared/utils/result'

export function createGym(payload: CreateGymPayload): Promise<Result<GymResponse>> {
  return httpClient.post<GymResponse>('/api/v1/gyms', payload, { auth: true })
}

export function listGyms(): Promise<Result<UserWorkspace[]>> {
  return httpClient.get<UserWorkspace[]>('/api/v1/gyms', { auth: true })
}

export function getGym(gymID: string): Promise<Result<GymResponse>> {
  return httpClient.get<GymResponse>(`/api/v1/gyms/${gymID}`, { auth: true, workspace: true })
}

export function listMembers(gymID: string): Promise<Result<MemberResponse[]>> {
  return httpClient.get<MemberResponse[]>(`/api/v1/gyms/${gymID}/members`, {
    auth: true,
    workspace: true,
  })
}

export function inviteCoach(
  gymID: string,
  payload: InviteMemberPayload,
): Promise<Result<InvitationResponse>> {
  return httpClient.post<InvitationResponse>(`/api/v1/gyms/${gymID}/coaches`, payload, {
    auth: true,
    workspace: true,
  })
}

export function inviteAthlete(
  gymID: string,
  payload: InviteMemberPayload,
): Promise<Result<InvitationResponse>> {
  return httpClient.post<InvitationResponse>(`/api/v1/gyms/${gymID}/athletes`, payload, {
    auth: true,
    workspace: true,
  })
}
