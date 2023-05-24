package action

import (
	mutagenConfig "bunnyshell.com/dev/pkg/mutagen/config"
	"bunnyshell.com/dev/pkg/remote"
	"bunnyshell.com/dev/pkg/remote/container"
	"bunnyshell.com/sdk"
)

type UpOptions struct {
	ContainerName string

	Profile *sdk.ProfileItem

	EnvironPairs []string

	WaitTimeout int64
}

type UpParameters struct {
	Resource sdk.ComponentResourceItem

	SyncMode mutagenConfig.Mode

	ManualSelectSingleResource bool

	LocalSyncPath  string
	RemoteSyncPath string

	PortMappings []string

	Options *UpOptions
}

func (params *UpParameters) FillFromOptions() error {
	if params.Options == nil {
		return nil
	}

	profile := params.Options.Profile

	if profile == nil {
		return nil
	}

	params.PortMappings = append(params.PortMappings, profile.GetPortMapping()...)

	syncPaths := profile.GetSyncPaths()
	if len(syncPaths) > 0 {
		if len(syncPaths) > 1 {
			return ErrOneSyncPathSupported
		}

		syncPath := syncPaths[0]

		if path, ok := syncPath.GetRemotePathOk(); ok {
			params.RemoteSyncPath = *path
		}

		if path, ok := syncPath.GetLocalPathOk(); ok {
			params.LocalSyncPath = *path
		}
	}

	return nil
}

type Up struct {
	Action

	remoteDev *remote.RemoteDevelopment
}

func NewUp(
	environment sdk.EnvironmentItem,
) *Up {
	return &Up{
		Action: *NewAction(environment),
	}
}

func (up *Up) Run(parameters *UpParameters) error {
	remoteDev, err := up.Action.GetRemoteDev(parameters.Resource)
	if err != nil {
		return err
	}

	return up.run(remoteDev, parameters)
}

func (up *Up) StartSSHTerminal() error {
	if up.remoteDev == nil {
		return ErrRemoteDevNotInitialized
	}

	return up.remoteDev.StartSSHTerminal()
}

func (up *Up) Wait() error {
	if up.remoteDev == nil {
		return ErrRemoteDevNotInitialized
	}

	return up.remoteDev.Wait()
}

func (up *Up) run(
	remoteDev *remote.RemoteDevelopment,
	parameters *UpParameters,
) error {
	if err := up.loadRemoteDevOptions(remoteDev, parameters.Options); err != nil {
		return err
	}

	if err := up.ensureSyncPaths(parameters); err != nil {
		return err
	}

	remoteDev.
		WithSyncMode(parameters.SyncMode).
		WithLocalSyncPath(parameters.LocalSyncPath).
		WithRemoteSyncPath(parameters.RemoteSyncPath)

	if err := remoteDev.SelectContainer(); err != nil {
		return err
	}

	if err := remoteDev.PrepareSSHTunnels(parameters.PortMappings); err != nil {
		return err
	}

	up.remoteDev = remoteDev

	return remoteDev.Up()
}

func (up *Up) loadRemoteDevOptions(remoteDev *remote.RemoteDevelopment, options *UpOptions) error {
	if options == nil {
		return nil
	}

	remoteDev.ContainerName = options.ContainerName

	if options.WaitTimeout != 0 {
		remoteDev.WithWaitTimeout(options.WaitTimeout)
	}

	if err := up.loadProfileIntoRemoteDev(remoteDev, options.Profile); err != nil {
		return err
	}

	if len(options.EnvironPairs) > 0 {
		for _, pair := range options.EnvironPairs {
			if err := remoteDev.ContainerConfig.Environ.AddFromDefinition(pair); err != nil {
				return err
			}
		}
	}

	return nil
}

func (up *Up) loadProfileIntoRemoteDev(remoteDev *remote.RemoteDevelopment, profile *sdk.ProfileItem) error {
	if profile == nil {
		return nil
	}

	containerConfig := &remoteDev.ContainerConfig

	containerConfig.Command = profile.Command

	if profile.Environ != nil {
		for key, value := range *profile.Environ {
			containerConfig.Environ.Set(key, value)
		}
	}

	requirements := profile.GetRequirements().ResourceRequirementItem
	if requirements.HasLimits() {
		resourceList := requirements.GetLimits().ResourceListItem

		if err := up.setResourceLimits(containerConfig, resourceList); err != nil {
			return err
		}
	}

	if requirements.HasRequests() {
		resourceList := requirements.GetRequests().ResourceListItem

		if err := up.setResourceRequests(containerConfig, resourceList); err != nil {
			return err
		}
	}

	return nil
}

func (up *Up) setResourceLimits(containerConfig *container.Config, resourceList *sdk.ResourceListItem) error {
	if resourceList.HasCpu() {
		if err := containerConfig.Resources.SetLimitsCPU(resourceList.GetCpu()); err != nil {
			return err
		}
	}

	if resourceList.HasMemory() {
		if err := containerConfig.Resources.SetLimitsMemory(resourceList.GetMemory()); err != nil {
			return err
		}
	}

	return nil
}

func (up *Up) setResourceRequests(containerConfig *container.Config, resourceList *sdk.ResourceListItem) error {
	if resourceList.HasCpu() {
		if err := containerConfig.Resources.SetRequestsCPU(resourceList.GetCpu()); err != nil {
			return err
		}
	}

	if resourceList.HasMemory() {
		if err := containerConfig.Resources.SetRequestsMemory(resourceList.GetMemory()); err != nil {
			return err
		}
	}

	return nil
}
