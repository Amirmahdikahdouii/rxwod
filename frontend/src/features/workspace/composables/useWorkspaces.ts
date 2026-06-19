import { computed, ref } from 'vue'
import { useSession } from '@/features/auth/composables/useSession'
import {
  createGym as createGymRequest,
  inviteAthlete as inviteAthleteRequest,
  inviteCoach as inviteCoachRequest,
  listMembers as listMembersRequest,
  removeMember as removeMemberRequest,
  updateMemberRole as updateMemberRoleRequest,
} from '@/features/workspace/api/workspaceApi'
import type { InvitationResponse, MemberResponse } from '@/features/workspace/model/workspaceTypes'
import { canManageMembers } from '@/features/workspace/model/workspaceTypes'
import { err, ok, type Result } from '@/shared/utils/result'

const members = ref<MemberResponse[]>([])
const loadingMembers = ref(false)
const updatingMemberId = ref<string | null>(null)
const removingMemberId = ref<string | null>(null)
const workspaceError = ref<string | null>(null)
const workspaceSuccess = ref<string | null>(null)

const SUCCESS_DISMISS_MS = 4000
let successDismissTimer: ReturnType<typeof setTimeout> | null = null

function clearWorkspaceSuccess() {
  workspaceSuccess.value = null
  if (successDismissTimer) {
    clearTimeout(successDismissTimer)
    successDismissTimer = null
  }
}

function setWorkspaceSuccess(message: string) {
  workspaceSuccess.value = message
  if (successDismissTimer) {
    clearTimeout(successDismissTimer)
  }
  successDismissTimer = setTimeout(() => {
    workspaceSuccess.value = null
    successDismissTimer = null
  }, SUCCESS_DISMISS_MS)
}

export function useWorkspaces() {
  const session = useSession()
  const canManageActiveWorkspace = computed(() => canManageMembers(session.activeWorkspaceRole.value))

  async function refreshMembers(): Promise<Result<MemberResponse[]>> {
    const workspaceID = session.activeWorkspaceId.value
    if (!workspaceID) {
      return err('Choose a workspace first')
    }

    loadingMembers.value = true
    workspaceError.value = null
    const response = await listMembersRequest(workspaceID)
    loadingMembers.value = false

    if (!response.ok) {
      workspaceError.value = response.error
      return err(response.error)
    }

    members.value = response.value
    return ok(response.value)
  }

  async function createGym(name: string): Promise<Result<void>> {
    workspaceError.value = null
    clearWorkspaceSuccess()
    const response = await createGymRequest({ name: name.trim() })
    if (!response.ok) {
      workspaceError.value = response.error
      return err(response.error)
    }

    await session.loadMe()
    session.setActiveWorkspace(response.value.id)
    setWorkspaceSuccess(`${response.value.name} is ready.`)
    return ok(undefined)
  }

  async function invite(
    role: 'coach' | 'athlete',
    email: string,
  ): Promise<Result<InvitationResponse>> {
    const workspaceID = session.activeWorkspaceId.value
    if (!workspaceID) {
      return err('Choose a workspace first')
    }

    workspaceError.value = null
    clearWorkspaceSuccess()
    const response =
      role === 'coach'
        ? await inviteCoachRequest(workspaceID, { email: email.trim() })
        : await inviteAthleteRequest(workspaceID, { email: email.trim() })

    if (!response.ok) {
      workspaceError.value = response.error
      return err(response.error)
    }

    setWorkspaceSuccess(`${response.value.email} was invited as ${response.value.role}.`)
    await refreshMembers()
    return ok(response.value)
  }

  async function updateMemberRole(
    userID: string,
    role: 'coach' | 'athlete',
  ): Promise<Result<MemberResponse>> {
    const workspaceID = session.activeWorkspaceId.value
    if (!workspaceID) {
      return err('Choose a workspace first')
    }

    workspaceError.value = null
    clearWorkspaceSuccess()
    updatingMemberId.value = userID
    const response = await updateMemberRoleRequest(workspaceID, userID, { role })
    updatingMemberId.value = null

    if (!response.ok) {
      workspaceError.value = response.error
      return err(response.error)
    }

    setWorkspaceSuccess(`Member role updated to ${response.value.role}.`)
    await refreshMembers()
    return ok(response.value)
  }

  async function removeMember(userID: string): Promise<Result<void>> {
    const workspaceID = session.activeWorkspaceId.value
    if (!workspaceID) {
      return err('Choose a workspace first')
    }

    workspaceError.value = null
    clearWorkspaceSuccess()
    removingMemberId.value = userID
    const response = await removeMemberRequest(workspaceID, userID)
    removingMemberId.value = null

    if (!response.ok) {
      workspaceError.value = response.error
      return err(response.error)
    }

    setWorkspaceSuccess('Member removed from workspace.')
    await refreshMembers()
    return ok(undefined)
  }

  return {
    members,
    loadingMembers,
    updatingMemberId,
    removingMemberId,
    workspaceError,
    workspaceSuccess,
    canManageActiveWorkspace,
    clearWorkspaceSuccess,
    refreshMembers,
    createGym,
    invite,
    updateMemberRole,
    removeMember,
  }
}
