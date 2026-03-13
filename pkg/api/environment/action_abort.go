package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type AbortOptions struct {
	common.ItemOptions
}

func NewAbortOptions(id string) *AbortOptions {
	return &AbortOptions{
		ItemOptions: *common.NewItemOptions(id),
	}
}

func Abort(options *AbortOptions) (*sdk.EventItem, error) {
	model, resp, err := AbortRaw(options)
	if resp != nil && resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func AbortRaw(options *AbortOptions) (*sdk.EventItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentAbort(ctx, options.ID)

	return request.Execute()
}
