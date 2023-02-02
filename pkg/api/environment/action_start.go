package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type StartOptions struct {
	common.ActionOptions
}

func NewStartOptions(id string) *StartOptions {
	return &StartOptions{
		ActionOptions: *common.NewActionOptions(id),
	}
}

func Start(options *StartOptions) (*sdk.EventItem, error) {
	model, resp, err := StartRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, err
}

func StartRaw(options *StartOptions) (*sdk.EventItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentApi.EnvironmentStart(ctx, options.ID)

	return request.Execute()
}
