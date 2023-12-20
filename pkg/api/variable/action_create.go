package variable

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

	sdk.EnvironmentVariableCreateAction

	Value    string
	IsSecret enum.Bool
}

func NewCreateOptions() *CreateOptions {
	variableCreateActionCreateOptions := sdk.NewEnvironmentVariableCreateActionWithDefaults()

	return &CreateOptions{
		EnvironmentVariableCreateAction: *variableCreateActionCreateOptions,
	}
}

func (co *CreateOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&co.Name, "name", co.Name, "Unique name for the environment variable")
	util.MarkFlagRequiredWithHelp(flags.Lookup("name"), "A unique name within the environment for the new environment variable")

	flags.StringVar(&co.Value, "value", co.Value, "The value of the project variable")

	isSecretFlag := enum.BoolFlag(
		&co.IsSecret,
		"secret",
		"Whether the project variable is secret or not",
	)
	flags.AddFlag(isSecretFlag)
	isSecretFlag.NoOptDefVal = "true"
}

func Create(options *CreateOptions) (*sdk.EnvironmentVariableItem, error) {
	options.EnvironmentVariableCreateAction.SetValue(options.Value)

	if options.IsSecret == enum.BoolTrue {
		options.EnvironmentVariableCreateAction.SetIsSecret(true)
	}

	model, resp, err := CreateRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func CreateRaw(options *CreateOptions) (*sdk.EnvironmentVariableItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).
		EnvironmentVariableAPI.EnvironmentVariableCreate(ctx).
		EnvironmentVariableCreateAction(options.EnvironmentVariableCreateAction)

	return request.Execute()
}
