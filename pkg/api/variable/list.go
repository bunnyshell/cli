package variable

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

	Name string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		ListOptions: *common.NewListOptions(),
	}
}

func (lo *ListOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&lo.Name, "name", lo.Name, "Filter by Name")

	lo.ListOptions.UpdateFlagSet(flags)
}

func List(options *ListOptions) (*sdk.PaginatedEnvironmentVariableCollection, error) {
	model, resp, err := ListRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func ListRaw(options *ListOptions) (*sdk.PaginatedEnvironmentVariableCollection, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentVariableApi.EnvironmentVariableList(ctx)

	return applyOptions(request, options).Execute()
}

func applyOptions(request sdk.ApiEnvironmentVariableListRequest, options *ListOptions) sdk.ApiEnvironmentVariableListRequest {
	if options == nil {
		return request
	}

	if options.Page > 1 {
		request = request.Page(options.Page)
	}

	if options.Organization != "" {
		request = request.Organization(options.Organization)
	}

	if options.Environment != "" {
		request = request.Environment(options.Environment)
	}

	if options.Name != "" {
		request = request.Name(options.Name)
	}

	return request
}
