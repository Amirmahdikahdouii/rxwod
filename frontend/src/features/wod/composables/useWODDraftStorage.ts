import type { StageFormState } from '@/features/wod/model/wodTypes'

export interface StoredWODDraft {
  name: string
  description: string
  scheduledDate: string
  stages: StageFormState[]
  mode: 'create' | 'edit'
  wodId: string | null
  savedAt: string
}

const STORAGE_PREFIX = 'rxwod.wodDraft'

export function draftStorageKey(workspaceId: string, wodId?: string | null): string {
  if (wodId) {
    return `${STORAGE_PREFIX}.${workspaceId}.${wodId}`
  }
  return `${STORAGE_PREFIX}.${workspaceId}`
}

export function loadStoredDraft(key: string): StoredWODDraft | null {
  const raw = localStorage.getItem(key)
  if (!raw) {
    return null
  }
  try {
    return JSON.parse(raw) as StoredWODDraft
  } catch {
    return null
  }
}

export function saveStoredDraft(key: string, draft: StoredWODDraft): void {
  localStorage.setItem(key, JSON.stringify(draft))
}

export function clearStoredDraft(key: string): void {
  localStorage.removeItem(key)
}

export function hasMeaningfulDraft(draft: StoredWODDraft): boolean {
  return Boolean(
    draft.name.trim() ||
      draft.description.trim() ||
      draft.scheduledDate ||
      draft.stages.some(
        (stage) =>
          stage.instructions.trim() ||
          stage.movements.some(
            (movement) => movement.name.trim() || movement.prescription?.trim() || movement.label?.trim(),
          ),
      ),
  )
}

export function formatDraftAge(savedAt: string): string {
  const saved = new Date(savedAt)
  const diffMs = Date.now() - saved.getTime()
  const minutes = Math.floor(diffMs / 60000)
  if (minutes < 1) {
    return 'just now'
  }
  if (minutes < 60) {
    return `${minutes} minute${minutes === 1 ? '' : 's'} ago`
  }
  const hours = Math.floor(minutes / 60)
  if (hours < 24) {
    return `${hours} hour${hours === 1 ? '' : 's'} ago`
  }
  return saved.toLocaleString()
}

export function draftsMatch(
  draft: StoredWODDraft,
  current: Pick<StoredWODDraft, 'name' | 'description' | 'scheduledDate' | 'stages'>,
): boolean {
  const snapshot = (value: Pick<StoredWODDraft, 'name' | 'description' | 'scheduledDate' | 'stages'>) =>
    JSON.stringify({
      name: value.name,
      description: value.description,
      scheduledDate: value.scheduledDate,
      stages: value.stages,
    })

  return snapshot(draft) === snapshot(current)
}
