package up

import (
	"errors"
	"time"

	"bunnyshell.com/cli/pkg/k8s/bridge"
	"bunnyshell.com/cli/pkg/remote_development/action"
	"bunnyshell.com/cli/pkg/remote_development/config"
	"bunnyshell.com/sdk"
)

type Options struct {
	ManualSelectSingleResource bool

	manager *config.Manager

	resourceLoader *bridge.ResourceLoader

	waitTimeout time.Duration

	resourcePath  string
	containerName string

	command []string

	environPairs []string

	limitCPU    string
	limitMemory string

	requestCPU    string
	requestMemory string

	localSyncPath  string
	remoteSyncPath string

	syncMode SyncMode

	portMappings []string
}

func NewOptions(
	manager *config.Manager,
	resourceLoader *bridge.ResourceLoader,
) *Options {
	return &Options{
		manager: manager,

		resourceLoader: resourceLoader,

		waitTimeout: defaultWaitTimeout,

		syncMode: TwoWayResolved,

		portMappings: []string{},
	}
}

func (up *Options) SetCommand(command []string) {
	up.command = command
}

func (up *Options) Validate() error {
	if err := up.manager.Validate(); err != nil {
		return err
	}

	return nil
}

func (up *Options) ToParameters() (*action.UpParameters, error) {
	up.resourceLoader.ManualSelectSingleResource = up.ManualSelectSingleResource

	parameters := &action.UpParameters{
		SyncMode: SyncModeToMutagenMode[up.syncMode],

		ManualSelectSingleResource: up.ManualSelectSingleResource,

		PortMappings: up.portMappings,

		Options: &action.UpOptions{
			WaitTimeout: int64(up.waitTimeout.Seconds()),

			EnvironPairs: up.environPairs,
		},
	}

	if err := up.loadProfile(parameters); err != nil {
		return nil, err
	}

	up.fillFromFlags(parameters)

	if err := up.makeAbsolutePaths(parameters); err != nil {
		return nil, err
	}

	if err := up.loadResource(); err != nil {
		return nil, err
	}

	parameters.Resource = *up.resourceLoader.GetResource()

	if err := parameters.FillFromOptions(); err != nil {
		return nil, err
	}

	return parameters, nil
}

func (up *Options) loadProfile(parameters *action.UpParameters) error {
	profile, err := up.getProfile()
	if err != nil {
		if errors.Is(err, ErrNoRDevConfig) {
			// no profile let intearctive mode handle things
			return nil
		}

		return err
	}

	parameters.Options.Profile = profile

	return nil
}

func (up *Options) loadResource() error {
	if !up.resourceLoader.IsLoaded() {
		return ErrResourceLoaderNotHydrated
	}

	if up.resourceLoader.GetResource() != nil {
		return nil
	}

	if up.resourcePath != "" {
		return up.resourceLoader.SelectResourceFromString(up.resourcePath)
	}

	return up.resourceLoader.SelectResource()
}

func (up *Options) makeAbsolutePaths(parameters *action.UpParameters) error {
	if err := up.manager.MakeAbsolute(&parameters.LocalSyncPath); err != nil {
		return err
	}

	if parameters.Options.Profile == nil {
		return nil
	}

	for _, syncPath := range parameters.Options.Profile.GetSyncPaths() {
		if err := up.manager.MakeAbsolute(syncPath.LocalPath.Get()); err != nil {
			return err
		}
	}

	return nil
}

func (up *Options) fillFromFlags(parameters *action.UpParameters) {
	if up.localSyncPath != "" {
		ensureProfileSyncPath(parameters).SetLocalPath(up.localSyncPath)
	}

	if up.remoteSyncPath != "" {
		ensureProfileSyncPath(parameters).SetRemotePath(up.remoteSyncPath)
	}

	if up.limitCPU != "" {
		ensureLimits(parameters).SetCpu(up.limitCPU)
	}

	if up.limitMemory != "" {
		ensureLimits(parameters).SetMemory(up.limitMemory)
	}

	if up.requestCPU != "" {
		ensureRequests(parameters).SetCpu(up.requestCPU)
	}

	if up.requestMemory != "" {
		ensureRequests(parameters).SetMemory(up.requestMemory)
	}

	if up.containerName != "" {
		parameters.Options.ContainerName = up.containerName
	}

	if len(up.command) > 0 {
		ensureProfile(parameters).Command = up.command
	}
}

func ensureProfile(parameters *action.UpParameters) *sdk.ProfileItem {
	if parameters.Options.Profile == nil {
		parameters.Options.Profile = &sdk.ProfileItem{}
	}

	return parameters.Options.Profile
}

func ensureRequirements(parameters *action.UpParameters) *sdk.ResourceRequirementItem {
	profile := ensureProfile(parameters)

	if !profile.HasRequirements() {
		profile.SetRequirements(sdk.ProfileItemRequirements{
			ResourceRequirementItem: &sdk.ResourceRequirementItem{},
		})
	}

	return profile.GetRequirements().ResourceRequirementItem
}

func ensureLimits(parameters *action.UpParameters) *sdk.ResourceListItem {
	requirements := ensureRequirements(parameters)

	if !requirements.HasLimits() {
		requirements.SetLimits(sdk.ResourceRequirementItemLimits{
			ResourceListItem: &sdk.ResourceListItem{},
		})
	}

	return requirements.GetLimits().ResourceListItem
}

func ensureRequests(parameters *action.UpParameters) *sdk.ResourceListItem {
	requirements := ensureRequirements(parameters)

	if !requirements.HasRequests() {
		requirements.SetRequests(sdk.ResourceRequirementItemRequests{
			ResourceListItem: &sdk.ResourceListItem{},
		})
	}

	return requirements.GetRequests().ResourceListItem
}

func ensureProfileSyncPath(parameters *action.UpParameters) *sdk.SyncPathItem {
	profile := ensureProfile(parameters)

	if !profile.HasSyncPaths() {
		profile.SyncPaths = []sdk.SyncPathItem{
			*sdk.NewSyncPathItem(),
		}
	}

	return &profile.SyncPaths[0]
}
