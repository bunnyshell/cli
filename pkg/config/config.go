package config

import (
	"time"
)

type Config struct {
	Debug bool `json:"debug" yaml:"debug"`

	OutputFormat string        `json:"outputFormat,omitempty" yaml:"outputFormat,omitempty"`
	Timeout      time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	DefaultProfile string        `json:"defaultProfile,omitempty" yaml:"defaultProfile,omitempty"`
	Profiles       NamedProfiles `json:"profiles,omitempty" yaml:"profiles,omitempty"`
}

func (config *Config) setDefaultProfile(name string) error {
	if _, ok := config.Profiles[name]; !ok {
		return ErrUnknownProfile
	}

	config.DefaultProfile = name

	return nil
}

func (config *Config) getProfile(name string) (*Profile, error) {
	profile, ok := config.Profiles[name]
	if !ok {
		return nil, ErrUnknownProfile
	}

	return &profile, nil
}

func (config *Config) addProfile(profile Profile) error {
	name := profile.Name

	if config.Profiles == nil {
		config.Profiles = NamedProfiles{}
	} else if _, ok := config.Profiles[name]; ok {
		return ErrDuplicateProfile
	}

	config.Profiles[name] = profile

	return nil
}

func (config *Config) removeProfile(name string) error {
	delete(config.Profiles, name)

	if config.DefaultProfile == name {
		config.DefaultProfile = ""
	}

	return nil
}

func (config *Config) profileNames() []string {
	profiles := []string{}

	for name := range config.Profiles {
		profiles = append(profiles, name)
	}

	return profiles
}
