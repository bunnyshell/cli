package event

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type ItemOptions struct {
	common.ItemOptions
}

func NewItemOptions(id string) *ItemOptions {
	return &ItemOptions{
		ItemOptions: *common.NewItemOptions(id),
	}
}

func Get(options *ItemOptions) (*sdk.EventItem, error) {
	model, resp, err := GetRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func GetRaw(options *ItemOptions) (*sdk.EventItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EventApi.EventView(ctx, options.ID)

	return request.Execute()
}
