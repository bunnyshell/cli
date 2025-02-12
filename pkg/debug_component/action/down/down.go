package down

import (
	"bunnyshell.com/cli/pkg/k8s/bridge"
	"bunnyshell.com/cli/pkg/debug_component/action"
)

type Options struct {
	ManualSelectSingleResource bool

	resourceLoader *bridge.ResourceLoader

	resourcePath string
}

func NewOptions(
	resourceLoader *bridge.ResourceLoader,
) *Options {
	return &Options{
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
