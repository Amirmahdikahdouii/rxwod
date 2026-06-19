import { describe, expect, it } from 'vitest'
import { heatmapLevel } from '@/features/wod/utils/calendarUtils'
import { canEditWOD, canPublishWOD } from '@/features/workspace/model/workspaceTypes'

describe('heatmapLevel', () => {
  it('maps published counts to intensity levels', () => {
    expect(heatmapLevel(0)).toBe(0)
    expect(heatmapLevel(1)).toBe(1)
    expect(heatmapLevel(2)).toBe(2)
    expect(heatmapLevel(3)).toBe(3)
    expect(heatmapLevel(4)).toBe(4)
    expect(heatmapLevel(10)).toBe(4)
  })
})

describe('canEditWOD', () => {
  it('allows owners to edit any program', () => {
    expect(canEditWOD('owner', { createdBy: 'coach-1', status: 'PUBLISHED' }, 'owner-1')).toBe(true)
  })

  it('allows coaches to edit their own drafts only', () => {
    expect(canEditWOD('coach', { createdBy: 'coach-1', status: 'DRAFT' }, 'coach-1')).toBe(true)
    expect(canEditWOD('coach', { createdBy: 'coach-2', status: 'DRAFT' }, 'coach-1')).toBe(false)
    expect(canEditWOD('coach', { createdBy: 'coach-1', status: 'PUBLISHED' }, 'coach-1')).toBe(false)
  })
})

describe('canPublishWOD', () => {
  it('allows owners and coaches to publish', () => {
    expect(canPublishWOD('owner')).toBe(true)
    expect(canPublishWOD('coach')).toBe(true)
    expect(canPublishWOD('athlete')).toBe(false)
  })
})
