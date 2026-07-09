package wodresult

import "errors"

var (
	ErrInvalidScore = errors.New("score value must be zero or greater")
	ErrNotFound     = errors.New("wod result not found")
)
