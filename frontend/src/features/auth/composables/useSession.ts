import { computed, ref } from 'vue'
import { getMe, login as loginRequest, register as registerRequest } from '@/features/auth/api/authApi'
import type {
  LoginPayload,
  RegisterPayload,
  UserProfile,
  UserWorkspace,
} from '@/features/auth/model/authTypes'
import { httpClient } from '@/shared/api/httpClient'
import { err, ok, type Result } from '@/shared/utils/result'

const accessTokenKey = 'rxwod.accessToken'
const refreshTokenKey = 'rxwod.refreshToken'
const activeWorkspaceKey = 'rxwod.activeWorkspaceId'

const accessToken = ref<string | null>(localStorage.getItem(accessTokenKey))
const refreshToken = ref<string | null>(localStorage.getItem(refreshTokenKey))
const currentUser = ref<UserProfile | null>(null)
const workspaces = ref<UserWorkspace[]>([])
const activeWorkspaceId = ref<string | null>(localStorage.getItem(activeWorkspaceKey))
const loadingSession = ref(false)
const sessionReady = ref(false)

const isAuthenticated = computed(() => Boolean(accessToken.value))
const activeWorkspace = computed(
  () => workspaces.value.find((workspace) => workspace.id === activeWorkspaceId.value) ?? null,
)
const activeWorkspaceRole = computed(() => activeWorkspace.value?.role ?? null)

function persistTokens(nextAccessToken: string, nextRefreshToken?: string) {
  accessToken.value = nextAccessToken
  localStorage.setItem(accessTokenKey, nextAccessToken)

  if (nextRefreshToken) {
    refreshToken.value = nextRefreshToken
    localStorage.setItem(refreshTokenKey, nextRefreshToken)
  }
}

function selectDefaultWorkspace() {
  if (workspaces.value.length === 0) {
    setActiveWorkspace(null)
    return
  }

  const savedWorkspace = workspaces.value.find(
    (workspace) => workspace.id === activeWorkspaceId.value,
  )
  setActiveWorkspace(savedWorkspace?.id ?? workspaces.value[0].id)
}

function setActiveWorkspace(workspaceId: string | null) {
  activeWorkspaceId.value = workspaceId
  if (workspaceId) {
    localStorage.setItem(activeWorkspaceKey, workspaceId)
  } else {
    localStorage.removeItem(activeWorkspaceKey)
  }
}

function clearSession() {
  accessToken.value = null
  refreshToken.value = null
  currentUser.value = null
  workspaces.value = []
  activeWorkspaceId.value = null
  localStorage.removeItem(accessTokenKey)
  localStorage.removeItem(refreshTokenKey)
  localStorage.removeItem(activeWorkspaceKey)
}

async function loadMe(): Promise<Result<UserProfile>> {
  if (!accessToken.value) {
    sessionReady.value = true
    return err('Login is required')
  }

  loadingSession.value = true
  const response = await getMe()
  loadingSession.value = false
  sessionReady.value = true

  if (!response.ok) {
    clearSession()
    return err(response.error)
  }

  currentUser.value = response.value.user
  workspaces.value = response.value.gyms
  selectDefaultWorkspace()
  return ok(response.value.user)
}

async function login(payload: LoginPayload): Promise<Result<UserProfile>> {
  const response = await loginRequest(payload)
  if (!response.ok) {
    return err(response.error)
  }
  persistTokens(response.value.accessToken, response.value.refreshToken)
  return loadMe()
}

async function register(payload: RegisterPayload): Promise<Result<UserProfile>> {
  const response = await registerRequest(payload)
  if (!response.ok) {
    return err(response.error)
  }
  persistTokens(response.value.accessToken, response.value.refreshToken)
  return loadMe()
}

function logout() {
  clearSession()
}

httpClient.configure({
  tokenProvider: () => accessToken.value,
  workspaceProvider: () => activeWorkspaceId.value,
  unauthorizedHandler: () => {
    clearSession()
  },
})

export function useSession() {
  return {
    accessToken,
    refreshToken,
    currentUser,
    workspaces,
    activeWorkspaceId,
    activeWorkspace,
    activeWorkspaceRole,
    loadingSession,
    sessionReady,
    isAuthenticated,
    loadMe,
    login,
    register,
    logout,
    setActiveWorkspace,
  }
}
