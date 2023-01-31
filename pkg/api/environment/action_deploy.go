package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

func Deploy(environment string) (*sdk.EventItem, error) {
	model, resp, err := DeployRaw(environment)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, err
}

func DeployRaw(environment string) (*sdk.EventItem, *http.Response, error) {
	ctx, cancel := lib.GetContext()
	defer cancel()

	request := lib.GetAPI().EnvironmentApi.EnvironmentDeploy(ctx, environment)

	return request.Execute()
}
