package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type CloneOptions struct {
	common.ActionOptions

	sdk.EnvironmentCloneAction
}

func NewCloneOptions(id string) *CloneOptions {
	return &CloneOptions{
		ActionOptions: *common.NewActionOptions(id),

		EnvironmentCloneAction: *sdk.NewEnvironmentCloneActionWithDefaults(),
	}
}

func (co *CloneOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&co.EnvironmentCloneAction.Name, "name", co.EnvironmentCloneAction.Name, "Environment Clone Name")
}

func Clone(options *CloneOptions) (*sdk.EnvironmentItem, error) {
	model, resp, err := CloneRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, err
}

func CloneRaw(options *CloneOptions) (*sdk.EnvironmentItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentApi.EnvironmentClone(ctx, options.ID)

	request = request.EnvironmentCloneAction(options.EnvironmentCloneAction)

	return request.Execute()
}
