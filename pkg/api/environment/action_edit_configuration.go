package environment

import (
	"errors"
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

type EditConfigurationOptions struct {
	DeployOptions

	sdk.EnvironmentEditConfiguration

	EditConfigurationData

	WithDeploy bool

	genesisSourceOptions GenesisSourceOptions
}

type EditConfigurationData struct {
	K8SIntegration string
}

func NewEditConfigurationOptions(environment string) *EditConfigurationOptions {
	environmentEditConfiguration := sdk.NewEnvironmentEditConfiguration()

	return &EditConfigurationOptions{
		DeployOptions: *NewDeployOptions(environment),

		EnvironmentEditConfiguration: *environmentEditConfiguration,

		genesisSourceOptions: *NewGenesisSourceOptions(),
	}
}

func (eco *EditConfigurationOptions) UpdateCommandFlags(command *cobra.Command) {
	flags := command.Flags()

	data := &eco.EditConfigurationData

	flags.StringVar(&data.K8SIntegration, "k8s", data.K8SIntegration, "Set Kubernetes integration for the environment (if not set)")

	flags.BoolVar(&eco.WithDeploy, "deploy", eco.WithDeploy, "Deploy the environment after update")

	eco.DeployOptions.UpdateFlagSet(flags)

	eco.genesisSourceOptions.updateCommandFlags(command, "update")
}

func (eco *EditConfigurationOptions) Validate() error {
	return eco.genesisSourceOptions.validate()
}

func (eco *EditConfigurationOptions) AttachGenesis() error {
	fromGit, fromGitSpec, fromString, fromTemplate, err := eco.genesisSourceOptions.getGenesis()
	if err != nil {
		return err
	}

	eco.Configuration = &sdk.EnvironmentEditConfigurationConfiguration{
		FromGit:      fromGit,
		FromGitSpec:  fromGitSpec,
		FromTemplate: fromTemplate,
		FromString:   fromString,
	}

	return nil
}

func (eco *EditConfigurationOptions) HandleError(cmd *cobra.Command, err error) error {
	var apiError api.Error

	if errors.As(err, &apiError) {
		return eco.genesisSourceOptions.handleError(cmd, apiError)
	}

	return lib.FormatCommandError(cmd, err)
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

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentEditConfiguration(ctx, options.ID)

	return applyEditConfigurationoptions(request, options).Execute()
}

func applyEditConfigurationoptions(
	request sdk.ApiEnvironmentEditConfigurationRequest,
	options *EditConfigurationOptions,
) sdk.ApiEnvironmentEditConfigurationRequest {
	return request.EnvironmentEditConfiguration(options.EnvironmentEditConfiguration)
}
