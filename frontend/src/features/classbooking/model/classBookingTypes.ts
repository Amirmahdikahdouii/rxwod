export interface ClassSessionResponse {
  id: string
  gymId: string
  wodId?: string
  startTime: string
  endTime: string
  capacity: number
  coachId: string
  bookedCount: number
  myBookingStatus?: 'BOOKED' | 'CANCELLED'
  createdAt: string
  updatedAt: string
}

export interface CreateClassSessionResponse {
  session: ClassSessionResponse
  autoBookedCount: number
  autoBookedMembershipIds: string[]
}

export interface ClassBookingResponse {
  id: string
  sessionId: string
  gymMembershipId: string
  status: 'BOOKED' | 'CANCELLED' | 'ATTENDED'
  createdAt: string
  updatedAt: string
}

export interface DefaultSessionResponse {
  id: string
  gymMembershipId: string
  dayOfWeek: number
  timeSlot: string
  createdAt: string
  updatedAt: string
}

export interface CreateClassSessionPayload {
  wodId?: string
  startTime: string
  endTime: string
  capacity: number
}

export interface SetDefaultSessionPayload {
  dayOfWeek: number
  timeSlot: string
}

export const DAY_OF_WEEK_LABELS = [
  'Sunday',
  'Monday',
  'Tuesday',
  'Wednesday',
  'Thursday',
  'Friday',
  'Saturday',
] as const
