package remote_development

import (
	"fmt"

	"bunnyshell.com/cli/pkg/environment"
	"bunnyshell.com/dev/pkg/remote"

	remoteDevMutagenConfig "bunnyshell.com/dev/pkg/mutagen/config"
)

var (
	ErrNoOrganizationSelected = fmt.Errorf("you need to select an organization first")
)

type RemoteDevelopment struct {
	remoteDev           *remote.RemoteDevelopment
	environmentResource *environment.EnvironmentResource

	environmentWorkspaceDir string

	kubeConfigPath string

	syncMode       remoteDevMutagenConfig.Mode
	localSyncPath  string
	remoteSyncPath string

	portMappings []string

	waitTimeout int64
}

func NewRemoteDevelopment() *RemoteDevelopment {
	return &RemoteDevelopment{
		remoteDev:           remote.NewRemoteDevelopment(),
		environmentResource: environment.NewEnvironmentResource(),
	}
}

func (r *RemoteDevelopment) WithEnvironmentResource(environmentResource *environment.EnvironmentResource) *RemoteDevelopment {
	r.environmentResource = environmentResource

	return r
}

func (r *RemoteDevelopment) WithEnvironmentWorkspaceDir(environmentWorkspaceDir string) *RemoteDevelopment {
	r.environmentWorkspaceDir = environmentWorkspaceDir

	return r
}

func (r *RemoteDevelopment) WithKubeConfigPath(kubeConfigPath string) *RemoteDevelopment {
	r.kubeConfigPath = kubeConfigPath

	return r
}

func (r *RemoteDevelopment) WithLocalSyncPath(localSyncPath string) *RemoteDevelopment {
	r.localSyncPath = localSyncPath

	return r
}

func (r *RemoteDevelopment) WithRemoteSyncPath(remoteSyncPath string) *RemoteDevelopment {
	r.remoteSyncPath = remoteSyncPath

	return r
}

func (r *RemoteDevelopment) WithPortMappings(portMappings []string) *RemoteDevelopment {
	r.portMappings = portMappings

	return r
}

func (r *RemoteDevelopment) WithWaitTimeout(waitTimeout int64) *RemoteDevelopment {
	r.waitTimeout = waitTimeout

	return r
}

func (r *RemoteDevelopment) WithSyncMode(syncMode remoteDevMutagenConfig.Mode) *RemoteDevelopment {
	r.syncMode = syncMode

	return r
}
