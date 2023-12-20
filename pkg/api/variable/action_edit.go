package variable

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/config/enum"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type EditOptions struct {
	common.ItemOptions

	sdk.EnvironmentVariableEditAction

	EditData
}

type EditData struct {
	Value    string
	IsSecret enum.Bool
}

func NewEditOptions(id string) *EditOptions {
	return &EditOptions{
		ItemOptions: *common.NewItemOptions(id),

		EditData: EditData{},

		EnvironmentVariableEditAction: *sdk.NewEnvironmentVariableEditActionWithDefaults(),
	}
}

func (eso *EditOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	data := &eso.EditData

	flags.StringVar(&data.Value, "value", data.Value, "Update the environment variable value")

	isSecretFlag := enum.BoolFlag(
		&eso.EditData.IsSecret,
		"secret",
		"Whether the environment variable is secret or not",
	)
	flags.AddFlag(isSecretFlag)
	isSecretFlag.NoOptDefVal = "true"
}

func Edit(options *EditOptions) (*sdk.EnvironmentVariableItem, error) {
	model, resp, err := EditRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func EditRaw(options *EditOptions) (*sdk.EnvironmentVariableItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentVariableAPI.EnvironmentVariableEdit(ctx, options.ID)

	return applyEditOptions(request, options).Execute()
}

func applyEditOptions(
	request sdk.ApiEnvironmentVariableEditRequest,
	options *EditOptions,
) sdk.ApiEnvironmentVariableEditRequest {
	if options.EditData.IsSecret != enum.BoolNone {
		options.EnvironmentVariableEditAction.SetIsSecret(options.EditData.IsSecret == enum.BoolTrue)
	}

	request = request.EnvironmentVariableEditAction(options.EnvironmentVariableEditAction)

	return request
}
