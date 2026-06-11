<script setup lang="ts">
import MovementListEditor from '@/features/wod/components/MovementListEditor.vue'
import ScorePreview from '@/features/wod/components/ScorePreview.vue'
import WODDynamicConfigForm from '@/features/wod/components/WODDynamicConfigForm.vue'
import WODTypeSelector from '@/features/wod/components/WODTypeSelector.vue'
import { useWODForm } from '@/features/wod/composables/useWODForm'
import BaseInput from '@/shared/components/BaseInput.vue'
import type { MovementInput } from '@/features/wod/model/wodTypes'

const {
  name,
  description,
  type,
  config,
  movements,
  loading,
  error,
  result,
  addMovement,
  removeMovement,
  updateConfigField,
  submit,
} = useWODForm()

function updateMovement(index: number, field: keyof MovementInput, value: string | number | undefined) {
  const next = [...movements.value]
  next[index] = { ...next[index], [field]: value }
  movements.value = next
}
</script>

<template>
  <div class="container stack-lg">
    <header class="page-header">
      <h1 class="page-title">Create WOD</h1>
      <p class="page-subtitle">Build your workout — configure type, scoring, and movements.</p>
    </header>

    <form class="card stack-lg" @submit.prevent="submit">
      <section class="card-section stack">
        <h2 class="section-title">Basics</h2>
        <BaseInput v-model="name" label="Name" placeholder="Open 24.1 Style AMRAP" />
        <BaseInput v-model="description" label="Description" placeholder="Workout notes" />
      </section>

      <section class="card-section stack">
        <h2 class="section-title">Workout Type</h2>
        <div class="type-config-grid">
          <WODTypeSelector v-model="type" />
          <ScorePreview :type="type" />
        </div>
      </section>

      <section class="card-section">
        <WODDynamicConfigForm :type="type" :config="config" @update:field="updateConfigField" />
      </section>

      <section class="card-section">
        <MovementListEditor
          :movements="movements"
          @add="addMovement"
          @remove="removeMovement"
          @update:movement="updateMovement"
        />
      </section>

      <button type="submit" class="btn-full" :disabled="loading">
        {{ loading ? 'Saving...' : 'Save WOD' }}
      </button>

      <div v-if="error" class="alert alert--error" role="alert">{{ error }}</div>
      <div v-if="result" class="alert alert--success" role="status">
        Saved {{ result.name }} ({{ result.type }}) with scoring {{ result.scoringKind }}.
      </div>
    </form>
  </div>
</template>
