package classsession

import "time"

type CreateClassSessionCommand struct {
	WodID     *string
	StartTime time.Time
	EndTime   time.Time
	Capacity  int
}

type ListClassSessionsCommand struct {
	From time.Time
	To   time.Time
}
