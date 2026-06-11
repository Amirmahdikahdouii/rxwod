package wod

import "errors"

var (
	ErrMissingConfigField = errors.New("required config field is missing for wod type")
)
