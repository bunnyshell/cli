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

type CreateOptions struct {
	common.Options

	sdk.ProjectVariableCreateAction

	IsSecret bool
}

func NewCreateOptions() *CreateOptions {
	projectVariableCreateOptions := sdk.NewProjectVariableCreateAction("", "")

	return &CreateOptions{
		ProjectVariableCreateAction: *projectVariableCreateOptions,
	}
}

func (co *CreateOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&co.Name, "name", co.Name, "Unique name for the project variable")
	util.MarkFlagRequiredWithHelp(flags.Lookup("name"), "A unique name within the project for the new project variable")

	flags.BoolVar(&co.IsSecret, "secret", co.IsSecret, "Whether the project variable is secret or not")
}

func Create(options *CreateOptions) (*sdk.ProjectVariableItem, error) {
	if util.IsFlagPassed("secret") {
		options.ProjectVariableCreateAction.IsSecret.Set(&options.IsSecret)
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
