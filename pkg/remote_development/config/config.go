package config

import "fmt"

type Config struct {
	Profiles NamedProfiles `json:"profiles" yaml:"profiles"`
}

func (config *Config) getProfile(name string) (*Profile, error) {
	profile, ok := config.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownProfile, name)
	}

	return &profile, nil
}

func (config *Config) profileNames() []string {
	profiles := []string{}

	for name := range config.Profiles {
		profiles = append(profiles, name)
	}

	return profiles
}
