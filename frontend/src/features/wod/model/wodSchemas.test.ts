import { describe, expect, it } from 'vitest'
import { configToPayload, defaultConfigForType } from '@/features/wod/model/wodSchemas'

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
})
