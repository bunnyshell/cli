package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/build_settings"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type EditBuildSettingsOptions struct {
	build_settings.EditOptions

	sdk.EnvironmentEditBuildSettingsAction
}

func NewEditBuildSettingsOptions(project string) *EditBuildSettingsOptions {
	return &EditBuildSettingsOptions{
		EditOptions: *build_settings.NewEditOptions(project),

		EnvironmentEditBuildSettingsAction: *sdk.NewEnvironmentEditBuildSettingsAction(),
	}
}

func EditBuildSettings(options *EditBuildSettingsOptions) (*sdk.EnvironmentItem, error) {
	model, resp, err := EditBuildSettingsRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EditBuildSettingsRaw(options *EditBuildSettingsOptions) (*sdk.EnvironmentItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	build_settings.ApplyEditOptionsToAction(&options.EnvironmentEditBuildSettingsAction, &options.EditData)

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.
		EnvironmentEditBuildSettings(ctx, options.ID).
		EnvironmentEditBuildSettingsAction(options.EnvironmentEditBuildSettingsAction)

	return request.Execute()
}
