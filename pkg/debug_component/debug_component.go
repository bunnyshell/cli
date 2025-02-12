package debug_component

import (
	"fmt"

	"bunnyshell.com/cli/pkg/environment"
	"bunnyshell.com/dev/pkg/debug"
)

var (
	ErrNoOrganizationSelected = fmt.Errorf("you need to select an organization first")
)

type DebugComponent struct {
	debugCmp            *debug.DebugComponent
	environmentResource *environment.EnvironmentResource

	environmentWorkspaceDir string

	kubeConfigPath string

	waitTimeout int64
}

func NewDebugComponent() *DebugComponent {
	return &DebugComponent{
		debugCmp:            debug.NewDebugComponent(),
		environmentResource: environment.NewEnvironmentResource(),
	}
}

func (d *DebugComponent) WithEnvironmentResource(environmentResource *environment.EnvironmentResource) *DebugComponent {
	d.environmentResource = environmentResource

	return d
}

func (d *DebugComponent) WithEnvironmentWorkspaceDir(environmentWorkspaceDir string) *DebugComponent {
	d.environmentWorkspaceDir = environmentWorkspaceDir

	return d
}

func (d *DebugComponent) WithKubeConfigPath(kubeConfigPath string) *DebugComponent {
	d.kubeConfigPath = kubeConfigPath

	return d
}

func (d *DebugComponent) WithWaitTimeout(waitTimeout int64) *DebugComponent {
	d.waitTimeout = waitTimeout

	return d
}
