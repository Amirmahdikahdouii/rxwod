import type { ScoringKind } from '@/features/wod/model/wodTypes'

// The backend stores scoreValue as an opaque non-negative int. These are the
// frontend-only encode/decode conventions for each ScoringKind:
// - TIME_TO_COMPLETE: total seconds (minutes * 60 + seconds)
// - ROUNDS_REPS: rounds * ROUNDS_REPS_MULTIPLIER + reps
// - TOTAL_REPS: raw rep count
const ROUNDS_REPS_MULTIPLIER = 1000

export interface TimeToCompleteParts {
  minutes: number
  seconds: number
}

export interface RoundsRepsParts {
  rounds: number
  reps: number
}

export function encodeTimeToComplete(parts: TimeToCompleteParts): number {
  return Math.max(0, parts.minutes) * 60 + Math.max(0, parts.seconds)
}

export function encodeRoundsReps(parts: RoundsRepsParts): number {
  return Math.max(0, parts.rounds) * ROUNDS_REPS_MULTIPLIER + Math.max(0, parts.reps)
}

export function scoringKindLabel(kind: ScoringKind): string {
  switch (kind) {
    case 'TIME_TO_COMPLETE':
      return 'For Time (MM:SS)'
    case 'ROUNDS_REPS':
      return 'AMRAP (Rounds + Reps)'
    case 'TOTAL_REPS':
      return 'Total Reps'
    case 'NONE':
    default:
      return 'No structured scoring'
  }
}

export function decodeScore(kind: ScoringKind, value: number): string {
  switch (kind) {
    case 'TIME_TO_COMPLETE': {
      const minutes = Math.floor(value / 60)
      const seconds = value % 60
      return `${minutes}:${String(seconds).padStart(2, '0')}`
    }
    case 'ROUNDS_REPS': {
      const rounds = Math.floor(value / ROUNDS_REPS_MULTIPLIER)
      const reps = value % ROUNDS_REPS_MULTIPLIER
      return reps > 0 ? `${rounds} rounds + ${reps} reps` : `${rounds} rounds`
    }
    case 'TOTAL_REPS':
      return `${value} reps`
    case 'NONE':
    default:
      return '—'
  }
}
