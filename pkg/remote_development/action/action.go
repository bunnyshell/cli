package action

import (
	"errors"
	"fmt"

	"bunnyshell.com/cli/pkg/remote_development/workspace"
	"bunnyshell.com/dev/pkg/remote"
	"bunnyshell.com/sdk"
)

var ErrResourceKindNotSupported = errors.New("resource kind not supported")

type Action struct {
	workspace *workspace.Workspace
}

func NewAction(
	environment sdk.EnvironmentItem,
) *Action {
	return &Action{
		workspace: workspace.NewWorkspace(environment.GetId()),
	}
}

func (action *Action) GetRemoteDev(resource sdk.ComponentResourceItem) (*remote.RemoteDevelopment, error) {
	kubeConfigFile, err := action.workspace.DownloadKubeConfig()
	if err != nil {
		return nil, err
	}

	remoteDev := remote.NewRemoteDevelopment().
		WithKubernetesClient(kubeConfigFile).
		WithNamespaceName(resource.GetNamespace())

	switch kind := resource.GetKind(); kind {
	case "Deployment":
		return remoteDev.WithDeploymentName(resource.GetName()), nil
	case "StatefulSet":
		return remoteDev.WithStatefulSetName(resource.GetName()), nil
	case "DaemonSet":
		return remoteDev.WithDaemonSetName(resource.GetName()), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrResourceKindNotSupported, kind)
	}
}
