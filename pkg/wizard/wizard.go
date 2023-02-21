package wizard

import (
	"bunnyshell.com/cli/pkg/config"
)

type Wizard struct {
	profile *config.Profile
}

func New(profile *config.Profile) *Wizard {
	return &Wizard{
		profile: profile,
	}
}
