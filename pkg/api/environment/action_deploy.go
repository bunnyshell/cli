package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

const (
	IncludedDepdendenciesNone    string = "none"
	IncludedDepdendenciesAll     string = "all"
	IncludedDepdendenciesMissing string = "missing"
)

type DeployOptions struct {
	common.PartialActionOptions

	IncludedDepdendencies string
}

func NewDeployOptions(id string, isPartial bool, components []string) *DeployOptions {
	return &DeployOptions{
		PartialActionOptions:  *common.NewPartialActionOptions(id, isPartial, components),
		IncludedDepdendencies: IncludedDepdendenciesNone,
	}
}

func (options *DeployOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	options.PartialActionOptions.UpdateFlagSet(flags)

	flags.StringVar(&options.IncludedDepdendencies, "included-dependencies", options.IncludedDepdendencies, "Include dependencies in the deployment (none, all, missing)")
}

func Deploy(options *DeployOptions) (*sdk.EventItem, error) {
	model, resp, err := DeployRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func DeployRaw(options *DeployOptions) (*sdk.EventItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentApi.EnvironmentDeploy(ctx, options.ID).
		EnvironmentPartialDeployAction(sdk.EnvironmentPartialDeployAction{
			IsPartial:            &options.IsPartial,
			Components:           options.Components,
			IncludedDependencies: &options.IncludedDepdendencies,
		})

	return request.Execute()
}
