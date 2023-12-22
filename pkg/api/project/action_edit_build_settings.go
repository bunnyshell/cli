package project

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/build_settings"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type EditBuildSettingsOptions struct {
	build_settings.EditOptions

	sdk.ProjectEditBuildSettingsAction
}

func NewEditBuildSettingsOptions(project string) *EditBuildSettingsOptions {
	return &EditBuildSettingsOptions{
		EditOptions: *build_settings.NewEditOptions(project),

		ProjectEditBuildSettingsAction: *sdk.NewProjectEditBuildSettingsActionWithDefaults(),
	}
}

func EditBuildSettings(options *EditBuildSettingsOptions) (*sdk.ProjectItem, error) {
	model, resp, err := EditBuildSettingsRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EditBuildSettingsRaw(options *EditBuildSettingsOptions) (*sdk.ProjectItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	build_settings.ApplyEditOptionsToAction(&options.ProjectEditBuildSettingsAction, &options.EditData)

	request := lib.GetAPIFromProfile(profile).ProjectAPI.
		ProjectEditBuildSettings(ctx, options.ID).
		ProjectEditBuildSettingsAction(options.ProjectEditBuildSettingsAction)

	return request.Execute()
}
