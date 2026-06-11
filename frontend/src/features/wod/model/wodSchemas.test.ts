import { describe, expect, it } from 'vitest'
import { configToPayload, defaultConfigForType, defaultStage, stageToPayload } from '@/features/wod/model/wodSchemas'

describe('wodSchemas', () => {
  it('builds AMRAP payload', () => {
    const config = defaultConfigForType('AMRAP')
    expect(configToPayload(config)).toEqual({ timeCapSeconds: 900 })
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

  it('creates default stage with warmup fortime defaults', () => {
    expect(defaultStage()).toEqual({
      kind: 'WARMUP',
      type: 'FORTIME',
      config: { type: 'FORTIME', rounds: 5 },
      movements: [{ position: 1, name: '', reps: 10 }],
    })
  })

  it('maps stage to API payload with trimmed movement names', () => {
    const stage = {
      ...defaultStage('METCON', 'AMRAP'),
      movements: [{ position: 1, name: '  Burpee  ', reps: 21 }],
    }

    expect(stageToPayload(stage)).toEqual({
      kind: 'METCON',
      type: 'AMRAP',
      config: { timeCapSeconds: 900 },
      movements: [{ position: 1, name: 'Burpee', reps: 21 }],
    })
  })
})
