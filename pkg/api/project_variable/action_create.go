package project_variable

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/config/enum"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type CreateOptions struct {
	common.Options

	sdk.ProjectVariableCreateAction

	Value    string
	IsSecret enum.Bool
}

func NewCreateOptions() *CreateOptions {
	projectVariableCreateOptions := sdk.NewProjectVariableCreateAction("", "", "")

	return &CreateOptions{
		ProjectVariableCreateAction: *projectVariableCreateOptions,

		IsSecret: enum.BoolNone,
	}
}

func (co *CreateOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&co.Name, "name", co.Name, "Unique name for the project variable")
	util.MarkFlagRequiredWithHelp(flags.Lookup("name"), "A unique name within the project for the new project variable")

	flags.StringVar(&co.Value, "value", co.Value, "The value of the project variable")
	util.AppendFlagHelp(flags.Lookup("value"), "A value for this project variable")
	util.MarkFlag(flags.Lookup("value"), util.FlagAllowBlank)

	isSecretFlag := enum.BoolFlag(
		&co.IsSecret,
		"secret",
		"Whether the project variable is secret or not",
	)
	flags.AddFlag(isSecretFlag)
	isSecretFlag.NoOptDefVal = "true"
}

func Create(options *CreateOptions) (*sdk.ProjectVariableItem, error) {
	options.ProjectVariableCreateAction.SetValue(options.Value)

	if options.IsSecret == enum.BoolTrue {
		options.ProjectVariableCreateAction.SetIsSecret(true)
	}

	model, resp, err := CreateRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func CreateRaw(options *CreateOptions) (*sdk.ProjectVariableItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).
		ProjectVariableAPI.ProjectVariableCreate(ctx).
		ProjectVariableCreateAction(options.ProjectVariableCreateAction)

	return request.Execute()
}
