import { computed, ref } from 'vue'
import { createWOD } from '@/features/wod/api/wodApi'
import { defaultConfigForType, defaultStage, stageToPayload } from '@/features/wod/model/wodSchemas'
import type {
  CreateWODResponse,
  MovementInput,
  StageFormState,
  StageKind,
  WODType,
} from '@/features/wod/model/wodTypes'

export function useWODForm() {
  const name = ref('')
  const description = ref('')
  const stages = ref<StageFormState[]>([defaultStage()])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const result = ref<CreateWODResponse | null>(null)

  const payload = computed(() => ({
    name: name.value.trim(),
    description: description.value.trim(),
    stages: stages.value.map(stageToPayload),
  }))

  function addStage() {
    stages.value.push(defaultStage('METCON', 'AMRAP'))
  }

  function removeStage(index: number) {
    if (stages.value.length === 1) {
      return
    }
    stages.value.splice(index, 1)
  }

  function moveStageUp(index: number) {
    if (index === 0) {
      return
    }
    const next = [...stages.value]
    const current = next[index]
    next[index] = next[index - 1]
    next[index - 1] = current
    stages.value = next
  }

  function moveStageDown(index: number) {
    if (index >= stages.value.length - 1) {
      return
    }
    const next = [...stages.value]
    const current = next[index]
    next[index] = next[index + 1]
    next[index + 1] = current
    stages.value = next
  }

  function updateStageKind(index: number, kind: StageKind) {
    const next = [...stages.value]
    next[index] = { ...next[index], kind }
    stages.value = next
  }

  function updateStageType(index: number, type: WODType) {
    const next = [...stages.value]
    next[index] = { ...next[index], type, config: defaultConfigForType(type) }
    stages.value = next
  }

  function updateStageConfigField(index: number, key: string, value: number | undefined) {
    const next = [...stages.value]
    next[index] = {
      ...next[index],
      config: { ...next[index].config, [key]: value } as StageFormState['config'],
    }
    stages.value = next
  }

  function addMovement(stageIndex: number) {
    const next = [...stages.value]
    const stage = next[stageIndex]
    stage.movements.push({ position: stage.movements.length + 1, name: '', reps: 10 })
    stages.value = next
  }

  function removeMovement(stageIndex: number, movementIndex: number) {
    const next = [...stages.value]
    const stage = next[stageIndex]
    if (stage.movements.length === 1) {
      return
    }
    stage.movements.splice(movementIndex, 1)
    stages.value = next
  }

  function updateMovement(
    stageIndex: number,
    movementIndex: number,
    field: keyof MovementInput,
    value: string | number | undefined,
  ) {
    const next = [...stages.value]
    const movements = [...next[stageIndex].movements]
    movements[movementIndex] = { ...movements[movementIndex], [field]: value }
    next[stageIndex] = { ...next[stageIndex], movements }
    stages.value = next
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
    stages,
    loading,
    error,
    result,
    payload,
    addStage,
    removeStage,
    moveStageUp,
    moveStageDown,
    updateStageKind,
    updateStageType,
    updateStageConfigField,
    addMovement,
    removeMovement,
    updateMovement,
    submit,
  }
}
