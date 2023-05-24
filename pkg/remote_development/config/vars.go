package config

import (
	"errors"
)

const (
	defaultConfigDirParam  = ".../.bunnyshell"
	defaultConfigFileParam = "rdev.yaml"
)

var (
	MainManager = NewManager()

	ErrUnknownProfile = errors.New("profile not found")
	ErrConfigLoad     = errors.New("unable to load config")
	ErrInvalidValue   = errors.New("invalid value")
)
