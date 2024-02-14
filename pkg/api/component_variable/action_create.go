package component_variable

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

	sdk.ServiceComponentVariableCreateAction

	Value    string
	IsSecret enum.Bool
}

func NewCreateOptions() *CreateOptions {
	componentVariableCreateOptions := sdk.NewServiceComponentVariableCreateAction("", "", "")

	return &CreateOptions{
		ServiceComponentVariableCreateAction: *componentVariableCreateOptions,

		IsSecret: enum.BoolNone,
	}
}

func (co *CreateOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&co.Name, "name", co.Name, "Unique name for the component variable")
	util.MarkFlagRequiredWithHelp(flags.Lookup("name"), "A unique name within the component for the new variable")

	flags.StringVar(&co.Value, "value", co.Value, "The value of the component variable")
	util.AppendFlagHelp(flags.Lookup("value"), "A value for this component variable")
	util.MarkFlag(flags.Lookup("value"), util.FlagAllowBlank)

	isSecretFlag := enum.BoolFlag(
		&co.IsSecret,
		"secret",
		"Whether the component variable is secret or not",
	)
	flags.AddFlag(isSecretFlag)
	isSecretFlag.NoOptDefVal = "true"
}

func Create(options *CreateOptions) (*sdk.ServiceComponentVariableItem, error) {
	options.ServiceComponentVariableCreateAction.SetValue(options.Value)

	if options.IsSecret == enum.BoolTrue {
		options.ServiceComponentVariableCreateAction.SetIsSecret(true)
	}

	model, resp, err := CreateRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func CreateRaw(options *CreateOptions) (*sdk.ServiceComponentVariableItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).
		ServiceComponentVariableAPI.ServiceComponentVariableCreate(ctx).
		ServiceComponentVariableCreateAction(options.ServiceComponentVariableCreateAction)

	return request.Execute()
}
