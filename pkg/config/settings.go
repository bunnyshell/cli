package config

import (
	"time"
)

type Settings struct {
	ConfigFile string

	Debug      bool
	NoProgress bool

	NonInteractive bool

	Profile Profile

	Verbosity int

	OutputFormat string
	Timeout      time.Duration
}

func NewSettings() *Settings {
	return &Settings{
		Timeout:      defaultTimeout,
		OutputFormat: defaultFormat,
	}
}

func (settings *Settings) IsStylish() bool {
	return settings.OutputFormat == "stylish"
}
