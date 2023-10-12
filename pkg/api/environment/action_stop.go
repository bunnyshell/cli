package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type StopOptions struct {
	common.PartialActionOptions
}

func NewStopOptions(id string) *StopOptions {
	return &StopOptions{
		PartialActionOptions: *common.NewPartialActionOptions(id),
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

	isPartialAction := options.IsPartial()

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentStop(ctx, options.ID).
		EnvironmentPartialAction(sdk.EnvironmentPartialAction{
			IsPartial:  &isPartialAction,
			Components: options.GetActionComponents(),
		})

	return request.Execute()
}
