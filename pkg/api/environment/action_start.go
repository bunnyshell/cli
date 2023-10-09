package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type StartOptions struct {
	common.PartialActionOptions

	WithDependencies bool
}

func NewStartOptions(id string, components []string) *StartOptions {
	return &StartOptions{
		PartialActionOptions: *common.NewPartialActionOptions(id, len(components) > 0, components),
	}
}

func (options *StartOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	options.PartialActionOptions.UpdateFlagSet(flags)

	flags.BoolVar(&options.WithDependencies, "with-dependencies", options.WithDependencies, "Start the component dependencies too.")
}

func Start(options *StartOptions) (*sdk.EventItem, error) {
	model, resp, err := StartRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func StartRaw(options *StartOptions) (*sdk.EventItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentApi.EnvironmentStart(ctx, options.ID).
		EnvironmentPartialStartAction(sdk.EnvironmentPartialStartAction{
			IsPartial:        &options.IsPartial,
			Components:       options.Components,
			WithDependencies: &options.WithDependencies,
		})

	return request.Execute()
}
