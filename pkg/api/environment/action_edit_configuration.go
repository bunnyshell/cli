package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type EditConfigurationOptions struct {
	DeployOptions

	sdk.EnvironmentEditConfiguration

	EditConfigurationData

	WithDeploy bool
}

type EditConfigurationData struct {
	K8SIntegration string
}

func NewEditConfigurationOptions(environment string) *EditConfigurationOptions {
	environmentEditConfiguration := sdk.NewEnvironmentEditConfiguration()

	return &EditConfigurationOptions{
		DeployOptions: *NewDeployOptions(environment, false, []string{}),

		EnvironmentEditConfiguration: *environmentEditConfiguration,
	}
}

func (eco *EditConfigurationOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	data := &eco.EditConfigurationData

	flags.StringVar(&data.K8SIntegration, "k8s", data.K8SIntegration, "Set Kubernetes integration for the environment (if not set)")

	flags.BoolVar(&eco.WithDeploy, "deploy", eco.WithDeploy, "Deploy the environment after update")

	eco.DeployOptions.UpdateFlagSet(flags)
}

func EditConfiguration(options *EditConfigurationOptions) (*sdk.EnvironmentItem, error) {
	model, resp, err := EditConfigurationRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EditConfigurationRaw(options *EditConfigurationOptions) (*sdk.EnvironmentItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentApi.EnvironmentEditConfiguration(ctx, options.ID)

	return applyEditConfigurationoptions(request, options).Execute()
}

func applyEditConfigurationoptions(
	request sdk.ApiEnvironmentEditConfigurationRequest,
	options *EditConfigurationOptions,
) sdk.ApiEnvironmentEditConfigurationRequest {
	return request.EnvironmentEditConfiguration(options.EnvironmentEditConfiguration)
}
