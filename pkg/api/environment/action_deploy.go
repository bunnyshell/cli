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

	return model, err
}

func DeployRaw(options *DeployOptions) (*sdk.EventItem, *http.Response, error) {
	ctx, cancel := lib.GetContext()
	defer cancel()

	request := lib.GetAPI().EnvironmentApi.EnvironmentDeploy(ctx, options.ID)

	return request.Execute()
}