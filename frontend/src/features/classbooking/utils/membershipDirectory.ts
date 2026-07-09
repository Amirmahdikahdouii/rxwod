import { listMembers } from '@/features/workspace/api/workspaceApi'
import { err, ok, type Result } from '@/shared/utils/result'

export interface MembershipEntry {
  displayName: string
  userId: string
}

export type MembershipDirectory = Map<string, MembershipEntry>

const PAGE_LIMIT = 100

export async function fetchMembershipDirectory(gymId: string): Promise<Result<MembershipDirectory>> {
  const directory: MembershipDirectory = new Map()
  let page = 1

  while (true) {
    const response = await listMembers(gymId, { page, limit: PAGE_LIMIT })
    if (!response.ok) {
      return err(response.error)
    }

    for (const member of response.value.data) {
      directory.set(member.membershipId, {
        displayName: member.displayName,
        userId: member.userId,
      })
    }

    if (page >= response.value.meta.totalPages) {
      break
    }
    page++
  }

  return ok(directory)
}
