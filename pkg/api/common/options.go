package common

import "bunnyshell.com/cli/pkg/config"

type Options struct {
	Profile *config.Profile
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) GetProfile() config.Profile {
	if o.Profile == nil {
		return config.GetSettings().Profile
	}

	return *o.Profile
}
