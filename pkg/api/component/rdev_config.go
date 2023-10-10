package component

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type RDevConfigOptions struct {
	common.ItemOptions
}

func NewRDevConfigOptions(id string) *RDevContextOptions {
	return &RDevContextOptions{
		ItemOptions: *common.NewItemOptions(id),
	}
}

func RDevConfig(options *RDevContextOptions) (*sdk.ComponentConfigItem, error) {
	model, resp, err := RDevConfigRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func RDevConfigRaw(options *RDevContextOptions) (*sdk.ComponentConfigItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).ComponentAPI.ComponentRemoteDevConfig(ctx, options.ID)

	return request.Execute()
}
