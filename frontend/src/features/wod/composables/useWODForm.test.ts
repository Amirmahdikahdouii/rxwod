import { defineComponent, nextTick } from 'vue'
import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { useWODForm } from '@/features/wod/composables/useWODForm'
import { defaultStage } from '@/features/wod/model/wodSchemas'
import type { WODDetail } from '@/features/wod/model/wodTypes'

function mountForm() {
  let form!: ReturnType<typeof useWODForm>
  const Comp = defineComponent({
    setup() {
      form = useWODForm()
      return () => null
    },
  })
  mount(Comp)
  return form
}

describe('useWODForm', () => {
  it('builds open prescription payload', () => {
    const form = mountForm()
    form.name.value = 'Monday Session'
    form.description.value = 'Full class plan'
    form.stages.value = [
      {
        ...defaultStage('WARMUP', 'OPEN'),
        movements: [{
          position: 1,
          label: 'D',
          name: 'Crow + Knee Lift Off',
          prescription: '1-2X4 lift offs per side, rest as needed.',
          sets: 3,
          reps: 4,
        }],
      },
      {
        ...defaultStage('STRENGTH', 'OPEN'),
        instructions: 'Complete in 20 minutes.',
        movements: [{
          position: 1,
          label: 'B',
          name: 'Close Grip Bench Press',
          prescription: '3RM',
        }],
      },
    ]

    expect(form.payload.value).toEqual({
      name: 'Monday Session',
      description: 'Full class plan',
      stages: [
        {
          kind: 'WARMUP',
          type: 'OPEN',
          instructions: '',
          config: {},
          movements: [{
            position: 1,
            label: 'D',
            name: 'Crow + Knee Lift Off',
            prescription: '1-2X4 lift offs per side, rest as needed.',
            sets: 3,
            reps: 4,
            notes: '',
          }],
        },
        {
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
        },
      ],
    })
  })

  it('defaults warmup stage to OPEN format', () => {
    const form = mountForm()
    expect(form.stages.value[0].type).toBe('OPEN')
  })

  it('adds metcon stages by default with AMRAP config', () => {
    const form = mountForm()
    form.addStage()

    expect(form.stages.value[1].kind).toBe('METCON')
    expect(form.stages.value[1].config).toEqual({
      type: 'AMRAP',
      timeCapSeconds: 900,
    })
  })

  it('adds selected non-metcon stages with OPEN config', () => {
    const form = mountForm()
    form.addStage('STRENGTH')

    expect(form.stages.value[1].kind).toBe('STRENGTH')
    expect(form.stages.value[1].config).toEqual({ type: 'OPEN' })
  })

  it('loads detail responses into edit mode', () => {
    const form = mountForm()
    const detail: WODDetail = {
      id: 'wod-1',
      name: 'Strength Day',
      description: 'Coach plan',
      status: 'DRAFT',
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
        }],
      }],
    }

    form.loadFromDetail(detail)

    expect(form.mode.value).toBe('edit')
    expect(form.wodId.value).toBe('wod-1')
    expect(form.name.value).toBe('Strength Day')
    expect(form.stages.value[0].movements[0].sets).toBe(5)
  })

  it('switches to structured format with AMRAP defaults for metcon', async () => {
    const form = mountForm()
    form.updateStageFormat(0, 'STRUCTURED')
    await nextTick()

    expect(form.stages.value[0].config).toEqual({
      type: 'FORTIME',
      rounds: 5,
    })
  })

  it('resets only the changed stage config when structured type changes', async () => {
    const form = mountForm()
    form.addStage()
    form.updateStageFormat(1, 'STRUCTURED')
    form.updateStageType(1, 'TABATA')
    await nextTick()

    expect(form.stages.value[0].type).toBe('OPEN')
    expect(form.stages.value[1].config).toEqual({
      type: 'TABATA',
      workSeconds: 20,
      restSeconds: 10,
      rounds: 8,
      cycles: 1,
    })
  })
})
