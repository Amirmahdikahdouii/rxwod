package wodresult

import "time"

type ResultDTO struct {
	ID              string
	WODID           string
	GymMembershipID string
	ScoreValue      int
	IsRx            bool
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type LeaderboardEntryDTO struct {
	Rank            int
	GymMembershipID string
	DisplayName     string
	ScoreValue      int
	IsRx            bool
	Notes           string
	UpdatedAt       time.Time
}

type LeaderboardDTO struct {
	WODID   string
	Entries []LeaderboardEntryDTO
}
