package project

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

	sdk.ProjectCreateAction
}

func NewCreateOptions() *CreateOptions {
	projectCreateOptions := sdk.NewProjectCreateAction("", "")

	labels := make(map[string]string)
	projectCreateOptions.Labels = &labels

	return &CreateOptions{
		ProjectCreateAction: *projectCreateOptions,
	}
}

func (co *CreateOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&co.Name, "name", co.Name, "Unique name for the project")
	util.MarkFlagRequiredWithHelp(flags.Lookup("name"), "A unique name within the organization for the new project")

	flags.StringToStringVar(co.Labels, "label", *co.Labels, "Set labels for the new project (key=value)")
}

func Create(options *CreateOptions) (*sdk.ProjectItem, error) {
	model, resp, err := CreateRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func CreateRaw(options *CreateOptions) (*sdk.ProjectItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).
		ProjectAPI.ProjectCreate(ctx).
		ProjectCreateAction(options.ProjectCreateAction)

	return request.Execute()
}
