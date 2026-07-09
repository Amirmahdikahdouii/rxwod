import { httpClient } from '@/shared/api/httpClient'
import type { UserWorkspace } from '@/features/auth/model/authTypes'
import type {
  CreateGymPayload,
  GymResponse,
  InvitationPreviewResponse,
  InvitationResponse,
  InviteMemberPayload,
  MemberResponse,
  UpdateMemberRolePayload,
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

export function updateMemberRole(
  gymID: string,
  userID: string,
  payload: UpdateMemberRolePayload,
): Promise<Result<MemberResponse>> {
  return httpClient.patch<MemberResponse>(`/api/v1/gyms/${gymID}/members/${userID}`, payload, {
    auth: true,
    workspace: true,
  })
}

export function removeMember(gymID: string, userID: string): Promise<Result<void>> {
  return httpClient.delete<void>(`/api/v1/gyms/${gymID}/members/${userID}`, {
    auth: true,
    workspace: true,
  })
}

export function getInvitationPreview(token: string): Promise<Result<InvitationPreviewResponse>> {
  return httpClient.get<InvitationPreviewResponse>(`/api/v1/invitations/${encodeURIComponent(token)}`)
}

export function acceptInvitation(gymId: string, token: string): Promise<Result<MemberResponse>> {
  return httpClient.post<MemberResponse>(`/api/v1/gyms/${gymId}/members/accept`, { token }, { auth: true })
}
