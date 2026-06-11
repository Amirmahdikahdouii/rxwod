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
  it('builds multi-stage payload', () => {
    const form = mountForm()
    form.name.value = 'Monday Session'
    form.description.value = 'Full class plan'
    form.stages.value = [
      {
        ...defaultStage('WARMUP', 'FORTIME'),
        movements: [{ position: 1, name: 'Jumping Jacks', reps: 20 }],
      },
      {
        ...defaultStage('METCON', 'AMRAP'),
        movements: [{ position: 1, name: 'Burpee', reps: 21 }],
      },
    ]

    expect(form.payload.value).toEqual({
      name: 'Monday Session',
      description: 'Full class plan',
      stages: [
        {
          kind: 'WARMUP',
          type: 'FORTIME',
          config: { rounds: 5 },
          movements: [{ position: 1, name: 'Jumping Jacks', reps: 20 }],
        },
        {
          kind: 'METCON',
          type: 'AMRAP',
          config: { timeCapSeconds: 900 },
          movements: [{ position: 1, name: 'Burpee', reps: 21 }],
        },
      ],
    })
  })

  it('adds and removes stages', () => {
    const form = mountForm()
    form.addStage()
    expect(form.stages.value).toHaveLength(2)

    form.removeStage(1)
    expect(form.stages.value).toHaveLength(1)
  })

  it('does not remove the last stage', () => {
    const form = mountForm()
    form.removeStage(0)
    expect(form.stages.value).toHaveLength(1)
  })

  it('resets only the changed stage config when type changes', async () => {
    const form = mountForm()
    form.addStage()
    form.updateStageType(1, 'TABATA')
    await nextTick()

    expect(form.stages.value[0].config).toEqual({
      type: 'FORTIME',
      rounds: 5,
    })
    expect(form.stages.value[1].config).toEqual({
      type: 'TABATA',
      workSeconds: 20,
      restSeconds: 10,
      rounds: 8,
      cycles: 1,
    })
  })
})
