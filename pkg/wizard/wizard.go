package wizard

import (
	"context"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type Wizard struct {
	profile *config.Profile
	client  *sdk.APIClient
}

func New(profile *config.Profile) *Wizard {
	return &Wizard{
		profile: profile,
		client:  lib.GetAPIFromProfile(*profile),
	}
}

func (w *Wizard) getContext() (context.Context, context.CancelFunc) {
	return lib.GetContextFromProfile(*w.profile)
}
