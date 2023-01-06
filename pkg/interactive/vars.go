package interactive

import (
	"errors"

	"bunnyshell.com/cli/pkg/config"
)

var (
	ErrInvalidValue   = errors.New("invalid value")
	ErrRequiredValue  = errors.New("required value")
	ErrNonInteractive = errors.New("refusing to run with non-interactive flag")
)

func getSettings() *config.Settings {
	return config.GetSettings()
}

func IsStylishInteractive() bool {
	if getSettings().OutputFormat != "stylish" {
		return false
	}

	return !getSettings().NonInteractive
}
