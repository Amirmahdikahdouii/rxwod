import { describe, expect, it } from 'vitest'
import {
  configToPayload,
  defaultConfigForType,
  defaultMovement,
  defaultStage,
  stageToPayload,
} from '@/features/wod/model/wodSchemas'

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
        notes: '',
      }],
    })
  })
})
