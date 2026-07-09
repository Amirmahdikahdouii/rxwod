import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import BasePagination from '@/shared/components/BasePagination.vue'

describe('BasePagination', () => {
  it('hides itself when totalPages is 1', () => {
    const wrapper = mount(BasePagination, {
      props: {
        page: 1,
        totalPages: 1,
      },
    })

    expect(wrapper.find('nav').exists()).toBe(false)
  })

  it('emits update:page when clicking next', async () => {
    const wrapper = mount(BasePagination, {
      props: {
        page: 2,
        totalPages: 5,
      },
    })

    await wrapper.get('[aria-label="Next page"]').trigger('click')

    expect(wrapper.emitted('update:page')).toEqual([[3]])
  })

  it('emits update:page when clicking previous', async () => {
    const wrapper = mount(BasePagination, {
      props: {
        page: 3,
        totalPages: 5,
      },
    })

    await wrapper.get('[aria-label="Previous page"]').trigger('click')

    expect(wrapper.emitted('update:page')).toEqual([[2]])
  })

  it('emits update:page when clicking a page number', async () => {
    const wrapper = mount(BasePagination, {
      props: {
        page: 1,
        totalPages: 5,
      },
    })

    const pageButtons = wrapper.findAll('.pagination__page')
    await pageButtons[pageButtons.length - 1].trigger('click')

    expect(wrapper.emitted('update:page')).toEqual([[5]])
  })

  it('disables previous on the first page and next on the last page', () => {
    const wrapper = mount(BasePagination, {
      props: {
        page: 1,
        totalPages: 3,
      },
    })

    expect(wrapper.get('[aria-label="Previous page"]').attributes('disabled')).toBeDefined()
    expect(wrapper.get('[aria-label="Next page"]').attributes('disabled')).toBeUndefined()
  })

  it('shows range summary when total and limit are provided', () => {
    const wrapper = mount(BasePagination, {
      props: {
        page: 2,
        totalPages: 3,
        total: 42,
        limit: 20,
      },
    })

    expect(wrapper.text()).toContain('Showing 21-40 of 42')
  })
})
