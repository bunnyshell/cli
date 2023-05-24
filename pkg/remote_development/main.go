package remote_development

import (
	"fmt"
	"os"

	"bunnyshell.com/cli/pkg/interactive"
	mutagenConfig "bunnyshell.com/dev/pkg/mutagen/config"
	"bunnyshell.com/dev/pkg/util"
)

func (r *RemoteDevelopment) Up() error {
	if err := r.ensureEnvironmentWorkspaceDir(); err != nil {
		return err
	}

	if err := r.ensureEnvironmentKubeConfig(); err != nil {
		return err
	}

	componentResource := r.environmentResource.ComponentResource

	r.remoteDev.
		WithKubernetesClient(r.kubeConfigPath).
		WithNamespaceName(componentResource.GetNamespace()).
		WithWaitTimeout(r.waitTimeout).
		WithSyncMode(r.syncMode)

	switch componentResource.GetKind() {
	case "Deployment":
		r.remoteDev.WithDeploymentName(componentResource.GetName())
	case "StatefulSet":
		r.remoteDev.WithStatefulSetName(componentResource.GetName())
	case "DaemonSet":
		r.remoteDev.WithDaemonSetName(componentResource.GetName())
	default:
		return fmt.Errorf("resource kind \"%s\" is not supported", componentResource.GetKind())
	}

	if err := r.remoteDev.SelectContainer(); err != nil {
		return err
	}

	if err := r.ensureSyncPaths(); err != nil {
		return err
	}

	r.remoteDev.WithLocalSyncPath(r.localSyncPath)
	r.remoteDev.WithRemoteSyncPath(r.remoteSyncPath)

	err := r.remoteDev.PrepareSSHTunnels(r.portMappings)
	if err != nil {
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

	componentResource := r.environmentResource.ComponentResource

	r.remoteDev.
		WithKubernetesClient(r.kubeConfigPath).
		WithNamespaceName(componentResource.GetNamespace())

	switch componentResource.GetKind() {
	case "Deployment":
		r.remoteDev.WithDeploymentName(componentResource.GetName())
	case "StatefulSet":
		r.remoteDev.WithStatefulSetName(componentResource.GetName())
	case "DaemonSet":
		r.remoteDev.WithDaemonSetName(componentResource.GetName())
	default:
		return fmt.Errorf("resource kind \"%s\" is not supported", componentResource.GetKind())
	}

	return r.remoteDev.Down()
}

func (r *RemoteDevelopment) StartSSHTerminal() error {
	return r.remoteDev.StartSSHTerminal()
}

func (r *RemoteDevelopment) Wait() error {
	return r.remoteDev.Wait()
}

func (r *RemoteDevelopment) ensureSyncPaths() error {
	if r.syncMode == mutagenConfig.None {
		return r.ensurePersistentWorkdir()
	}

	if err := r.ensureLocalSyncPath(); err != nil {
		return err
	}

	if err := r.ensureRemoteSyncPath(); err != nil {
		return err
	}

	return nil
}

func (r *RemoteDevelopment) ensureLocalSyncPath() error {
	if r.localSyncPath != "" {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	question := interactive.NewInput("Local Sync Path")
	question.Default = cwd
	question.Help = "Local path is the folder on your machine that will be synced into the container"
	question.SetValidate(util.IsDirectoryValidator)

	syncPath, err := question.AskString()
	if err != nil {
		return err
	}

	r.WithLocalSyncPath(syncPath)

	return nil
}

func (r *RemoteDevelopment) ensureRemoteSyncPath() error {
	if r.remoteSyncPath != "" {
		return nil
	}

	question := interactive.NewInput("Remote Sync Path")
	question.Help = "Remote path is the folder within the container where the application is loaded from\n" +
		"This is where the local files will be synced to\n" +
		"This folder will be mounted as a persistent volume to persist your changes across multiple development sessions."

	syncPath, err := question.AskString()
	if err != nil {
		return err
	}

	r.WithRemoteSyncPath(syncPath)

	return nil
}

func (r *RemoteDevelopment) ensurePersistentWorkdir() error {
	if r.remoteSyncPath != "" {
		return nil
	}

	question := interactive.NewInput("Persistent Workdir")
	question.Help = "Persistent workdir is the folder within the container where the application is loaded from\n" +
		"This folder will be mounted as a persistent volume to persist your changes across multiple development sessions."

	syncPath, err := question.AskString()
	if err != nil {
		return err
	}

	r.WithRemoteSyncPath(syncPath)

	return nil
}
