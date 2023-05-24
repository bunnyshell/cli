package environment

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

	Organization          string
	Project               string
	KubernetesIntegration string

	Type            string
	ClusterStatus   string
	OperationStatus string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		ListOptions: *common.NewListOptions(),
	}
}

func (lo *ListOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	flags.StringVar(&lo.Type, "type", lo.Type, "Filter by Type")
	flags.StringVar(&lo.ClusterStatus, "clusterStatus", lo.ClusterStatus, "Filter by Cluster Status")
	flags.StringVar(&lo.OperationStatus, "operationStatus", lo.OperationStatus, "Filter by Operation Status")
	flags.StringVar(&lo.KubernetesIntegration, "k8sCluster", lo.KubernetesIntegration, "Filter by K8SIntegrationID")

	lo.ListOptions.UpdateFlagSet(flags)
}

func List(options *ListOptions) (*sdk.PaginatedEnvironmentCollection, error) {
	model, resp, err := ListRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func ListRaw(options *ListOptions) (*sdk.PaginatedEnvironmentCollection, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).EnvironmentApi.EnvironmentList(ctx)

	return applyOptions(request, options).Execute()
}

func applyOptions(request sdk.ApiEnvironmentListRequest, options *ListOptions) sdk.ApiEnvironmentListRequest {
	if options == nil {
		return request
	}

	if options.Page > 1 {
		request = request.Page(options.Page)
	}

	if options.Organization != "" {
		request = request.Organization(options.Organization)
	}

	if options.Project != "" {
		request = request.Project(options.Project)
	}

	if options.KubernetesIntegration != "" {
		request = request.KubernetesIntegration(options.KubernetesIntegration)
	}

	if options.ClusterStatus != "" {
		request = request.ClusterStatus(options.ClusterStatus)
	}

	if options.OperationStatus != "" {
		request = request.OperationStatus(options.OperationStatus)
	}

	if options.Type != "" {
		request = request.Type_(options.Type)
	}

	return request
}
