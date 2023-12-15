package project

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

	sdk.ProjectEditSettingsAction

	EditSettingsData
}

type EditSettingsData struct {
	Name string

	Labels       map[string]string
	LabelReplace bool
}

func NewEditSettingsOptions(project string) *EditSettingsOptions {
	return &EditSettingsOptions{
		ItemOptions: *common.NewItemOptions(project),

		EditSettingsData: EditSettingsData{},

		ProjectEditSettingsAction: *sdk.NewProjectEditSettingsActionWithDefaults(),
	}
}

func (eso *EditSettingsOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	data := &eso.EditSettingsData

	flags.StringVar(&data.Name, "name", data.Name, "Update project name")

	flags.StringToStringVar(&data.Labels, "label", data.Labels, "Set labels for the project (key=value)")
	flags.BoolVar(&data.LabelReplace, "label-replace", data.LabelReplace, "Set label strategy to replace (default: merge)")
}

func EditSettings(options *EditSettingsOptions) (*sdk.ProjectItem, error) {
	model, resp, err := EditSettingsRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EditSettingsRaw(options *EditSettingsOptions) (*sdk.ProjectItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).ProjectAPI.ProjectEditSettings(ctx, options.ID)

	return applyEditSettingsOptions(request, options).Execute()
}

func applyEditSettingsOptions(
	request sdk.ApiProjectEditSettingsRequest,
	options *EditSettingsOptions,
) sdk.ApiProjectEditSettingsRequest {
	if options.EditSettingsData.Name != "" {
		options.ProjectEditSettingsAction.SetName(options.EditSettingsData.Name)
	}

	if options.EditSettingsData.Labels != nil {
		labelsEdit := *sdk.NewEdit()
		if options.EditSettingsData.LabelReplace {
			labelsEdit.SetStrategy("replace")
		}

		labelsEdit.SetValues(options.EditSettingsData.Labels)

		options.ProjectEditSettingsAction.SetLabels(labelsEdit)
	} else if options.EditSettingsData.LabelReplace {
		labelsEdit := *sdk.NewEdit()
		labelsEdit.SetStrategy("replace")

		options.ProjectEditSettingsAction.SetLabels(labelsEdit)
	}

	request = request.ProjectEditSettingsAction(options.ProjectEditSettingsAction)

	return request
}
