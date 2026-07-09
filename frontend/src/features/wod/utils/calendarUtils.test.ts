import { describe, expect, it } from 'vitest'
import { canDeleteWOD, canEditWOD, canPublishWOD, canViewWOD } from '@/features/workspace/model/workspaceTypes'

describe('canEditWOD', () => {
  it('allows owners to edit any program', () => {
    expect(canEditWOD('owner', { createdBy: 'coach-1', status: 'PUBLISHED' }, 'owner-1')).toBe(true)
  })

  it('blocks editing archived programs for all roles', () => {
    const archived = { createdBy: 'owner-1', status: 'ARCHIVED' }
    expect(canEditWOD('owner', archived, 'owner-1')).toBe(false)
    expect(canEditWOD('coach', archived, 'coach-1')).toBe(false)
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

  it('blocks athletes from viewing drafts and archived programs', () => {
    const draft = { status: 'DRAFT' }
    expect(canViewWOD('owner', draft)).toBe(true)
    expect(canViewWOD('coach', draft)).toBe(true)
    expect(canViewWOD('athlete', draft)).toBe(false)

    const archived = { status: 'ARCHIVED' }
    expect(canViewWOD('owner', archived)).toBe(true)
    expect(canViewWOD('coach', archived)).toBe(true)
    expect(canViewWOD('athlete', archived)).toBe(false)
  })
})

describe('canPublishWOD', () => {
  it('allows owners and coaches to publish', () => {
    expect(canPublishWOD('owner')).toBe(true)
    expect(canPublishWOD('coach')).toBe(true)
    expect(canPublishWOD('athlete')).toBe(false)
  })
})

describe('canDeleteWOD', () => {
  it('allows only owners to archive or delete programs', () => {
    expect(canDeleteWOD('owner')).toBe(true)
    expect(canDeleteWOD('coach')).toBe(false)
    expect(canDeleteWOD('athlete')).toBe(false)
    expect(canDeleteWOD(null)).toBe(false)
  })
})
