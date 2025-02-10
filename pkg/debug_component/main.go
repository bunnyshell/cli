package debug_component

import (
	"fmt"
	"os"

	"bunnyshell.com/cli/pkg/interactive"
	"bunnyshell.com/dev/pkg/util"
)

func (d *DebugComponent) Up() error {
	if err := d.ensureEnvironmentWorkspaceDir(); err != nil {
		return err
	}

	if err := d.ensureEnvironmentKubeConfig(); err != nil {
		return err
	}

	componentResource := d.environmentResource.ComponentResource

	d.debugCmp.
		WithKubernetesClient(d.kubeConfigPath).
		WithNamespaceName(componentResource.GetNamespace()).
		WithWaitTimeout(d.waitTimeout)

	switch componentResource.GetKind() {
	case "Deployment":
		d.debugCmp.WithDeploymentName(componentResource.GetName())
	case "StatefulSet":
		d.debugCmp.WithStatefulSetName(componentResource.GetName())
	case "DaemonSet":
		d.debugCmp.WithDaemonSetName(componentResource.GetName())
	default:
		return fmt.Errorf("resource kind \"%s\" is not supported", componentResource.GetKind())
	}

	if err := d.debugCmp.SelectContainer(); err != nil {
		return err
	}

	return d.debugCmp.Up()
}

func (d *DebugComponent) Down() error {
	if err := d.ensureEnvironmentWorkspaceDir(); err != nil {
		return err
	}

	if err := d.ensureEnvironmentKubeConfig(); err != nil {
		return err
	}

	componentResource := d.environmentResource.ComponentResource

	d.debugCmp.
		WithKubernetesClient(d.kubeConfigPath).
		WithNamespaceName(componentResource.GetNamespace())

	switch componentResource.GetKind() {
	case "Deployment":
		d.debugCmp.WithDeploymentName(componentResource.GetName())
	case "StatefulSet":
		d.debugCmp.WithStatefulSetName(componentResource.GetName())
	case "DaemonSet":
		d.debugCmp.WithDaemonSetName(componentResource.GetName())
	default:
		return fmt.Errorf("resource kind \"%s\" is not supported", componentResource.GetKind())
	}

	return d.debugCmp.Down()
}

func (d *DebugComponent) Wait() error {
	return d.debugCmp.Wait()
}
