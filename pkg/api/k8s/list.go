package k8s

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

	CloudProvider string
	Status        string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		ListOptions: *common.NewListOptions(),
	}
}

func (lo *ListOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&lo.CloudProvider, "cloudProvider", lo.CloudProvider, "Filter by Cloud Provider")
	flags.StringVar(&lo.Status, "status", lo.Status, "Filter by Status")

	lo.ListOptions.UpdateFlagSet(flags)
}

func List(options *ListOptions) (*sdk.PaginatedKubernetesIntegrationCollection, error) {
	model, resp, err := ListRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func ListRaw(options *ListOptions) (*sdk.PaginatedKubernetesIntegrationCollection, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).KubernetesIntegrationAPI.KubernetesIntegrationList(ctx)

	return applyOptions(request, options).Execute()
}

func applyOptions(request sdk.ApiKubernetesIntegrationListRequest, options *ListOptions) sdk.ApiKubernetesIntegrationListRequest {
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

	if options.CloudProvider != "" {
		request = request.CloudProvider(options.CloudProvider)
	}

	if options.Status != "" {
		request = request.Status(options.Status)
	}

	return request
}
