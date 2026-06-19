import { httpClient } from '@/shared/api/httpClient'
import type {
  CalendarDaySummary,
  CreateWODPayload,
  CreateWODResponse,
  WODDetail,
  WODStatus,
  WODSummary,
} from '@/features/wod/model/wodTypes'
import type { Result } from '@/shared/utils/result'

export function createWOD(payload: CreateWODPayload): Promise<Result<CreateWODResponse>> {
  return httpClient.post<CreateWODResponse>('/api/v1/wods', payload, { auth: true, workspace: true })
}

export function listWODs(options?: { status?: WODStatus }): Promise<Result<WODSummary[]>> {
  const query = options?.status ? `?status=${encodeURIComponent(options.status)}` : ''
  return httpClient.get<WODSummary[]>(`/api/v1/wods${query}`, { auth: true, workspace: true })
}

export function getWODCalendar(from: string, to: string): Promise<Result<CalendarDaySummary[]>> {
  const params = new URLSearchParams({ from, to })
  return httpClient.get<CalendarDaySummary[]>(`/api/v1/wods/calendar?${params.toString()}`, {
    auth: true,
    workspace: true,
  })
}

export function getWOD(id: string): Promise<Result<WODDetail>> {
  return httpClient.get<WODDetail>(`/api/v1/wods/${id}`, { auth: true, workspace: true })
}

export function updateWOD(id: string, payload: CreateWODPayload): Promise<Result<WODDetail>> {
  return httpClient.put<WODDetail>(`/api/v1/wods/${id}`, payload, { auth: true, workspace: true })
}

export function publishWOD(id: string): Promise<Result<WODDetail>> {
  return httpClient.post<WODDetail>(`/api/v1/wods/${id}/publish`, {}, { auth: true, workspace: true })
}
