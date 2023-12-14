package project

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type EditBuildSettingsOptions struct {
	common.ItemOptions

	sdk.ProjectEditBuildSettingsAction

	EditBuildSettingsData
}

type EditBuildSettingsData struct {
}

func NewEditBuildSettingsOptions(project string) *EditBuildSettingsOptions {
	return &EditBuildSettingsOptions{
		ItemOptions: *common.NewItemOptions(project),

		EditBuildSettingsData: EditBuildSettingsData{},

		ProjectEditBuildSettingsAction: *sdk.NewProjectEditBuildSettingsAction(),
	}
}

func (eso *EditBuildSettingsOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	// data := &eso.EditSettingsData

	// flags.StringVar(&data.Name, "name", data.Name, "Update project name")

	// flags.StringToStringVar(&data.Labels, "label", data.Labels, "Set labels for the project (key=value)")
	// flags.BoolVar(&data.LabelReplace, "label-replace", data.LabelReplace, "Set label strategy to replace (default: merge)")
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

	request := lib.GetAPIFromProfile(profile).ProjectAPI.ProjectEditBuildSettings(ctx, options.ID)

	return applyEditBuildSettingsOptions(request, options).Execute()
}

func applyEditBuildSettingsOptions(
	request sdk.ApiProjectEditBuildSettingsRequest,
	options *EditBuildSettingsOptions,
) sdk.ApiProjectEditBuildSettingsRequest {
	// if options.EditBuildSettingsData.Name != "" {
	// options.ProjectEditBuildSettingsAction.SetName(options.EditSettingsData.Name)
	// }

	request = request.ProjectEditBuildSettingsAction(options.ProjectEditBuildSettingsAction)

	return request
}
