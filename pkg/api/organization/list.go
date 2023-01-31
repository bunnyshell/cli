package organization

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
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		Page: 1,
	}
}

func List(options *ListOptions) (*sdk.PaginatedOrganizationCollection, error) {
	model, resp, err := ListRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func ListRaw(options *ListOptions) (*sdk.PaginatedOrganizationCollection, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).OrganizationApi.OrganizationList(ctx)

	return applyOptions(request, options).Execute()
}

func applyOptions(request sdk.ApiOrganizationListRequest, options *ListOptions) sdk.ApiOrganizationListRequest {
	if options == nil {
		return request
	}

	if options.Page > 1 {
		request = request.Page(options.Page)
	}

	return request
}
