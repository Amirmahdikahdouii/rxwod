import type { WorkspaceRole } from '@/features/workspace/model/workspaceTypes'

export interface TokenResponse {
  accessToken: string
  refreshToken?: string
  expiresIn: number
}

export interface UserProfile {
  id: string
  email: string
  displayName: string
}

export interface UserWorkspace {
  id: string
  name: string
  role: WorkspaceRole
}

export interface MeResponse {
  user: UserProfile
  gyms: UserWorkspace[]
}

export interface LoginPayload {
  email: string
  password: string
}

export interface RegisterPayload extends LoginPayload {
  displayName: string
}
