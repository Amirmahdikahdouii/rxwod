export type WODType = 'AMRAP' | 'FORTIME' | 'TABATA' | 'EMOM' | 'OPEN'

export type StageKind = 'WARMUP' | 'STRENGTH' | 'CORE' | 'METCON' | 'COOLDOWN'

export type StageFormat = 'OPEN' | 'STRUCTURED'

export type ScoringKind = 'ROUNDS_REPS' | 'TIME_TO_COMPLETE' | 'TOTAL_REPS' | 'NONE'

export type WODStatus = 'DRAFT' | 'PUBLISHED' | 'ARCHIVED'

export type LoadUnit = 'kg' | 'lb' | 'bodyweight'

export interface MovementInput {
  position: number
  label?: string
  name: string
  prescription?: string
  sets?: number
  reps?: number
  loadValue?: number
  loadUnit?: LoadUnit
  notes?: string
}

export type WODFormConfig =
  | { type: 'OPEN' }
  | { type: 'AMRAP'; timeCapSeconds: number }
  | { type: 'FORTIME'; rounds: number; timeCapSeconds?: number }
  | { type: 'TABATA'; workSeconds: number; restSeconds: number; rounds: number; cycles: number }
  | { type: 'EMOM'; intervalSeconds: number; rounds: number }

export interface StageFormState {
  kind: StageKind
  type: WODType
  instructions: string
  config: WODFormConfig
  movements: MovementInput[]
}

export interface StagePayload {
  kind: StageKind
  type: WODType
  instructions: string
  config: Record<string, number | undefined>
  movements: MovementInput[]
}

export interface CreateWODPayload {
  name: string
  description: string
  scheduledDate?: string
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
  scheduledDate?: string
}

export interface WODSummary {
  id: string
  name: string
  status: WODStatus
  stageCount: number
  stages: StageSummary[]
  createdBy: string
  scheduledDate?: string
  publishedAt?: string
  createdAt: string
  updatedAt: string
}

export interface StageDetail extends StageSummary {
  id: string
  instructions: string
  config: Record<string, number | undefined>
  movements: Array<MovementInput & { id: string }>
}

export interface WODDetail {
  id: string
  name: string
  description: string
  status: WODStatus
  stages: StageDetail[]
  createdBy: string
  scheduledDate?: string
  publishedAt?: string
  createdAt: string
  updatedAt: string
}

export interface CalendarPlanSummary {
  id: string
  name: string
  status: WODStatus
}

export interface CalendarDaySummary {
  date: string
  publishedCount: number
  draftCount: number
  plans: CalendarPlanSummary[]
}

export const WOD_TYPES: WODType[] = ['OPEN', 'AMRAP', 'FORTIME', 'TABATA', 'EMOM']

export const STRUCTURED_WOD_TYPES: WODType[] = ['AMRAP', 'FORTIME', 'TABATA', 'EMOM']

export const STAGE_KINDS: StageKind[] = ['WARMUP', 'STRENGTH', 'CORE', 'METCON', 'COOLDOWN']

export const SCORING_BY_TYPE: Record<WODType, ScoringKind> = {
  OPEN: 'NONE',
  AMRAP: 'ROUNDS_REPS',
  FORTIME: 'TIME_TO_COMPLETE',
  TABATA: 'TOTAL_REPS',
  EMOM: 'ROUNDS_REPS',
}

export function isOpenFormat(type: WODType): boolean {
  return type === 'OPEN'
}

export function stageFormat(type: WODType): StageFormat {
  return isOpenFormat(type) ? 'OPEN' : 'STRUCTURED'
}

export function defaultTypeForKind(kind: StageKind): WODType {
  return kind === 'METCON' ? 'AMRAP' : 'OPEN'
}
