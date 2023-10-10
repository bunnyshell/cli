package pipeline

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
	Environment  string
	Event        string

	Status string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		ListOptions: *common.NewListOptions(),
	}
}

func (lo *ListOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&lo.Event, "event", lo.Event, "Filter by EventID")
	flags.StringVar(&lo.Status, "status", lo.Status, "Filter by Status")

	lo.ListOptions.UpdateFlagSet(flags)
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

	request := lib.GetAPIFromProfile(profile).PipelineAPI.PipelineList(ctx)

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

	if options.Status != "" {
		request = request.Status(options.Status)
	}

	return request
}
