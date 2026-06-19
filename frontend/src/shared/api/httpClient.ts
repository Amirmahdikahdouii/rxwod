import { err, ok, type Result } from '@/shared/utils/result'

const baseURL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

type TokenProvider = () => string | null
type WorkspaceProvider = () => string | null
type UnauthorizedHandler = () => void

let tokenProvider: TokenProvider = () => null
let workspaceProvider: WorkspaceProvider = () => null
let unauthorizedHandler: UnauthorizedHandler | null = null

interface RequestOptions extends RequestInit {
  auth?: boolean
  workspace?: boolean
}

function buildHeaders(init?: RequestOptions): HeadersInit {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  const incomingHeaders = init?.headers
  if (incomingHeaders instanceof Headers) {
    incomingHeaders.forEach((value, key) => {
      headers[key] = value
    })
  } else if (Array.isArray(incomingHeaders)) {
    for (const [key, value] of incomingHeaders) {
      headers[key] = value
    }
  } else if (incomingHeaders) {
    Object.assign(headers, incomingHeaders)
  }

  if (init?.auth) {
    const token = tokenProvider()
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
  }

  if (init?.workspace) {
    const workspaceID = workspaceProvider()
    if (workspaceID) {
      headers['X-Gym-ID'] = workspaceID
    }
  }

  return headers
}

async function request<T>(path: string, init?: RequestOptions): Promise<Result<T>> {
  try {
    const response = await fetch(`${baseURL}${path}`, {
      ...init,
      headers: buildHeaders(init),
    })

    if (!response.ok) {
      const body = await response.json().catch(() => ({ error: 'Request failed' }))
      if (response.status === 401 && unauthorizedHandler) {
        unauthorizedHandler()
      }
      return err(typeof body.error === 'string' ? body.error : 'Request failed')
    }

    if (response.status === 204) {
      return ok(undefined as T)
    }

    const contentType = response.headers.get('content-type')
    if (!contentType?.includes('application/json')) {
      return ok(undefined as T)
    }

    const data = (await response.json()) as T
    return ok(data)
  } catch (error) {
    return err(error instanceof Error ? error.message : 'Network error')
  }
}

export const httpClient = {
  configure: (config: {
    tokenProvider?: TokenProvider
    workspaceProvider?: WorkspaceProvider
    unauthorizedHandler?: UnauthorizedHandler
  }) => {
    if (config.tokenProvider) {
      tokenProvider = config.tokenProvider
    }
    if (config.workspaceProvider) {
      workspaceProvider = config.workspaceProvider
    }
    if (config.unauthorizedHandler) {
      unauthorizedHandler = config.unauthorizedHandler
    }
  },
  get: <T>(path: string, options?: RequestOptions) => request<T>(path, options),
  post: <T>(path: string, body: unknown, options?: RequestOptions) =>
    request<T>(path, { ...options, method: 'POST', body: JSON.stringify(body) }),
  put: <T>(path: string, body: unknown, options?: RequestOptions) =>
    request<T>(path, { ...options, method: 'PUT', body: JSON.stringify(body) }),
  patch: <T>(path: string, body: unknown, options?: RequestOptions) =>
    request<T>(path, { ...options, method: 'PATCH', body: JSON.stringify(body) }),
  delete: <T>(path: string, options?: RequestOptions) =>
    request<T>(path, { ...options, method: 'DELETE' }),
}
