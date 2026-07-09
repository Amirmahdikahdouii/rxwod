package wodresult

type WODResultID string
type ScoreValue int
type Notes string

func (id WODResultID) String() string {
	return string(id)
}

func validateScoreValue(value ScoreValue) error {
	if value < 0 {
		return ErrInvalidScore
	}
	return nil
}
