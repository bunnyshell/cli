package variable_group

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

	sdk.EnvironItemCreateAction

	Value    string
	IsSecret enum.Bool
}

func NewCreateOptions() *CreateOptions {
	return &CreateOptions{
		EnvironItemCreateAction: *sdk.NewEnvironItemCreateActionWithDefaults(),
	}
}

func (co *CreateOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&co.GroupName, "group", co.GroupName, "Environment variable group name")
	util.MarkFlagRequiredWithHelp(flags.Lookup("group"), "The group in which the environment variable should be created")

	flags.StringVar(&co.Name, "name", co.Name, "Unique name for the environment variable")
	util.MarkFlagRequiredWithHelp(flags.Lookup("name"), "A unique name within the environment for the new environment variable")

	flags.StringVar(&co.Value, "value", co.Value, "The value of the environment variable")
	util.AppendFlagHelp(flags.Lookup("value"), "A value for this environment variable")
	util.MarkFlag(flags.Lookup("value"), util.FlagAllowBlank)

	isSecretFlag := enum.BoolFlag(
		&co.IsSecret,
		"secret",
		"Whether the environment variable is secret or not",
	)
	flags.AddFlag(isSecretFlag)
	isSecretFlag.NoOptDefVal = "true"
}

func Create(options *CreateOptions) (*sdk.EnvironItemItem, error) {
	options.EnvironItemCreateAction.SetValue(options.Value)

	if options.IsSecret == enum.BoolTrue {
		options.EnvironItemCreateAction.SetIsSecret(true)
	}

	model, resp, err := CreateRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func CreateRaw(options *CreateOptions) (*sdk.EnvironItemItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).
		EnvironItemAPI.EnvironItemCreate(ctx).
		EnvironItemCreateAction(options.EnvironItemCreateAction)

	return request.Execute()
}
