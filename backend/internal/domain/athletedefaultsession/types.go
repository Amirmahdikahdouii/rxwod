package athletedefaultsession

import "time"

type AthleteDefaultSessionID string

type DayOfWeek int

type TimeSlot string

func (id AthleteDefaultSessionID) String() string {
	return string(id)
}

func validateDayOfWeek(day DayOfWeek) error {
	if day < 0 || day > 6 {
		return ErrInvalidDayOfWeek
	}
	return nil
}

func validateTimeSlot(slot TimeSlot) error {
	if _, err := time.Parse("15:04", string(slot)); err != nil {
		return ErrInvalidTimeSlot
	}
	return nil
}
