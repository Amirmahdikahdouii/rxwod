import type { MovementInput, StageFormState, StageKind, StagePayload, WODFormConfig, WODType } from './wodTypes'
import { defaultTypeForKind } from './wodTypes'

export interface FieldSchema {
  key: string
  label: string
  type: 'number'
  required?: boolean
  min?: number
}

export const WOD_CONFIG_SCHEMAS: Record<Exclude<WODType, 'OPEN'>, FieldSchema[]> = {
  AMRAP: [
    { key: 'timeCapSeconds', label: 'Time Cap (seconds)', type: 'number', required: true, min: 1 },
  ],
  FORTIME: [
    { key: 'rounds', label: 'Rounds', type: 'number', required: true, min: 1 },
    { key: 'timeCapSeconds', label: 'Time Cap (seconds)', type: 'number', min: 1 },
  ],
  TABATA: [
    { key: 'workSeconds', label: 'Work (seconds)', type: 'number', required: true, min: 1 },
    { key: 'restSeconds', label: 'Rest (seconds)', type: 'number', required: true, min: 0 },
    { key: 'rounds', label: 'Rounds', type: 'number', required: true, min: 1 },
    { key: 'cycles', label: 'Cycles', type: 'number', required: true, min: 1 },
  ],
  EMOM: [
    { key: 'intervalSeconds', label: 'Interval (seconds)', type: 'number', required: true, min: 1 },
    { key: 'rounds', label: 'Rounds', type: 'number', required: true, min: 1 },
  ],
}

export function defaultMovement(): MovementInput {
  return { position: 1, label: '', name: '', prescription: '', reps: undefined }
}

export function defaultConfigForType(type: WODType): WODFormConfig {
  switch (type) {
    case 'OPEN':
      return { type: 'OPEN' }
    case 'AMRAP':
      return { type: 'AMRAP', timeCapSeconds: 900 }
    case 'FORTIME':
      return { type: 'FORTIME', rounds: 5 }
    case 'TABATA':
      return { type: 'TABATA', workSeconds: 20, restSeconds: 10, rounds: 8, cycles: 1 }
    case 'EMOM':
      return { type: 'EMOM', intervalSeconds: 60, rounds: 10 }
  }
}

export function configToPayload(config: WODFormConfig): Record<string, number | undefined> {
  switch (config.type) {
    case 'OPEN':
      return {}
    case 'AMRAP':
      return { timeCapSeconds: config.timeCapSeconds }
    case 'FORTIME':
      return { rounds: config.rounds, timeCapSeconds: config.timeCapSeconds }
    case 'TABATA':
      return {
        workSeconds: config.workSeconds,
        restSeconds: config.restSeconds,
        rounds: config.rounds,
        cycles: config.cycles,
      }
    case 'EMOM':
      return { intervalSeconds: config.intervalSeconds, rounds: config.rounds }
  }
}

export function defaultStage(kind: StageKind = 'WARMUP', type?: WODType): StageFormState {
  const resolvedType = type ?? defaultTypeForKind(kind)
  return {
    kind,
    type: resolvedType,
    instructions: '',
    config: defaultConfigForType(resolvedType),
    movements: [defaultMovement()],
  }
}

export function stageToPayload(stage: StageFormState): StagePayload {
  return {
    kind: stage.kind,
    type: stage.type,
    instructions: stage.instructions.trim(),
    config: configToPayload(stage.config),
    movements: stage.movements.map((movement, index) => ({
      ...movement,
      position: index + 1,
      label: movement.label?.trim() ?? '',
      name: movement.name.trim(),
      prescription: movement.prescription?.trim() ?? '',
      notes: movement.notes?.trim() ?? '',
    })),
  }
}
