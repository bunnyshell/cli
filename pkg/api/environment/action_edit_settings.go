package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type EditSettingsOptions struct {
	common.ItemOptions

	sdk.EnvironmentEditSettings

	EditSettingsData
}

type EditSettingsData struct {
	Name string

	RemoteDevelopmentAllowed bool
	AutoUpdate               bool

	K8SIntegration string
}

func NewEditSettingsOptions(environment string) *EditSettingsOptions {
	environmentEditSettings := sdk.NewEnvironmentEditSettings()

	return &EditSettingsOptions{
		ItemOptions: *common.NewItemOptions(environment),

		EditSettingsData: EditSettingsData{
			RemoteDevelopmentAllowed: true,
			AutoUpdate:               true,
		},

		EnvironmentEditSettings: *environmentEditSettings,
	}
}

func (eso *EditSettingsOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	data := &eso.EditSettingsData

	flags.StringVar(&data.Name, "name", data.Name, "Update environment name")
	flags.BoolVar(&data.RemoteDevelopmentAllowed, "remote-development", data.RemoteDevelopmentAllowed, "Allow remote development for the environment")
	flags.BoolVar(&data.AutoUpdate, "auto-update", data.AutoUpdate, "Automatically update the environment when components git refs change")

	flags.StringVar(&data.K8SIntegration, "k8s", data.K8SIntegration, "Set Kubernetes integration for the environment (if not set)")
}

func EditSettings(options *EditSettingsOptions) (*sdk.EnvironmentItem, error) {
	model, resp, err := EditSettingsRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EditSettingsRaw(options *EditSettingsOptions) (*sdk.EnvironmentItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentApi.EnvironmentEditSettings(ctx, options.ID)

	return applyEditSettingsOptions(request, options).Execute()
}

func applyEditSettingsOptions(
	request sdk.ApiEnvironmentEditSettingsRequest,
	options *EditSettingsOptions,
) sdk.ApiEnvironmentEditSettingsRequest {
	if options.EditSettingsData.Name != "" {
		options.EnvironmentEditSettings.SetName(options.EditSettingsData.Name)
	}

	if options.EditSettingsData.K8SIntegration != "" {
		options.EnvironmentEditSettings.SetKubernetesIntegration(options.EditSettingsData.K8SIntegration)
	}

	options.EnvironmentEditSettings.SetAutoUpdate(options.EditSettingsData.AutoUpdate)
	options.EnvironmentEditSettings.SetRemoteDevelopmentAllowed(options.EditSettingsData.RemoteDevelopmentAllowed)

	request = request.EnvironmentEditSettings(options.EnvironmentEditSettings)

	return request
}
