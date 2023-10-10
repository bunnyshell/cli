package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type EditComponentOptions struct {
	DeployOptions

	EditComponentsData

	WithDeploy bool
}

type EditComponentsData struct {
	K8SIntegration string

	Component string

	SourceRepository string
	SourceBranch     string

	TargetRepository string
	TargetBranch     string
}

func NewEditComponentOptions() *EditComponentOptions {
	return &EditComponentOptions{
		DeployOptions: *NewDeployOptions(""),
	}
}

func (eo *EditComponentOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	data := eo.EditComponentsData

	flags.BoolVar(&eo.WithDeploy, "deploy", eo.WithDeploy, "Deploy the environment after update")

	flags.StringVar(&data.K8SIntegration, "k8s", data.K8SIntegration, "Set Kubernetes integration for the environment (if not set)")

	eo.DeployOptions.UpdateFlagSet(flags)
}

func EditComponents(options *EditComponentOptions) (*sdk.EnvironmentItem, error) {
	model, resp, err := EditComponentsRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EditComponentsRaw(options *EditComponentOptions) (*sdk.EnvironmentItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentEditComponents(ctx, options.ID)

	return appyEditComponentsOptions(request, options).Execute()
}

func appyEditComponentsOptions(
	request sdk.ApiEnvironmentEditComponentsRequest,
	options *EditComponentOptions,
) sdk.ApiEnvironmentEditComponentsRequest {
	action := sdk.NewEnvironmentEditComponentsAction(
		getFilter(&options.EditComponentsData),
		getGitInfo(&options.EditComponentsData),
	)

	request = request.EnvironmentEditComponentsAction(*action)

	return request
}

func getFilter(data *EditComponentsData) sdk.EnvironmentEditComponentsActionFilter {
	if data.Component != "" {
		filter := sdk.NewFilterName()
		filter.SetName(data.Component)

		return sdk.FilterNameAsEnvironmentEditComponentsActionFilter(filter)
	}

	filter := sdk.NewFilterGit()

	if data.SourceRepository != "" {
		filter.SetRepository(data.SourceRepository)
	}

	if data.SourceBranch != "" {
		filter.SetBranch(data.SourceBranch)
	}

	return sdk.FilterGitAsEnvironmentEditComponentsActionFilter(filter)
}

func getGitInfo(data *EditComponentsData) sdk.GitInfo {
	gitInfo := sdk.NewGitInfo()

	if data.TargetRepository != "" {
		gitInfo.SetRepository(data.TargetRepository)
	}

	if data.TargetBranch != "" {
		gitInfo.SetBranch(data.TargetBranch)
	}

	return *gitInfo
}
