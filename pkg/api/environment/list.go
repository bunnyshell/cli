package environment

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

	Organization          string
	Project               string
	KubernetesIntegration string

	Type            string
	ClusterStatus   string
	OperationStatus string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		Page: 1,
	}
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
