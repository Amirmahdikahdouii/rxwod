import { httpClient } from '@/shared/api/httpClient'
import type {
  ClassBookingResponse,
  ClassSessionResponse,
  CreateClassSessionPayload,
  CreateClassSessionResponse,
  DefaultSessionResponse,
  SetDefaultSessionPayload,
} from '@/features/classbooking/model/classBookingTypes'
import type { Result } from '@/shared/utils/result'

export function listSessions(from: string, to: string): Promise<Result<ClassSessionResponse[]>> {
  const params = new URLSearchParams({ from, to })
  return httpClient.get<ClassSessionResponse[]>(`/api/v1/sessions?${params.toString()}`, {
    auth: true,
    workspace: true,
  })
}

export function createSession(
  payload: CreateClassSessionPayload,
): Promise<Result<CreateClassSessionResponse>> {
  return httpClient.post<CreateClassSessionResponse>('/api/v1/sessions', payload, {
    auth: true,
    workspace: true,
  })
}

export function getRoster(sessionId: string): Promise<Result<ClassBookingResponse[]>> {
  return httpClient.get<ClassBookingResponse[]>(`/api/v1/sessions/${sessionId}/bookings`, {
    auth: true,
    workspace: true,
  })
}

export function bookSession(sessionId: string): Promise<Result<ClassBookingResponse>> {
  return httpClient.post<ClassBookingResponse>(
    `/api/v1/sessions/${sessionId}/book`,
    {},
    { auth: true, workspace: true },
  )
}

export function overbookSession(
  sessionId: string,
  athleteUserId: string,
): Promise<Result<ClassBookingResponse>> {
  return httpClient.post<ClassBookingResponse>(
    `/api/v1/sessions/${sessionId}/overbook`,
    { athleteUserId },
    { auth: true, workspace: true },
  )
}

export function cancelBooking(
  sessionId: string,
  athleteUserId?: string,
): Promise<Result<ClassBookingResponse>> {
  const body = athleteUserId ? { athleteUserId } : {}
  return httpClient.post<ClassBookingResponse>(
    `/api/v1/sessions/${sessionId}/cancel`,
    body,
    { auth: true, workspace: true },
  )
}

export function getMyDefaultSessions(): Promise<Result<DefaultSessionResponse[]>> {
  return httpClient.get<DefaultSessionResponse[]>('/api/v1/athletes/preferences/default-session', {
    auth: true,
    workspace: true,
  })
}

export function setDefaultSession(
  payload: SetDefaultSessionPayload,
): Promise<Result<DefaultSessionResponse>> {
  return httpClient.post<DefaultSessionResponse>(
    '/api/v1/athletes/preferences/default-session',
    payload,
    { auth: true, workspace: true },
  )
}

export function removeDefaultSession(id: string): Promise<Result<void>> {
  return httpClient.delete<void>(`/api/v1/athletes/preferences/default-session/${id}`, {
    auth: true,
    workspace: true,
  })
}
