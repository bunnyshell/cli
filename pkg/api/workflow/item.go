package workflow

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

func Get(profile config.Profile, id string) (*sdk.WorkflowItem, error) {
	model, resp, err := GetRaw(profile, id)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}
	return model, nil
}

func GetRaw(profile config.Profile, id string) (*sdk.WorkflowItem, *http.Response, error) {
	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).WorkflowAPI.WorkflowView(ctx, id)

	return request.Execute()
}
