package component

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type RDevContextOptions struct {
	common.ItemOptions

	Profile any
}

func NewRDevContextOptions(id string) *RDevContextOptions {
	return &RDevContextOptions{
		ItemOptions: *common.NewItemOptions(id),
	}
}

func RDevContext(options *RDevContextOptions) (*sdk.ComponentProfileItem, error) {
	model, resp, err := RDevContextRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func RDevContextRaw(options *RDevContextOptions) (*sdk.ComponentProfileItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).ComponentAPI.ComponentRemoteDevProfile(ctx, options.ID)

	return request.Body(options.Profile).Execute()
}
