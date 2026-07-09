<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  page: number
  totalPages: number
  total?: number
  limit?: number
  disabled?: boolean
}>()

const emit = defineEmits<{
  'update:page': [page: number]
}>()

const pageNumbers = computed(() => {
  const { page, totalPages } = props
  if (totalPages <= 1) {
    return [] as Array<number | 'ellipsis'>
  }

  const pages = new Set<number>()
  pages.add(1)
  pages.add(totalPages)
  pages.add(page)
  if (page > 1) {
    pages.add(page - 1)
  }
  if (page < totalPages) {
    pages.add(page + 1)
  }

  const sorted = [...pages].sort((a, b) => a - b)
  const result: Array<number | 'ellipsis'> = []

  for (let i = 0; i < sorted.length; i++) {
    const current = sorted[i]
    const previous = sorted[i - 1]
    if (i > 0 && current - previous > 1) {
      result.push('ellipsis')
    }
    result.push(current)
  }

  return result
})

const rangeSummary = computed(() => {
  if (props.total === undefined || props.limit === undefined || props.total === 0) {
    return null
  }

  const start = (props.page - 1) * props.limit + 1
  const end = Math.min(props.page * props.limit, props.total)
  return `Showing ${start}-${end} of ${props.total}`
})

function goTo(targetPage: number) {
  if (
    props.disabled ||
    targetPage < 1 ||
    targetPage > props.totalPages ||
    targetPage === props.page
  ) {
    return
  }

  emit('update:page', targetPage)
}
</script>

<template>
  <nav v-if="totalPages > 1" class="pagination" aria-label="Pagination">
    <p v-if="rangeSummary" class="pagination__summary helper-text">{{ rangeSummary }}</p>

    <div class="pagination__controls row row--align-center">
      <button
        type="button"
        class="secondary compact-button pagination__button"
        :disabled="disabled || page <= 1"
        aria-label="Previous page"
        @click="goTo(page - 1)"
      >
        Previous
      </button>

      <div class="pagination__pages" role="group" aria-label="Page numbers">
        <template v-for="(item, index) in pageNumbers" :key="`${item}-${index}`">
          <span v-if="item === 'ellipsis'" class="pagination__ellipsis" aria-hidden="true">...</span>
          <button
            v-else
            type="button"
            class="pagination__page"
            :class="{ 'pagination__page--active': item === page }"
            :disabled="disabled || item === page"
            :aria-current="item === page ? 'page' : undefined"
            @click="goTo(item)"
          >
            {{ item }}
          </button>
        </template>
      </div>

      <button
        type="button"
        class="secondary compact-button pagination__button"
        :disabled="disabled || page >= totalPages"
        aria-label="Next page"
        @click="goTo(page + 1)"
      >
        Next
      </button>
    </div>
  </nav>
</template>
