import { describe, expect, it } from 'vitest'
import {
  configToPayload,
  defaultConfigForType,
  defaultMovement,
  defaultStage,
  detailToFormState,
  stageToPayload,
} from '@/features/wod/model/wodSchemas'
import type { WODDetail } from '@/features/wod/model/wodTypes'

describe('wodSchemas', () => {
  it('builds AMRAP payload', () => {
    const config = defaultConfigForType('AMRAP')
    expect(configToPayload(config)).toEqual({ timeCapSeconds: 900 })
  })

  it('builds OPEN payload as empty config', () => {
    const config = defaultConfigForType('OPEN')
    expect(configToPayload(config)).toEqual({})
  })

  it('builds TABATA payload', () => {
    const config = defaultConfigForType('TABATA')
    expect(configToPayload(config)).toEqual({
      workSeconds: 20,
      restSeconds: 10,
      rounds: 8,
      cycles: 1,
    })
  })

  it('creates default warmup stage as OPEN', () => {
    expect(defaultStage('WARMUP')).toEqual({
      kind: 'WARMUP',
      type: 'OPEN',
      instructions: '',
      config: { type: 'OPEN' },
      movements: [{ position: 1, label: '', name: '', prescription: '', reps: undefined }],
    })
  })

  it('creates default metcon stage as AMRAP', () => {
    expect(defaultStage('METCON').type).toBe('AMRAP')
  })

  it('maps stage to API payload with prescriptions', () => {
    const stage = {
      ...defaultStage('STRENGTH', 'OPEN'),
      instructions: 'Complete in 20 minutes.',
      movements: [{
        ...defaultMovement(),
        label: 'B',
        name: 'Close Grip Bench Press',
        prescription: '3RM',
        sets: 5,
        reps: 3,
      }],
    }

    expect(stageToPayload(stage)).toEqual({
      kind: 'STRENGTH',
      type: 'OPEN',
      instructions: 'Complete in 20 minutes.',
      config: {},
      movements: [{
        position: 1,
        label: 'B',
        name: 'Close Grip Bench Press',
        prescription: '3RM',
        sets: 5,
        reps: 3,
        notes: '',
      }],
    })
  })

  it('omits zero load value from API payload', () => {
    const stage = {
      ...defaultStage('WARMUP', 'OPEN'),
      movements: [{
        ...defaultMovement(),
        label: 'A',
        name: 'Wall facing handstand plate stepup',
        reps: 20,
        sets: 1,
        loadValue: 0,
      }],
    }

    expect(stageToPayload(stage).movements[0]).toEqual({
      position: 1,
      label: 'A',
      name: 'Wall facing handstand plate stepup',
      prescription: '',
      reps: 20,
      sets: 1,
      notes: '',
    })
  })

  it('hydrates detail responses into form state', () => {
    const detail: WODDetail = {
      id: 'wod-1',
      name: 'Strength Day',
      description: 'Coach plan',
      status: 'DRAFT',
      createdBy: 'coach-1',
      createdAt: '2026-01-01T00:00:00Z',
      updatedAt: '2026-01-01T00:00:00Z',
      stages: [{
        id: 'stage-1',
        kind: 'STRENGTH',
        position: 1,
        instructions: 'Complete in 20 minutes.',
        type: 'OPEN',
        scoringKind: 'NONE',
        config: {},
        movements: [{
          id: 'movement-1',
          position: 1,
          label: 'A',
          name: 'Back Squat',
          prescription: 'Heavy triple',
          sets: 5,
          reps: 3,
          notes: 'Use a rack',
        }],
      }],
    }

    expect(detailToFormState(detail)).toEqual({
      name: 'Strength Day',
      description: 'Coach plan',
      stages: [{
        kind: 'STRENGTH',
        type: 'OPEN',
        instructions: 'Complete in 20 minutes.',
        config: { type: 'OPEN' },
        movements: [{
          position: 1,
          label: 'A',
          name: 'Back Squat',
          prescription: 'Heavy triple',
          sets: 5,
          reps: 3,
          loadValue: undefined,
          loadUnit: undefined,
          notes: 'Use a rack',
        }],
      }],
    })
  })
})
