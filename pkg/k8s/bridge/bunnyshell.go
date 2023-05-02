package bridge

import (
	"fmt"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/wizard"
	"bunnyshell.com/sdk"
)

type EnvironmentComponent struct {
	loaded bool

	Environment *sdk.EnvironmentItem
	Component   *sdk.ComponentItem
}

func NewEnvironmentComponent() *EnvironmentComponent {
	return &EnvironmentComponent{
		loaded: false,
	}
}

func (ec *EnvironmentComponent) IsLoaded() bool {
	return ec.loaded
}

func (ec *EnvironmentComponent) Load(profile config.Profile) error {
	if err := ec.loadProfile(&profile); err != nil {
		return err
	}

	ec.loaded = true

	return nil
}

func (ec *EnvironmentComponent) loadProfile(profile *config.Profile) error {
	environment, component, err := getEnvironmentComponent(profile)
	if err != nil {
		return err
	}

	ec.Environment = environment
	ec.Component = component

	return nil
}

func getEnvironmentComponent(profile *config.Profile) (*sdk.EnvironmentItem, *sdk.ComponentItem, error) {
	wiz := wizard.New(profile)

	if wiz.HasComponent() {
		comp, err := wiz.GetComponent()
		if err != nil {
			return nil, nil, err
		}

		env, err := wiz.GetEnvironment()
		if err != nil {
			return nil, nil, err
		}

		return env, comp, nil
	}

	env, err := getOrLoadEnvironment(wiz)
	if err != nil {
		return nil, nil, err
	}

	comp, err := wiz.GetComponent()
	if err != nil {
		return nil, nil, err
	}

	return env, comp, nil
}

func getOrLoadEnvironment(wiz *wizard.Wizard) (*sdk.EnvironmentItem, error) {
	if wiz.HasEnvironment() {
		return wiz.GetEnvironment()
	}

	if !wiz.HasProject() {
		if _, err := wiz.GetOrganization(); err != nil {
			return nil, err
		}
	}

	if _, err := wiz.GetProject(); err != nil {
		return nil, err
	}

	return wiz.GetEnvironment()
}

func resourcesToNames(resources []sdk.ComponentResourceItem) []string {
	names := make([]string, len(resources))

	for i, resource := range resources {
		names[i] = fmt.Sprintf(
			"%s / %s / %s",
			resource.GetNamespace(),
			resource.GetKind(),
			resource.GetName(),
		)
	}

	return names
}
