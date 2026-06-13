import { httpClient } from '@/shared/api/httpClient'
import type { CreateWODPayload, CreateWODResponse, WODDetail, WODSummary } from '@/features/wod/model/wodTypes'
import type { Result } from '@/shared/utils/result'

export function createWOD(payload: CreateWODPayload): Promise<Result<CreateWODResponse>> {
  return httpClient.post<CreateWODResponse>('/api/v1/wods', payload)
}

export function listWODs(): Promise<Result<WODSummary[]>> {
  return httpClient.get<WODSummary[]>('/api/v1/wods')
}

export function getWOD(id: string): Promise<Result<WODDetail>> {
  return httpClient.get<WODDetail>(`/api/v1/wods/${id}`)
}

export function updateWOD(id: string, payload: CreateWODPayload): Promise<Result<WODDetail>> {
  return httpClient.put<WODDetail>(`/api/v1/wods/${id}`, payload)
}
