import { computed, ref } from 'vue'
import { createWOD, getWOD, publishWOD, updateWOD } from '@/features/wod/api/wodApi'
import {
  defaultConfigForType,
  defaultMovement,
  defaultStage,
  detailToFormState,
  stageToPayload,
} from '@/features/wod/model/wodSchemas'
import type {
  CreateWODResponse,
  MovementInput,
  StageFormState,
  StageFormat,
  StageKind,
  WODDetail,
  WODType,
} from '@/features/wod/model/wodTypes'
import { defaultTypeForKind } from '@/features/wod/model/wodTypes'

type WODFormMode = 'create' | 'edit'
type SaveAction = 'draft' | 'publish' | null

export function useWODForm(initialMode: WODFormMode = 'create') {
  const mode = ref<WODFormMode>(initialMode)
  const wodId = ref<string | null>(null)
  const name = ref('')
  const description = ref('')
  const scheduledDate = ref('')
  const stages = ref<StageFormState[]>([defaultStage()])
  const loading = ref(false)
  const savingAction = ref<SaveAction>(null)
  const initialLoading = ref(false)
  const error = ref<string | null>(null)
  const result = ref<CreateWODResponse | WODDetail | null>(null)
  const lastAction = ref<SaveAction>(null)

  const payload = computed(() => ({
    name: name.value.trim(),
    description: description.value.trim(),
    scheduledDate: scheduledDate.value || undefined,
    stages: stages.value.map(stageToPayload),
  }))

  const isDirty = computed(() =>
    Boolean(
      name.value.trim() ||
        description.value.trim() ||
        scheduledDate.value ||
        stages.value.some(
          (stage) =>
            stage.instructions.trim() ||
            stage.movements.some(
              (movement) => movement.name.trim() || movement.prescription?.trim() || movement.label?.trim(),
            ),
        ),
    ),
  )

  function addStage(kind: StageKind = 'METCON') {
    stages.value.push(defaultStage(kind))
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
    const stage = next[index]
    const type = defaultTypeForKind(kind)
    next[index] = {
      ...stage,
      kind,
      type,
      config: defaultConfigForType(type),
    }
    stages.value = next
  }

  function updateStageFormat(index: number, format: StageFormat) {
    const next = [...stages.value]
    const stage = next[index]
    if (format === 'OPEN') {
      next[index] = { ...stage, type: 'OPEN', config: defaultConfigForType('OPEN') }
    } else {
      const type: WODType = stage.kind === 'METCON' ? 'AMRAP' : 'FORTIME'
      next[index] = { ...stage, type, config: defaultConfigForType(type) }
    }
    stages.value = next
  }

  function updateStageInstructions(index: number, value: string) {
    const next = [...stages.value]
    next[index] = { ...next[index], instructions: value }
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
    stage.movements.push({ ...defaultMovement(), position: stage.movements.length + 1 })
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

  function loadFromDetail(detail: WODDetail) {
    const next = detailToFormState(detail)
    mode.value = 'edit'
    wodId.value = detail.id
    name.value = next.name
    description.value = detail.description
    scheduledDate.value = detail.scheduledDate ?? ''
    stages.value = next.stages.length > 0 ? next.stages : [defaultStage()]
    result.value = null
    error.value = null
  }

  function applyDraftState(draft: {
    name: string
    description: string
    scheduledDate: string
    stages: StageFormState[]
  }) {
    name.value = draft.name
    description.value = draft.description
    scheduledDate.value = draft.scheduledDate
    stages.value = draft.stages.length > 0 ? draft.stages : [defaultStage()]
    result.value = null
    error.value = null
  }

  async function initEdit(id: string) {
    mode.value = 'edit'
    wodId.value = id
    initialLoading.value = true
    error.value = null
    result.value = null

    const response = await getWOD(id)
    initialLoading.value = false

    if (!response.ok) {
      error.value = response.error
      return null
    }

    loadFromDetail(response.value)
    return response.value
  }

  async function saveDraft() {
    return persist('draft')
  }

  async function publishProgram() {
    if (!scheduledDate.value) {
      error.value = 'Program date is required before publishing.'
      return false
    }
    return persist('publish')
  }

  async function persist(action: SaveAction) {
    loading.value = true
    savingAction.value = action
    error.value = null
    result.value = null
    lastAction.value = action

    const saveResponse =
      mode.value === 'edit' && wodId.value
        ? await updateWOD(wodId.value, payload.value)
        : await createWOD(payload.value)

    if (!saveResponse.ok) {
      loading.value = false
      savingAction.value = null
      error.value = saveResponse.error
      return false
    }

    if (action === 'publish') {
      const targetId = mode.value === 'edit' && wodId.value ? wodId.value : saveResponse.value.id
      const publishResponse = await publishWOD(targetId)
      loading.value = false
      savingAction.value = null

      if (!publishResponse.ok) {
        error.value = publishResponse.error
        if (mode.value === 'create') {
          mode.value = 'edit'
          wodId.value = saveResponse.value.id
        }
        return false
      }

      result.value = publishResponse.value
      if (mode.value === 'create') {
        mode.value = 'edit'
        wodId.value = publishResponse.value.id
      }
      return true
    }

    loading.value = false
    savingAction.value = null
    result.value = saveResponse.value
    if (mode.value === 'create') {
      mode.value = 'edit'
      wodId.value = saveResponse.value.id
    }
    return true
  }

  return {
    mode,
    wodId,
    name,
    description,
    scheduledDate,
    stages,
    loading,
    savingAction,
    initialLoading,
    error,
    result,
    lastAction,
    isDirty,
    payload,
    addStage,
    removeStage,
    moveStageUp,
    moveStageDown,
    updateStageKind,
    updateStageFormat,
    updateStageInstructions,
    updateStageType,
    updateStageConfigField,
    addMovement,
    removeMovement,
    updateMovement,
    loadFromDetail,
    applyDraftState,
    initEdit,
    saveDraft,
    publishProgram,
  }
}
