export type WODType = 'AMRAP' | 'FORTIME' | 'TABATA' | 'EMOM'

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

export interface CreateWODPayload {
  name: string
  type: WODType
  description: string
  config: Record<string, number | undefined>
  movements: MovementInput[]
}

export interface CreateWODResponse {
  id: string
  name: string
  type: WODType
  status: WODStatus
  scoringKind: ScoringKind
}

export interface WODSummary {
  id: string
  name: string
  type: WODType
  status: WODStatus
  scoringKind: ScoringKind
  createdAt: string
  updatedAt: string
}

export interface WODDetail extends WODSummary {
  description: string
  config: Record<string, number | undefined>
  movements: Array<MovementInput & { id: string }>
}

export const WOD_TYPES: WODType[] = ['AMRAP', 'FORTIME', 'TABATA', 'EMOM']

export const SCORING_BY_TYPE: Record<WODType, ScoringKind> = {
  AMRAP: 'ROUNDS_REPS',
  FORTIME: 'TIME_TO_COMPLETE',
  TABATA: 'TOTAL_REPS',
  EMOM: 'ROUNDS_REPS',
}
