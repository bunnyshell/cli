package remote_development

func (r *RemoteDevelopment) Up() error {
	if err := r.ensureEnvironmentWorkspaceDir(); err != nil {
		return err
	}

	if err := r.ensureEnvironmentKubeConfig(); err != nil {
		return err
	}

	r.remoteDev.
		WithKubernetesClient(r.kubeConfigPath).
		WithNamespaceFromKubeConfig().
		WithDeploymentName(r.component.GetName()).
		WithRemoteSyncPath(r.component.GetSyncPath())

	if err := r.remoteDev.SelectContainer(); err != nil {
		return err
	}

	if r.localSyncPath != "" {
		r.remoteDev.WithLocalSyncPath(r.localSyncPath)
	} else if err := r.remoteDev.SelectLocalSyncPath(); err != nil {
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
		WithNamespaceFromKubeConfig().
		WithDeploymentName(r.component.GetName())

	return r.remoteDev.Down()
}

func (r *RemoteDevelopment) StartSSHTerminal() error {
	return r.remoteDev.StartSSHTerminal()
}

func (r *RemoteDevelopment) Wait() error {
	return r.remoteDev.Wait()
}
