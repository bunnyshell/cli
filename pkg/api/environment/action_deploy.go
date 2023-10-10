package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type DeployOptions struct {
	common.ActionOptions
}

func NewDeployOptions(id string) *DeployOptions {
	return &DeployOptions{
		ActionOptions: *common.NewActionOptions(id),
	}
}

func Deploy(options *DeployOptions) (*sdk.EventItem, error) {
	model, resp, err := DeployRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func DeployRaw(options *DeployOptions) (*sdk.EventItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentAPI.EnvironmentDeploy(ctx, options.ID)

	return request.Execute()
}
