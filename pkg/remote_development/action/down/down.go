package down

import (
	"bunnyshell.com/cli/pkg/k8s/bridge"
	"bunnyshell.com/cli/pkg/remote_development/action"
	"bunnyshell.com/cli/pkg/remote_development/config"
)

type Options struct {
	ManualSelectSingleResource bool

	manager *config.Manager

	resourceLoader *bridge.ResourceLoader

	resourcePath string
}

func NewOptions(
	manager *config.Manager,
	resourceLoader *bridge.ResourceLoader,
) *Options {
	return &Options{
		manager: manager,

		resourceLoader: resourceLoader,
	}
}

func (down *Options) ToParameters() (*action.DownParameters, error) {
	down.resourceLoader.ManualSelectSingleResource = down.ManualSelectSingleResource

	if err := down.loadResource(); err != nil {
		return nil, err
	}

	parameters := &action.DownParameters{
		Resource: *down.resourceLoader.GetResource(),
	}

	return parameters, nil
}

func (down *Options) loadResource() error {
	if !down.resourceLoader.IsLoaded() {
		return ErrResourceLoaderNotHydrated
	}

	if down.resourceLoader.GetResource() != nil {
		return nil
	}

	if down.resourcePath != "" {
		return down.resourceLoader.SelectResourceFromString(down.resourcePath)
	}

	return down.resourceLoader.SelectResource()
}
