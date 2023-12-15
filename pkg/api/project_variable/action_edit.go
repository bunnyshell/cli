package project_variable

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type EditOptions struct {
	common.ItemOptions

	sdk.ProjectVariableEditAction

	EditData
}

type EditData struct {
	Value    string
	IsSecret bool
}

func NewEditOptions(id string) *EditOptions {
	return &EditOptions{
		ItemOptions: *common.NewItemOptions(id),

		EditData: EditData{},

		ProjectVariableEditAction: *sdk.NewProjectVariableEditAction(),
	}
}

func (eso *EditOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	data := &eso.EditData

	flags.StringVar(&data.Value, "value", data.Value, "Update the project variable value")
	flags.BoolVar(&data.IsSecret, "secret", data.IsSecret, "Whether the project variable is secret or not")
}

func Edit(options *EditOptions) (*sdk.ProjectVariableItem, error) {
	model, resp, err := EditRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EditRaw(options *EditOptions) (*sdk.ProjectVariableItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).ProjectVariableAPI.ProjectVariableEdit(ctx, options.ID)

	return applyEditOptions(request, options).Execute()
}

func applyEditOptions(
	request sdk.ApiProjectVariableEditRequest,
	options *EditOptions,
) sdk.ApiProjectVariableEditRequest {
	if util.IsFlagPassed("value") {
		options.ProjectVariableEditAction.SetValue(options.EditData.Value)
	}

	if util.IsFlagPassed("secret") {
		options.ProjectVariableEditAction.SetIsSecret(options.EditData.IsSecret)
	}

	request = request.ProjectVariableEditAction(options.ProjectVariableEditAction)

	return request
}
