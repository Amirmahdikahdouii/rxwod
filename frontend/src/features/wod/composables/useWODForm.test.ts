import { defineComponent, nextTick } from 'vue'
import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { useWODForm } from '@/features/wod/composables/useWODForm'

function mountForm(initialType: 'AMRAP' | 'TABATA' = 'AMRAP') {
  let form!: ReturnType<typeof useWODForm>
  const Comp = defineComponent({
    setup() {
      form = useWODForm(initialType)
      return () => null
    },
  })
  mount(Comp)
  return form
}

describe('useWODForm', () => {
  it('builds AMRAP payload', () => {
    const form = mountForm('AMRAP')
    form.name.value = 'Test AMRAP'
    form.description.value = 'desc'
    form.movements.value = [{ position: 1, name: 'Burpee', reps: 21 }]

    expect(form.payload.value).toEqual({
      name: 'Test AMRAP',
      type: 'AMRAP',
      description: 'desc',
      config: { timeCapSeconds: 900 },
      movements: [{ position: 1, name: 'Burpee', reps: 21 }],
    })
  })

  it('resets config when type changes', async () => {
    const form = mountForm('AMRAP')
    form.type.value = 'TABATA'
    await nextTick()

    expect(form.config.value).toEqual({
      type: 'TABATA',
      workSeconds: 20,
      restSeconds: 10,
      rounds: 8,
      cycles: 1,
    })
  })
})
