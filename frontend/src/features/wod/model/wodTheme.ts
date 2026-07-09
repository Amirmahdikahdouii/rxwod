import type { StageKind, WODType } from './wodTypes'

export const WOD_TYPE_BADGE_CLASS: Record<WODType, string> = {
  OPEN: 'badge-open',
  AMRAP: 'badge-amrap',
  FORTIME: 'badge-fortime',
  TABATA: 'badge-tabata',
  EMOM: 'badge-emom',
}

export const STAGE_KIND_BADGE_CLASS: Record<StageKind, string> = {
  WARMUP: 'badge-stage-warmup',
  STRENGTH: 'badge-stage-strength',
  CORE: 'badge-stage-core',
  METCON: 'badge-stage-metcon',
  COOLDOWN: 'badge-stage-cooldown',
}

export const STAGE_KIND_LABEL: Record<StageKind, string> = {
  WARMUP: 'Warm-up',
  STRENGTH: 'Strength',
  CORE: 'Core',
  METCON: 'Metcon',
  COOLDOWN: 'Cooldown',
}

export function stageDisplayLabel(stage: { kind: StageKind; type: WODType }): string {
  return `${STAGE_KIND_LABEL[stage.kind]} (${stage.type})`
}
