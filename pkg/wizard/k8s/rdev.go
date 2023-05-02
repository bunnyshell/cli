package k8s

import (
	"fmt"

	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/sdk"
)

const (
	Deployment  = "Deployment"
	StatefulSet = "StatefulSet"
	DaemonSet   = "DaemonSet"
)

type RDevListOptions struct {
	Component string
}

func NewRDevListOptions(component string) *RDevListOptions {
	return &RDevListOptions{
		Component: component,
	}
}

func RDevList(options *RDevListOptions) ([]sdk.ComponentResourceItem, error) {
	resources, err := component.Resources(component.NewResourceOptions(options.Component))
	if err != nil {
		return nil, err
	}

	rdevResources := []sdk.ComponentResourceItem{}

	for _, resource := range resources {
		switch resource.GetKind() {
		case Deployment, StatefulSet, DaemonSet:
			rdevResources = append(rdevResources, resource)
		}
	}

	return rdevResources, nil
}

func RDevSelect(options *PodListOptions) (*sdk.ComponentResourceItem, error) {
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

	index, _, err := interactive.Choose("Choose a resource", resourceToSelectorItems(pods))
	if err != nil {
		return nil, err
	}

	return &pods[index], nil
}

func resourceToSelectorItems(pods []sdk.ComponentResourceItem) []string {
	items := []string{}

	for _, resource := range pods {
		items = append(items, fmt.Sprintf(
			"%s / %s / %s",
			resource.GetNamespace(),
			resource.GetKind(),
			resource.GetName(),
		))
	}

	return items
}
