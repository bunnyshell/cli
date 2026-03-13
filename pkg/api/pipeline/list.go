package pipeline

import (
	"fmt"
	"net/http"
	"strings"

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
	Sort   []string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		ListOptions: *common.NewListOptions(),
	}
}

func (lo *ListOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&lo.Event, "event", lo.Event, "Filter by EventID")
	flags.StringVar(&lo.Status, "status", lo.Status, "Filter by Status")
	flags.StringArrayVar(&lo.Sort, "sort", lo.Sort, "Sort by field and direction (repeatable); format: createdAt:asc|desc")

	lo.ListOptions.UpdateFlagSet(flags)
}

func List(options *ListOptions) (*sdk.PaginatedWorkflowCollection, error) {
	model, resp, err := ListRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func ListRaw(options *ListOptions) (*sdk.PaginatedWorkflowCollection, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).WorkflowAPI.WorkflowList(ctx)

	request, err := applyOptions(request, options)
	if err != nil {
		return nil, nil, err
	}

	return request.Execute()
}

func applyOptions(request sdk.ApiWorkflowListRequest, options *ListOptions) (sdk.ApiWorkflowListRequest, error) {
	if options == nil {
		return request, nil
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

	for _, sortValue := range options.Sort {
		field, direction, found := strings.Cut(sortValue, ":")
		if !found || field == "" || direction == "" {
			return request, fmt.Errorf(`invalid sort value %q, expected format "createdAt:asc|desc"`, sortValue)
		}

		if field != "createdAt" {
			return request, fmt.Errorf(`unsupported sort field %q, supported fields: createdAt`, field)
		}

		switch strings.ToLower(direction) {
		case "asc", "desc":
			request = request.OrderCreatedAt(strings.ToLower(direction))
		default:
			return request, fmt.Errorf(`unsupported sort direction %q, supported directions: asc, desc`, direction)
		}
	}

	return request, nil
}
