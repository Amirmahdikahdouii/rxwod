import { defineComponent, nextTick } from 'vue'
import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { useWODForm } from '@/features/wod/composables/useWODForm'
import { defaultStage } from '@/features/wod/model/wodSchemas'

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
