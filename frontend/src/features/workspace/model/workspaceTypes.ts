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
}

export interface CreateGymPayload {
  name: string
}

export interface InviteMemberPayload {
  email: string
}

export const ROLE_LABELS: Record<WorkspaceRole, string> = {
  owner: 'Owner',
  coach: 'Coach',
  athlete: 'Athlete',
}

export function canCreateWOD(role: WorkspaceRole | null): boolean {
  return role === 'owner' || role === 'coach'
}

export function canEditWOD(role: WorkspaceRole | null): boolean {
  return role === 'owner'
}

export function canManageMembers(role: WorkspaceRole | null): boolean {
  return role === 'owner'
}
