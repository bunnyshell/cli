package util

import (
	"errors"
	"time"
)

var (
	defaultDuration = 100 * time.Millisecond

	ErrInvalidValue = errors.New("invalid value")
)
