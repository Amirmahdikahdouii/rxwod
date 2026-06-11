export type WODType = 'AMRAP' | 'FORTIME' | 'TABATA' | 'EMOM'

export type StageKind = 'WARMUP' | 'STRENGTH' | 'CORE' | 'METCON' | 'COOLDOWN'

export type ScoringKind = 'ROUNDS_REPS' | 'TIME_TO_COMPLETE' | 'TOTAL_REPS'

export type WODStatus = 'DRAFT' | 'PUBLISHED' | 'ARCHIVED'

export type LoadUnit = 'kg' | 'lb' | 'bodyweight'

export interface MovementInput {
  position: number
  name: string
  reps?: number
  loadValue?: number
  loadUnit?: LoadUnit
  notes?: string
}

export type WODFormConfig =
  | { type: 'AMRAP'; timeCapSeconds: number }
  | { type: 'FORTIME'; rounds: number; timeCapSeconds?: number }
  | { type: 'TABATA'; workSeconds: number; restSeconds: number; rounds: number; cycles: number }
  | { type: 'EMOM'; intervalSeconds: number; rounds: number }

export interface StageFormState {
  kind: StageKind
  type: WODType
  config: WODFormConfig
  movements: MovementInput[]
}

export interface StagePayload {
  kind: StageKind
  type: WODType
  config: Record<string, number | undefined>
  movements: MovementInput[]
}

export interface CreateWODPayload {
  name: string
  description: string
  stages: StagePayload[]
}

export interface StageSummary {
  kind: StageKind
  position: number
  type: WODType
  scoringKind: ScoringKind
}

export interface CreateWODResponse {
  id: string
  name: string
  status: WODStatus
  stageCount: number
  stages: StageSummary[]
}

export interface WODSummary {
  id: string
  name: string
  status: WODStatus
  stageCount: number
  stages: StageSummary[]
  createdAt: string
  updatedAt: string
}

export interface StageDetail extends StageSummary {
  id: string
  config: Record<string, number | undefined>
  movements: Array<MovementInput & { id: string }>
}

export interface WODDetail {
  id: string
  name: string
  description: string
  status: WODStatus
  stages: StageDetail[]
  createdAt: string
  updatedAt: string
}

export const WOD_TYPES: WODType[] = ['AMRAP', 'FORTIME', 'TABATA', 'EMOM']

export const STAGE_KINDS: StageKind[] = ['WARMUP', 'STRENGTH', 'CORE', 'METCON', 'COOLDOWN']

export const SCORING_BY_TYPE: Record<WODType, ScoringKind> = {
  AMRAP: 'ROUNDS_REPS',
  FORTIME: 'TIME_TO_COMPLETE',
  TABATA: 'TOTAL_REPS',
  EMOM: 'ROUNDS_REPS',
}
