<script setup lang="ts">
defineProps<{
  modelValue: number | undefined
  label: string
  min?: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: number | undefined]
}>()

function onInput(event: Event) {
  const raw = (event.target as HTMLInputElement).value
  emit('update:modelValue', raw === '' ? undefined : Number(raw))
}
</script>

<template>
  <div class="field">
    <label>{{ label }}</label>
    <input type="number" :value="modelValue" :min="min" @input="onInput" />
  </div>
</template>
