package remote_development

import "fmt"

func (r *RemoteDevelopment) Up() error {
	if err := r.ensureEnvironmentWorkspaceDir(); err != nil {
		return err
	}

	if err := r.ensureEnvironmentKubeConfig(); err != nil {
		return err
	}

	r.remoteDev.
		WithKubernetesClient(r.kubeConfigPath).
		WithNamespaceName(r.componentResource.GetNamespace())

	switch r.componentResource.GetKind() {
	case "Deployment":
		r.remoteDev.WithDeploymentName(r.componentResource.GetName())
	case "StatefulSet":
		r.remoteDev.WithStatefulSetName(r.componentResource.GetName())
	case "DaemonSet":
		r.remoteDev.WithDaemonSetName(r.componentResource.GetName())
	default:
		return fmt.Errorf("resource kind \"%s\" is not supported", r.componentResource.GetKind())
	}

	resource, err := r.remoteDev.GetResource()
	if err != nil {
		return err
	}
	if r.remoteDev.IsActiveForResource(resource) {
		return fmt.Errorf("the selected resource is already under remote development")
	}

	if err := r.remoteDev.SelectContainer(); err != nil {
		return err
	}

	if r.localSyncPath != "" {
		r.remoteDev.WithLocalSyncPath(r.localSyncPath)
	} else if err := r.remoteDev.SelectLocalSyncPath(); err != nil {
		return err
	}

	if r.remoteSyncPath != "" {
		r.remoteDev.WithRemoteSyncPath(r.remoteSyncPath)
	} else if err := r.remoteDev.SelectRemoteSyncPath(); err != nil {
		return err
	}

	return r.remoteDev.Up()
}

func (r *RemoteDevelopment) Down() error {
	if err := r.ensureEnvironmentWorkspaceDir(); err != nil {
		return err
	}

	if err := r.ensureEnvironmentKubeConfig(); err != nil {
		return err
	}

	r.remoteDev.
		WithKubernetesClient(r.kubeConfigPath).
		WithNamespaceName(r.componentResource.GetNamespace())

	switch r.componentResource.GetKind() {
	case "Deployment":
		r.remoteDev.WithDeploymentName(r.componentResource.GetName())
	case "StatefulSet":
		r.remoteDev.WithStatefulSetName(r.componentResource.GetName())
	case "DaemonSet":
		r.remoteDev.WithDaemonSetName(r.componentResource.GetName())
	default:
		return fmt.Errorf("resource kind \"%s\" is not supported", r.componentResource.GetKind())
	}

	return r.remoteDev.Down()
}

func (r *RemoteDevelopment) StartSSHTerminal() error {
	return r.remoteDev.StartSSHTerminal()
}

func (r *RemoteDevelopment) Wait() error {
	return r.remoteDev.Wait()
}
