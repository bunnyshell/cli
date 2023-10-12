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

func NewDeployOptions(id string) *DeployOptions {
	return &DeployOptions{
		PartialActionOptions:  *common.NewPartialActionOptions(id),
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

	isPartialAction := options.IsPartial()

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentDeploy(ctx, options.ID).
		EnvironmentPartialDeployAction(sdk.EnvironmentPartialDeployAction{
			IsPartial:            &isPartialAction,
			Components:           options.GetActionComponents(),
			IncludedDependencies: &options.IncludedDepdendencies,
		})

	return request.Execute()
}
