package project

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/pflag"
)

type ListOptions struct {
	common.ListOptions

	Organization string

	Search string

	Labels map[string]string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		ListOptions: *common.NewListOptions(),
	}
}

func (lo *ListOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&lo.Search, "search", lo.Search, "Search by name")

	flags.StringToStringVar(&lo.Labels, "label", lo.Labels, "Filter by label (key=value)")

	lo.ListOptions.UpdateFlagSet(flags)
}

func List(options *ListOptions) (*sdk.PaginatedProjectCollection, error) {
	model, resp, err := ListRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func ListRaw(options *ListOptions) (*sdk.PaginatedProjectCollection, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).ProjectAPI.ProjectList(ctx)

	return applyOptions(request, options).Execute()
}

func applyOptions(request sdk.ApiProjectListRequest, options *ListOptions) sdk.ApiProjectListRequest {
	if options == nil {
		return request
	}

	if options.Page > 1 {
		request = request.Page(options.Page)
	}

	if options.Search != "" {
		request = request.Search(options.Search)
	}

	if options.Organization != "" {
		request = request.Organization(options.Organization)
	}

	if len(options.Labels) > 0 {
		request = request.Labels(options.Labels)
	}

	return request
}
