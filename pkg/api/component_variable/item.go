package component_variable

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

func NewItemOptions(id string) *common.ItemOptions {
	return common.NewItemOptions(id)
}

func Get(options *common.ItemOptions) (*sdk.ServiceComponentVariableItem, error) {
	model, resp, err := GetRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func GetRaw(options *common.ItemOptions) (*sdk.ServiceComponentVariableItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).ServiceComponentVariableAPI.ServiceComponentVariableView(ctx, options.ID)

	return request.Execute()
}
