import { beforeEach, describe, expect, it } from 'vitest'
import {
  clearStoredDraft,
  draftStorageKey,
  draftsMatch,
  hasMeaningfulDraft,
  loadStoredDraft,
  saveStoredDraft,
  type StoredWODDraft,
} from '@/features/wod/composables/useWODDraftStorage'
import { defaultStage } from '@/features/wod/model/wodSchemas'

const key = draftStorageKey('gym-1')

const sampleDraft: StoredWODDraft = {
  name: 'Monday Session',
  description: 'Notes',
  scheduledDate: '2026-06-20',
  stages: [defaultStage()],
  mode: 'create',
  wodId: null,
  savedAt: '2026-06-19T10:00:00.000Z',
}

describe('useWODDraftStorage', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('saves and loads drafts from localStorage', () => {
    saveStoredDraft(key, sampleDraft)
    expect(loadStoredDraft(key)).toEqual(sampleDraft)
  })

  it('clears stored drafts', () => {
    saveStoredDraft(key, sampleDraft)
    clearStoredDraft(key)
    expect(loadStoredDraft(key)).toBeNull()
  })

  it('detects meaningful draft content', () => {
    expect(hasMeaningfulDraft(sampleDraft)).toBe(true)
    expect(hasMeaningfulDraft({ ...sampleDraft, name: '', description: '', scheduledDate: '', stages: [defaultStage()] })).toBe(false)
  })

  it('compares draft snapshots', () => {
    expect(draftsMatch(sampleDraft, sampleDraft)).toBe(true)
    expect(draftsMatch(sampleDraft, { ...sampleDraft, name: 'Changed' })).toBe(false)
  })
})
