package action

import (
	"errors"
	"fmt"

	"bunnyshell.com/cli/pkg/remote_development/workspace"
	"bunnyshell.com/dev/pkg/debug"
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

func (action *Action) GetDebugCmp(resource sdk.ComponentResourceItem) (*debug.DebugComponent, error) {
	kubeConfigFile, err := action.workspace.DownloadKubeConfig()
	if err != nil {
		return nil, err
	}

	debugCmp := debug.NewDebugComponent().
		WithKubernetesClient(kubeConfigFile).
		WithNamespaceName(resource.GetNamespace())

	switch kind := resource.GetKind(); kind {
	case "Deployment":
		return debugCmp.WithDeploymentName(resource.GetName()), nil
	case "StatefulSet":
		return debugCmp.WithStatefulSetName(resource.GetName()), nil
	case "DaemonSet":
		return debugCmp.WithDaemonSetName(resource.GetName()), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrResourceKindNotSupported, kind)
	}
}
