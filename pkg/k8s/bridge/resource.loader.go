package bridge

import (
	"fmt"
	"strings"

	"bunnyshell.com/cli/pkg/interactive"
	k8sWizard "bunnyshell.com/cli/pkg/wizard/k8s"
	"bunnyshell.com/sdk"
)

type ResourceLoader struct {
	EnvironmentComponent

	ManualSelectSingleResource bool

	resources []sdk.ComponentResourceItem

	selectedIndex int
}

func NewResourceLoader() *ResourceLoader {
	return &ResourceLoader{
		EnvironmentComponent: *NewEnvironmentComponent(),

		ManualSelectSingleResource: false,

		selectedIndex: -1,
	}
}

func (loader *ResourceLoader) LoadResources() error {
	return loader.ensureResources()
}

func (loader *ResourceLoader) CountResources() int {
	if loader.resources == nil {
		return 0
	}

	return len(loader.resources)
}

func (loader *ResourceLoader) GetResource() *sdk.ComponentResourceItem {
	if loader.selectedIndex == -1 {
		return nil
	}

	return &loader.resources[loader.selectedIndex]
}

func (loader *ResourceLoader) SelectResourceFromString(spec string) error {
	resourceSpec := NewResourceSpec(strings.ToLower(spec))
	if resourceSpec == nil {
		return fmt.Errorf("%w: %s", ErrInvalidResourceSpec, spec)
	}

	return loader.SelectResourceFromSpec(resourceSpec)
}

func (loader *ResourceLoader) SelectResourceFromSpec(spec *ResourceSpec) error {
	if err := loader.ensureResources(); err != nil {
		return err
	}

	for index, resourceItem := range loader.resources {
		if spec.Match(resourceItem) {
			loader.selectedIndex = index

			return nil
		}
	}

	return fmt.Errorf("%w (no match)", ErrNoComponentResources)
}

func (loader *ResourceLoader) SelectResource() error {
	if err := loader.ensureResources(); err != nil {
		return err
	}

	if len(loader.resources) == 1 && !loader.ManualSelectSingleResource {
		loader.selectedIndex = 0

		return nil
	}

	index, _, err := interactive.Choose("Select resource", resourcesToNames(loader.resources))
	if err != nil {
		return err
	}

	loader.selectedIndex = index

	return nil
}

func (loader *ResourceLoader) ensureResources() error {
	if !loader.EnvironmentComponent.IsLoaded() {
		return ErrNotLoaded
	}

	if loader.resources != nil {
		return nil
	}

	resources, err := k8sWizard.RDevList(k8sWizard.NewRDevListOptions(loader.Component.GetId()))
	if err != nil {
		return err
	}

	if len(resources) == 0 {
		return ErrNoComponentResources
	}

	loader.resources = resources

	return nil
}
