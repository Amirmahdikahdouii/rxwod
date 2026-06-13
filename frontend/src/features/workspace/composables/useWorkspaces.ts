import { computed, ref } from 'vue'
import { useSession } from '@/features/auth/composables/useSession'
import {
  createGym as createGymRequest,
  inviteAthlete as inviteAthleteRequest,
  inviteCoach as inviteCoachRequest,
  listMembers as listMembersRequest,
} from '@/features/workspace/api/workspaceApi'
import type { InvitationResponse, MemberResponse } from '@/features/workspace/model/workspaceTypes'
import { canManageMembers } from '@/features/workspace/model/workspaceTypes'
import { err, ok, type Result } from '@/shared/utils/result'

const members = ref<MemberResponse[]>([])
const loadingMembers = ref(false)
const workspaceError = ref<string | null>(null)
const workspaceSuccess = ref<string | null>(null)

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
    workspaceSuccess.value = null
    const response = await createGymRequest({ name: name.trim() })
    if (!response.ok) {
      workspaceError.value = response.error
      return err(response.error)
    }

    await session.loadMe()
    session.setActiveWorkspace(response.value.id)
    workspaceSuccess.value = `${response.value.name} is ready.`
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
    workspaceSuccess.value = null
    const response =
      role === 'coach'
        ? await inviteCoachRequest(workspaceID, { email: email.trim() })
        : await inviteAthleteRequest(workspaceID, { email: email.trim() })

    if (!response.ok) {
      workspaceError.value = response.error
      return err(response.error)
    }

    workspaceSuccess.value = `${response.value.email} was invited as ${response.value.role}.`
    await refreshMembers()
    return ok(response.value)
  }

  return {
    members,
    loadingMembers,
    workspaceError,
    workspaceSuccess,
    canManageActiveWorkspace,
    refreshMembers,
    createGym,
    invite,
  }
}
