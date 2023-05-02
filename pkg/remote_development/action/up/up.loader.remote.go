package up

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/cli/pkg/k8s/bridge"
	"bunnyshell.com/sdk"
)

type ContainerMap = map[string]sdk.ContainerConfigItem

func (up *Options) getProfileFromComponentConfig() (*sdk.ContainerConfigItemProfile, error) {
	rdevConfig, err := component.RDevConfig(component.NewRDevConfigOptions(up.resourceLoader.Component.GetId()))
	if err != nil {
		return nil, err
	}

	containers, err := up.getContainers(rdevConfig)
	if err != nil {
		return nil, err
	}

	if containers == nil {
		return nil, ErrNoRDevConfig
	}

	return up.getProfileFromContainerMap(containers)
}

func (up *Options) getContainers(rdevConfig *sdk.ComponentConfigItem) (ContainerMap, error) {
	if !rdevConfig.HasConfig() {
		return nil, ErrNoRDevConfig
	}

	config := rdevConfig.GetConfig()

	if config.ArrayOfSimpleResourceConfigItem != nil {
		return up.getContainersFromSimpleResource(*rdevConfig.Config.ArrayOfSimpleResourceConfigItem)
	}

	if config.ArrayOfExtendedResourceConfigItem != nil {
		return up.selectExtendedResource(*rdevConfig.Config.ArrayOfExtendedResourceConfigItem)
	}

	return nil, ErrUnknownConfigurationType
}

func (up *Options) getContainersFromSimpleResource(resourceConfigList []sdk.SimpleResourceConfigItem) (ContainerMap, error) {
	if len(resourceConfigList) == 0 {
		return nil, fmt.Errorf("%w of configuration for component", ErrEmptyList)
	}

	if len(resourceConfigList) != 1 {
		return nil, ErrTooManySimpleConfig
	}

	if err := up.resourceLoader.LoadResources(); err != nil {
		return nil, err
	}

	if up.resourceLoader.CountResources() != 1 {
		return nil, ErrTooManySimpleResources
	}

	if up.resourcePath != "" {
		if err := up.resourceLoader.SelectResourceFromString(up.resourcePath); err != nil {
			return nil, err
		}
	} else {
		if err := up.resourceLoader.SelectResource(); err != nil {
			return nil, err
		}
	}

	return *resourceConfigList[0].Containers, nil
}

func (up *Options) getProfileFromContainerMap(containers map[string]sdk.ContainerConfigItem) (*sdk.ContainerConfigItemProfile, error) {
	name, err := up.selectContainerName(containers)
	if err != nil {
		return nil, err
	}

	up.containerName = *name

	return containers[*name].Profile, nil
}

func (up *Options) selectExtendedResource(resourceConfigList []sdk.ExtendedResourceConfigItem) (ContainerMap, error) {
	resource, err := up.selectResource(resourceConfigList)
	if err != nil {
		return nil, err
	}

	resourceSpec := &bridge.ResourceSpec{
		Namespace: resource.GetNamespace(),
		Kind:      resource.GetKind(),
		Name:      resource.GetName(),
	}

	if up.resourcePath != "" && !resourceSpec.MatchString(up.resourcePath) {
		return nil, bridge.ErrNoComponentResources
	}

	return *resource.Containers, nil
}

func (up *Options) selectResource(resources []sdk.ExtendedResourceConfigItem) (*sdk.ExtendedResourceConfigItem, error) {
	if len(resources) == 0 {
		return nil, fmt.Errorf("%w of configuration for component", ErrEmptyList)
	}

	if len(resources) == 1 && !up.ManualSelectSingleResource {
		return &resources[0], nil
	}

	index, _, err := interactive.Choose("Choose a resource", resourcesToNames(resources))
	if err != nil {
		return nil, err
	}

	return &resources[index], nil
}

func resourcesToNames(resources []sdk.ExtendedResourceConfigItem) []string {
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
