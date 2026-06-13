import { httpClient } from '@/shared/api/httpClient'
import type {
  LoginPayload,
  MeResponse,
  RegisterPayload,
  TokenResponse,
} from '@/features/auth/model/authTypes'
import type { Result } from '@/shared/utils/result'

export function login(payload: LoginPayload): Promise<Result<TokenResponse>> {
  return httpClient.post<TokenResponse>('/api/v1/auth/login', payload)
}

export function register(payload: RegisterPayload): Promise<Result<TokenResponse>> {
  return httpClient.post<TokenResponse>('/api/v1/auth/register', payload)
}

export function refreshToken(refreshToken: string): Promise<Result<TokenResponse>> {
  return httpClient.post<TokenResponse>('/api/v1/auth/refresh', { refreshToken })
}

export function getMe(): Promise<Result<MeResponse>> {
  return httpClient.get<MeResponse>('/api/v1/me', { auth: true })
}
