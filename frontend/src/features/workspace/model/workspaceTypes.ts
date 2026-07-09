export type WorkspaceRole = 'owner' | 'coach' | 'athlete'

export interface GymResponse {
  id: string
  name: string
  ownerId: string
}

export interface MemberResponse {
  userId: string
  email: string
  displayName: string
  role: WorkspaceRole
  status: 'pending' | 'active'
}

export interface InvitationResponse {
  id: string
  gymId: string
  email: string
  role: WorkspaceRole
  token?: string
}

export interface InvitationPreviewResponse {
  gymId: string
  gymName: string
  email: string
  role: WorkspaceRole
  status: 'pending' | 'expired' | 'accepted' | 'revoked'
}

export interface CreateGymPayload {
  name: string
}

export interface InviteMemberPayload {
  email: string
}

export interface UpdateMemberRolePayload {
  role: 'coach' | 'athlete'
}

export const ROLE_LABELS: Record<WorkspaceRole, string> = {
  owner: 'Owner',
  coach: 'Coach',
  athlete: 'Athlete',
}

export function canCreateWOD(role: WorkspaceRole | null): boolean {
  return role === 'owner' || role === 'coach'
}

export function canEditWOD(
  role: WorkspaceRole | null,
  wod?: { createdBy: string; status: string },
  userId?: string | null,
): boolean {
  if (role === 'owner') {
    return true
  }
  if (role === 'coach' && wod && userId) {
    return wod.createdBy === userId && (wod.status === 'DRAFT' || wod.status === 'PUBLISHED')
  }
  return false
}

export function canViewWOD(
  role: WorkspaceRole | null,
  wod?: { status: string },
): boolean {
  if (!role) {
    return false
  }
  if (role === 'athlete' && wod?.status !== 'PUBLISHED') {
    return false
  }
  return true
}

export function canPublishWOD(role: WorkspaceRole | null): boolean {
  return role === 'owner' || role === 'coach'
}

export function canManageMembers(role: WorkspaceRole | null): boolean {
  return role === 'owner'
}

export function canManageMemberTarget(role: WorkspaceRole): boolean {
  return role === 'coach' || role === 'athlete'
}
