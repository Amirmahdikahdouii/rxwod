<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: 'default' | 'elevated'
    padding?: 'sm' | 'md' | 'lg'
    hoverable?: boolean
    disabled?: boolean
  }>(),
  {
    variant: 'default',
    padding: 'md',
    hoverable: false,
    disabled: false,
  },
)

defineEmits<{
  click: [event: MouseEvent]
}>()
</script>

<template>
  <component
    :is="hoverable ? 'button' : 'div'"
    class="ui-card"
    :class="[
      `ui-card--${variant}`,
      `ui-card--padding-${padding}`,
      { 'ui-card--hoverable': hoverable },
    ]"
    :type="hoverable ? 'button' : undefined"
    :disabled="hoverable ? disabled : undefined"
    @click="hoverable ? $emit('click', $event) : undefined"
  >
    <slot />
  </component>
</template>

<style src="./ui.css"></style>
