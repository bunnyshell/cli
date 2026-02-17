package workflow

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/sdk"
)

type ListOptions struct {
	Environment  string
	Organization string
	Status       string

	Profile config.Profile
}

func List(options *ListOptions) (*sdk.PaginatedWorkflowCollection, error) {
	model, resp, err := ListRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}
	return model, nil
}

func ListRaw(options *ListOptions) (*sdk.PaginatedWorkflowCollection, *http.Response, error) {
	ctx, cancel := lib.GetContextFromProfile(options.Profile)
	defer cancel()

	request := lib.GetAPIFromProfile(options.Profile).WorkflowAPI.WorkflowList(ctx)

	if options.Environment != "" {
		request = request.Environment(options.Environment)
	}
	if options.Organization != "" {
		request = request.Organization(options.Organization)
	}
	if options.Status != "" {
		request = request.Status(options.Status)
	}

	return request.Execute()
}
