package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/config/enum"
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

	RemoteDevelopmentAllowed enum.Bool
	AutoUpdate               enum.Bool

	K8SIntegration string
}

func NewEditSettingsOptions(environment string) *EditSettingsOptions {
	return &EditSettingsOptions{
		ItemOptions: *common.NewItemOptions(environment),

		EditSettingsData: EditSettingsData{
			RemoteDevelopmentAllowed: enum.BoolNone,
			AutoUpdate:               enum.BoolNone,
		},

		EnvironmentEditSettings: *sdk.NewEnvironmentEditSettings(),
	}
}

func (eso *EditSettingsOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	data := &eso.EditSettingsData

	flags.StringVar(&data.Name, "name", data.Name, "Update environment name")

	autoUpdateFlag := enum.BoolFlag(
		&data.AutoUpdate,
		"auto-update",
		"Automatically update the environment when components git refs change",
	)
	flags.AddFlag(autoUpdateFlag)
	autoUpdateFlag.NoOptDefVal = "true"

	rdevFlag := enum.BoolFlag(
		&data.RemoteDevelopmentAllowed,
		"remote-development",
		"Allow remote development for the environment",
	)
	flags.AddFlag(rdevFlag)
	rdevFlag.NoOptDefVal = "true"
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

	if options.EditSettingsData.RemoteDevelopmentAllowed != enum.BoolNone {
		options.EnvironmentEditSettings.SetRemoteDevelopmentAllowed(options.EditSettingsData.RemoteDevelopmentAllowed == enum.BoolTrue)
	}

	if options.EditSettingsData.AutoUpdate != enum.BoolNone {
		options.EnvironmentEditSettings.SetAutoUpdate(options.EditSettingsData.AutoUpdate == enum.BoolTrue)
	}

	request = request.EnvironmentEditSettings(options.EnvironmentEditSettings)

	return request
}
