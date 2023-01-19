package util

import (
	"errors"
	"time"
)

const (
	StrTrue  = "true"
	StrFalse = "false"
)

var (
	defaultDuration = 100 * time.Millisecond

	ErrInvalidValue = errors.New("invalid value")
)
