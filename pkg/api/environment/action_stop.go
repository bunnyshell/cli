package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type StopOptions struct {
	common.ActionOptions
}

func NewStopOptions(id string) *StopOptions {
	return &StopOptions{
		ActionOptions: *common.NewActionOptions(id),
	}
}

func Stop(options *StopOptions) (*sdk.EventItem, error) {
	model, resp, err := StopRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func StopRaw(options *StopOptions) (*sdk.EventItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentApi.EnvironmentStop(ctx, options.ID)

	return request.Execute()
}
