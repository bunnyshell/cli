package action

import (
	"bunnyshell.com/dev/pkg/debug"
	"bunnyshell.com/dev/pkg/remote/container"
	"bunnyshell.com/sdk"
)

type UpOptions struct {
	ContainerName string

	EnvironPairs []string

	Command []string

    LimitCPU    string
    LimitMemory string

    RequestCPU    string
    RequestMemory string

	WaitTimeout int64
}

type UpParameters struct {
	Resource sdk.ComponentResourceItem

	ManualSelectSingleResource bool

	ForceRecreateResource bool

	Options *UpOptions
}

func (params *UpParameters) FillFromOptions() error {
	if params.Options == nil {
		return nil
	}

	return nil
}

type Up struct {
	Action

	debugCmp *debug.DebugComponent
}

func NewUp(
	environment sdk.EnvironmentItem,
) *Up {
	return &Up{
		Action: *NewAction(environment),
	}
}

func (up *Up) Run(parameters *UpParameters) error {
	debugCmp, err := up.Action.GetDebugCmp(parameters.Resource)
	if err != nil {
		return err
	}

	return up.run(debugCmp, parameters)
}

func (up *Up) Wait() error {
	if up.debugCmp == nil {
		return ErrDebugCmpNotInitialized
	}

	return up.debugCmp.Wait()
}

func (up *Up) Close() error {
	if up.debugCmp == nil {
		return ErrDebugCmpNotInitialized
	}

	up.debugCmp.Close()

	return nil
}

func (up *Up) run(
	debugCmp *debug.DebugComponent,
	parameters *UpParameters,
) error {
	if err := up.loadDebugCmpOptions(debugCmp, parameters.Options); err != nil {
		return err
	}

	if err := debugCmp.SelectContainer(); err != nil {
		return err
	}

    if err := debugCmp.CanUp(parameters.ForceRecreateResource); err != nil {
        return err
    }

	up.debugCmp = debugCmp

	return debugCmp.Up()
}

func (up *Up) loadDebugCmpOptions(debugCmp *debug.DebugComponent, options *UpOptions) error {
	if options == nil {
		return nil
	}

	debugCmp.ContainerName = options.ContainerName

	if options.WaitTimeout != 0 {
		debugCmp.WithWaitTimeout(options.WaitTimeout)
	}

    up.setContainerResources(&debugCmp.ContainerConfig, options)

	if len(options.EnvironPairs) > 0 {
		for _, pair := range options.EnvironPairs {
			if err := debugCmp.ContainerConfig.Environ.AddFromDefinition(pair); err != nil {
				return err
			}
		}
	}

	return nil
}

func (up *Up) setContainerResources(containerConfig *container.Config, options *UpOptions) error {
	if options.RequestCPU != "" {
		if err := containerConfig.Resources.SetRequestsCPU(options.RequestCPU); err != nil {
			return err
		}
	}

	if options.RequestMemory != "" {
		if err := containerConfig.Resources.SetRequestsMemory(options.RequestMemory); err != nil {
			return err
		}
	}

    if options.LimitCPU != "" {
        if err := containerConfig.Resources.SetLimitsCPU(options.LimitCPU); err != nil {
            return err
        }
    }

    if options.LimitMemory != "" {
        if err := containerConfig.Resources.SetLimitsMemory(options.LimitMemory); err != nil {
            return err
        }
    }

	return nil
}

func (up *Up) GetSelectedContainerName() (string, error) {
	if up.debugCmp == nil {
		return "", ErrDebugCmpNotInitialized
	}

	return up.debugCmp.GetSelectedContainerName()
}