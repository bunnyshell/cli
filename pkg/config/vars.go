package config

import (
	"errors"
	"time"
)

const (
	defaultFormat  = "stylish"
	defaultTimeout = 30 * time.Second

	configDirPerm  = int(0o700)
	configFilePerm = int(0o600)
)

var (
	MainManager = NewManager()

	Formats = []string{
		"stylish",
		"json",
		"yaml",
	}
	FormatDescriptions = []string{
		"stylish\tOutput format for human consumption",
		"json\tOutput in JSON",
		"yaml\tOutput in YAML",
	}

	ErrConfigExists     = errors.New("configFile already exists")
	ErrUnknownProfile   = errors.New("profile not found")
	ErrDuplicateProfile = errors.New("profile already exists")
	ErrConfigLoad       = errors.New("unable to load config")
	ErrInvalidValue     = errors.New("invalid value")
)

func GetConfig() *Config {
	return MainManager.config
}

func GetOptions() *Options {
	return MainManager.options
}

func GetSettings() *Settings {
	return MainManager.settings
}
