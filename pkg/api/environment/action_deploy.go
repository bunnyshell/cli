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
	IncludedDependenciesNone    string = "none"
	IncludedDependenciesAll     string = "all"
	IncludedDependenciesMissing string = "missing"
)

type DeployOptions struct {
	common.PartialActionOptions

	IncludedDependencies string

	QueueIfSomethingInProgress bool
}

func NewDeployOptions(id string) *DeployOptions {
	return &DeployOptions{
		PartialActionOptions:       *common.NewPartialActionOptions(id),
		IncludedDependencies:       IncludedDependenciesNone,
		QueueIfSomethingInProgress: false,
	}
}

func (options *DeployOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	options.PartialActionOptions.UpdateFlagSet(flags)

	flags.StringVar(&options.IncludedDependencies, "included-dependencies", options.IncludedDependencies, "Include dependencies in the deployment (none, all, missing)")
	flags.BoolVar(&options.QueueIfSomethingInProgress, "queue", options.QueueIfSomethingInProgress, "Queue the deploy pipeline if another operation is in progress now")
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
			IsPartial:                  &isPartialAction,
			Components:                 options.GetActionComponents(),
			IncludedDependencies:       &options.IncludedDependencies,
			QueueIfSomethingInProgress: &options.QueueIfSomethingInProgress,
		})

	return request.Execute()
}
