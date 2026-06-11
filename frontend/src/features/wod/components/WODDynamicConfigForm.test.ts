import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import WODDynamicConfigForm from '@/features/wod/components/WODDynamicConfigForm.vue'
import { defaultConfigForType } from '@/features/wod/model/wodSchemas'

describe('WODDynamicConfigForm', () => {
  it('renders AMRAP fields', () => {
    const wrapper = mount(WODDynamicConfigForm, {
      props: {
        type: 'AMRAP',
        config: defaultConfigForType('AMRAP'),
      },
    })

    expect(wrapper.text()).toContain('Time Cap (seconds)')
    expect(wrapper.find('input[type="number"]').exists()).toBe(true)
  })

  it('renders TABATA fields', () => {
    const wrapper = mount(WODDynamicConfigForm, {
      props: {
        type: 'TABATA',
        config: defaultConfigForType('TABATA'),
      },
    })

    expect(wrapper.text()).toContain('Work (seconds)')
    expect(wrapper.text()).toContain('Rest (seconds)')
    expect(wrapper.text()).toContain('Cycles')
    expect(wrapper.findAll('input[type="number"]').length).toBe(4)
  })

  it('emits field updates', async () => {
    const wrapper = mount(WODDynamicConfigForm, {
      props: {
        type: 'AMRAP',
        config: defaultConfigForType('AMRAP'),
      },
    })

    await wrapper.find('input[type="number"]').setValue('1200')
    expect(wrapper.emitted('update:field')).toBeTruthy()
  })
})
