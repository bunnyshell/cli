package config

import (
	"errors"
)

const RDevConfigFile = ".bunnyshell/rdev.yaml"

var (
	MainManager = NewManager()

	ErrUnknownProfile = errors.New("profile not found")
	ErrConfigLoad     = errors.New("unable to load config")
	ErrInvalidValue   = errors.New("invalid value")
)
