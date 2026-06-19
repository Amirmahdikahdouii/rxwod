import { describe, expect, it } from 'vitest'
import { canEditWOD, canPublishWOD, canViewWOD } from '@/features/workspace/model/workspaceTypes'

describe('canEditWOD', () => {
  it('allows owners to edit any program', () => {
    expect(canEditWOD('owner', { createdBy: 'coach-1', status: 'PUBLISHED' }, 'owner-1')).toBe(true)
  })

  it('allows coaches to edit their own drafts and published programs', () => {
    expect(canEditWOD('coach', { createdBy: 'coach-1', status: 'DRAFT' }, 'coach-1')).toBe(true)
    expect(canEditWOD('coach', { createdBy: 'coach-1', status: 'PUBLISHED' }, 'coach-1')).toBe(true)
    expect(canEditWOD('coach', { createdBy: 'coach-2', status: 'DRAFT' }, 'coach-1')).toBe(false)
    expect(canEditWOD('coach', { createdBy: 'coach-2', status: 'PUBLISHED' }, 'coach-1')).toBe(false)
  })
})

describe('canViewWOD', () => {
  it('allows all roles to view published programs', () => {
    const published = { status: 'PUBLISHED' }
    expect(canViewWOD('owner', published)).toBe(true)
    expect(canViewWOD('coach', published)).toBe(true)
    expect(canViewWOD('athlete', published)).toBe(true)
  })

  it('blocks athletes from viewing drafts', () => {
    const draft = { status: 'DRAFT' }
    expect(canViewWOD('owner', draft)).toBe(true)
    expect(canViewWOD('coach', draft)).toBe(true)
    expect(canViewWOD('athlete', draft)).toBe(false)
  })
})

describe('canPublishWOD', () => {
  it('allows owners and coaches to publish', () => {
    expect(canPublishWOD('owner')).toBe(true)
    expect(canPublishWOD('coach')).toBe(true)
    expect(canPublishWOD('athlete')).toBe(false)
  })
})
