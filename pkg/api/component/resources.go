package component

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type ResourceOptions struct {
	common.ItemOptions
}

func NewResourceOptions(id string) *ResourceOptions {
	return &ResourceOptions{
		ItemOptions: *common.NewItemOptions(id),
	}
}

func Resources(options *ResourceOptions) ([]sdk.ComponentResourceItem, error) {
	model, resp, err := ResourcesRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func ResourcesRaw(options *ResourceOptions) ([]sdk.ComponentResourceItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).ComponentAPI.ComponentResources(ctx, options.ID)

	return request.Execute()
}
