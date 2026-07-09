package classbooking

type BookCommand struct {
	SessionID string
}

type OverbookCommand struct {
	SessionID     string
	AthleteUserID string
}

type CancelCommand struct {
	SessionID     string
	AthleteUserID *string
}
