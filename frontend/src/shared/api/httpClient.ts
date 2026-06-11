import { err, ok, type Result } from '@/shared/utils/result'

const baseURL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

async function request<T>(path: string, init?: RequestInit): Promise<Result<T>> {
  try {
    const response = await fetch(`${baseURL}${path}`, {
      headers: {
        'Content-Type': 'application/json',
        ...(init?.headers ?? {}),
      },
      ...init,
    })

    if (!response.ok) {
      const body = await response.json().catch(() => ({ error: 'Request failed' }))
      return err(typeof body.error === 'string' ? body.error : 'Request failed')
    }

    const data = (await response.json()) as T
    return ok(data)
  } catch (error) {
    return err(error instanceof Error ? error.message : 'Network error')
  }
}

export const httpClient = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: 'POST', body: JSON.stringify(body) }),
}
