package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type DeleteOptions struct {
	common.ActionOptions
}

func NewDeleteOptions(id string) *DeleteOptions {
	return &DeleteOptions{
		ActionOptions: *common.NewActionOptions(id),
	}
}

func Delete(options *DeleteOptions) (*sdk.EventItem, error) {
	model, resp, err := DeleteRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, err
}

func DeleteRaw(options *DeleteOptions) (*sdk.EventItem, *http.Response, error) {
	ctx, cancel := lib.GetContext()
	defer cancel()

	request := lib.GetAPI().EnvironmentApi.EnvironmentDelete(ctx, options.ID)

	return request.Execute()
}
