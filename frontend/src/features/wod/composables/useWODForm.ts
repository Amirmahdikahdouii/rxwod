import { computed, ref, watch } from 'vue'
import { createWOD } from '@/features/wod/api/wodApi'
import { configToPayload, defaultConfigForType } from '@/features/wod/model/wodSchemas'
import type { CreateWODResponse, MovementInput, WODFormConfig, WODType } from '@/features/wod/model/wodTypes'

export function useWODForm(initialType: WODType = 'AMRAP') {
  const name = ref('')
  const description = ref('')
  const type = ref<WODType>(initialType)
  const config = ref<WODFormConfig>(defaultConfigForType(initialType))
  const movements = ref<MovementInput[]>([
    { position: 1, name: '', reps: 10 },
  ])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const result = ref<CreateWODResponse | null>(null)

  watch(type, (nextType) => {
    config.value = defaultConfigForType(nextType)
  })

  const payload = computed(() => ({
    name: name.value.trim(),
    type: type.value,
    description: description.value.trim(),
    config: configToPayload(config.value),
    movements: movements.value.map((movement, index) => ({
      ...movement,
      position: index + 1,
      name: movement.name.trim(),
    })),
  }))

  function addMovement() {
    movements.value.push({ position: movements.value.length + 1, name: '', reps: 10 })
  }

  function removeMovement(index: number) {
    if (movements.value.length === 1) {
      return
    }
    movements.value.splice(index, 1)
  }

  function updateConfigField(key: string, value: number | undefined) {
    config.value = { ...config.value, [key]: value } as WODFormConfig
  }

  async function submit() {
    loading.value = true
    error.value = null
    result.value = null

    const response = await createWOD(payload.value)
    loading.value = false

    if (!response.ok) {
      error.value = response.error
      return
    }

    result.value = response.value
  }

  return {
    name,
    description,
    type,
    config,
    movements,
    loading,
    error,
    result,
    payload,
    addMovement,
    removeMovement,
    updateConfigField,
    submit,
  }
}
