package k8s

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/sdk"
)

type PodListOptions struct {
	Component string
}

func PodList(options *PodListOptions) ([]sdk.ComponentResourceItem, error) {
	resources, err := component.Resources(component.NewResourceOptions(options.Component))
	if err != nil {
		return nil, err
	}

	pods := []sdk.ComponentResourceItem{}

	for _, resource := range resources {
		if resource.GetKind() == "Pod" {
			pods = append(pods, resource)
		}
	}

	return pods, nil
}

func PodSelect(options *PodListOptions) (*sdk.ComponentResourceItem, error) {
	pods, err := PodList(options)
	if err != nil {
		return nil, err
	}

	if len(pods) == 0 {
		return nil, errEmptyList
	}

	if len(pods) == 1 {
		return &pods[0], nil
	}

	index, _, err := interactive.Choose("Choose a pod", podToSelectorItems(pods))
	if err != nil {
		return nil, err
	}

	return &pods[index], nil
}

func podToSelectorItems(pods []sdk.ComponentResourceItem) []string {
	items := []string{}

	for _, resource := range pods {
		items = append(items, fmt.Sprintf(
			"%s / %s",
			resource.GetNamespace(),
			resource.GetName(),
		))
	}

	return items
}
