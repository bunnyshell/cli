package up

import (
	"time"

	"bunnyshell.com/cli/pkg/k8s/bridge"
	"bunnyshell.com/cli/pkg/debug_component/action"
)

type Options struct {
	ManualSelectSingleResource bool

	ForceRecreateResource bool

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
}

func NewOptions(
	resourceLoader *bridge.ResourceLoader,
) *Options {
	return &Options{
		resourceLoader: resourceLoader,

		waitTimeout: defaultWaitTimeout,
	}
}

func (up *Options) SetCommand(command []string) {
	up.command = command
}

func (up *Options) Validate() error {
	return nil
}

func (up *Options) ToParameters() (*action.UpParameters, error) {
	up.resourceLoader.ManualSelectSingleResource = up.ManualSelectSingleResource

	parameters := &action.UpParameters{
		ManualSelectSingleResource: up.ManualSelectSingleResource,
		ForceRecreateResource: up.ForceRecreateResource,

		Options: &action.UpOptions{
			WaitTimeout: int64(up.waitTimeout.Seconds()),

			EnvironPairs: up.environPairs,
		},
	}

	up.fillFromFlags(parameters)

	if err := up.loadResource(); err != nil {
		return nil, err
	}

	parameters.Resource = *up.resourceLoader.GetResource()

	if err := parameters.FillFromOptions(); err != nil {
		return nil, err
	}

	return parameters, nil
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

func (up *Options) fillFromFlags(parameters *action.UpParameters) {
	if up.limitCPU != "" {
		parameters.Options.LimitCPU = up.limitCPU
	}

	if up.limitMemory != "" {
		parameters.Options.LimitMemory = up.limitMemory
	}

	if up.requestCPU != "" {
		parameters.Options.RequestCPU = up.requestCPU
	}

	if up.requestMemory != "" {
		parameters.Options.RequestMemory = up.requestMemory
	}

	if up.containerName != "" {
		parameters.Options.ContainerName = up.containerName
	}

	if len(up.command) > 0 {
		parameters.Options.Command = up.command
	}
}
