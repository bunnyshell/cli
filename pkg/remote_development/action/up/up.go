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

func (up *Options) ToParameters() (*action.UpParameters, error) {
	up.resourceLoader.ManualSelectSingleResource = up.ManualSelectSingleResource

	parameters := &action.UpParameters{
		SyncMode: SyncModeToMutagenMode[up.syncMode],

		PortMappings: up.portMappings,

		Options: &action.UpOptions{
			WaitTimeout: int64(up.waitTimeout.Seconds()),

			EnvironPairs: up.environPairs,
		},
	}

	if err := up.loadProfile(parameters); err != nil {
		return nil, err
	}

	if up.localSyncPath != "" {
		ensureProfileSyncPath(parameters).SetLocalPath(up.localSyncPath)
	}

	if up.remoteSyncPath != "" {
		ensureProfileSyncPath(parameters).SetRemotePath(up.remoteSyncPath)
	}

	if up.containerName != "" {
		parameters.Options.ContainerName = up.containerName
	}

	if len(up.command) > 0 {
		if parameters.Options.Profile == nil {
			parameters.Options.Profile = &sdk.ProfileItem{}
		}

		parameters.Options.Profile.Command = up.command
	}

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

	for _, syncPath := range parameters.Options.Profile.SyncPaths {
		if err := up.manager.MakeAbsolute(syncPath.LocalPath.Get()); err != nil {
			return err
		}
	}

	return nil
}

func ensureProfileSyncPath(parameters *action.UpParameters) *sdk.SyncPathItem {
	if parameters.Options.Profile == nil {
		parameters.Options.Profile = &sdk.ProfileItem{}
	}

	if parameters.Options.Profile.SyncPaths == nil {
		parameters.Options.Profile.SyncPaths = []sdk.SyncPathItem{
			*sdk.NewSyncPathItem(),
		}
	}

	return &parameters.Options.Profile.SyncPaths[0]
}
