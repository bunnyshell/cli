package k8s

import (
	"context"
	"fmt"

	"bunnyshell.com/cli/pkg/interactive"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type ContainerListOptions struct {
	Namespace string
	PodName   string
	Container string

	Client v1.PodsGetter
}

type ContainerItem struct {
	fmt.Stringer

	*corev1.Container

	State corev1.ContainerState
}

func (c *ContainerItem) String() string {
	status := getStatus(c.State)

	return fmt.Sprintf("[%s] %s", status, c.Container.Name)
}

func ContainerList(options *ContainerListOptions) ([]*ContainerItem, error) {
	pod, err := options.Client.Pods(options.Namespace).Get(context.TODO(), options.PodName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	containers := []*ContainerItem{}

	for index := range pod.Spec.Containers {
		containers = append(containers, &ContainerItem{
			Container: &pod.Spec.Containers[index],

			State: pod.Status.ContainerStatuses[index].State,
		})
	}

	return containers, nil
}

func ContainerSelect(options *ContainerListOptions) (*ContainerItem, error) {
	containers, err := ContainerList(options)
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, errEmptyList
	}

	if len(containers) == 1 {
		return containers[0], nil
	}

	index, _, err := interactive.Choose("Choose a container", containersToSelectorItems(containers))
	if err != nil {
		return nil, err
	}

	return containers[index], nil
}

func containersToSelectorItems(containers []*ContainerItem) []string {
	names := make([]string, len(containers))

	for index, container := range containers {
		names[index] = container.String()
	}

	return names
}

func getStatus(state corev1.ContainerState) string {
	switch {
	case state.Waiting != nil:
		return state.Waiting.Reason
	case state.Running != nil:
		return "running"
	case state.Terminated != nil:
		return state.Terminated.Reason
	}

	return "unknown"
}
