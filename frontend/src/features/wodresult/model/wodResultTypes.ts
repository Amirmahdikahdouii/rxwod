export interface SubmitResultPayload {
  scoreValue: number
  isRx: boolean
  notes: string
}

export interface WODResultResponse {
  id: string
  wodId: string
  gymMembershipId: string
  scoreValue: number
  isRx: boolean
  notes: string
  createdAt: string
  updatedAt: string
}

export interface LeaderboardEntry {
  rank: number
  gymMembershipId: string
  displayName: string
  scoreValue: number
  isRx: boolean
  notes: string
  updatedAt: string
}

export interface LeaderboardResponse {
  wodId: string
  entries: LeaderboardEntry[]
}
