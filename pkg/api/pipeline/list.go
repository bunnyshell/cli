package pipeline

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type ListOptions struct {
	common.Options

	Page int32

	Organization string
	Environment  string
	Event        string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		Page: 1,
	}
}

func List(options *ListOptions) (*sdk.PaginatedPipelineCollection, error) {
	model, resp, err := ListRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func ListRaw(options *ListOptions) (*sdk.PaginatedPipelineCollection, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).PipelineApi.PipelineList(ctx)

	return applyOptions(request, options).Execute()
}

func applyOptions(request sdk.ApiPipelineListRequest, options *ListOptions) sdk.ApiPipelineListRequest {
	if options == nil {
		return request
	}

	if options.Page > 1 {
		request = request.Page(options.Page)
	}

	if options.Environment != "" {
		request = request.Environment(options.Environment)
	}

	if options.Event != "" {
		request = request.Event(options.Event)
	}

	if options.Organization != "" {
		request = request.Organization(options.Organization)
	}

	return request
}
