import { httpClient } from '@/shared/api/httpClient'
import type {
  ForgotPasswordPayload,
  LoginPayload,
  MeResponse,
  RegisterPayload,
  ResetPasswordPayload,
  TokenResponse,
  VerifyEmailPayload,
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

export function forgotPassword(payload: ForgotPasswordPayload): Promise<Result<void>> {
  return httpClient.post<void>('/api/v1/auth/forgot-password', payload)
}

export function resetPassword(payload: ResetPasswordPayload): Promise<Result<void>> {
  return httpClient.post<void>('/api/v1/auth/reset-password', payload)
}

export function verifyEmail(payload: VerifyEmailPayload): Promise<Result<void>> {
  return httpClient.post<void>('/api/v1/auth/verify-email', payload)
}

export function resendVerificationEmail(): Promise<Result<void>> {
  return httpClient.post<void>('/api/v1/auth/resend-verification', undefined, { auth: true })
}
