package workflow_job

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

	Workflow string
	Status   []string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		ListOptions: *common.NewListOptions(),
	}
}

func (lo *ListOptions) UpdateFlagSet(flags *pflag.FlagSet) {
	lo.ListOptions.UpdateFlagSet(flags)
}

func List(options *ListOptions) (*sdk.PaginatedWorkflowJobCollection, error) {
	model, resp, err := ListRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func AllJobs(options *ListOptions) ([]sdk.WorkflowJobCollection, error) {
	var result []sdk.WorkflowJobCollection

	for {
		model, err := List(options)
		if err != nil {
			return nil, err
		}

		if model.Embedded != nil {
			result = append(result, model.Embedded.Item...)
		}

		if !model.HasLinks() || !model.Links.HasNext() {
			return result, nil
		}

		options.Page++
	}
}

func ListRaw(options *ListOptions) (*sdk.PaginatedWorkflowJobCollection, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).WorkflowJobAPI.WorkflowJobList(ctx)

	return applyOptions(request, options).Execute()
}

func applyOptions(request sdk.ApiWorkflowJobListRequest, options *ListOptions) sdk.ApiWorkflowJobListRequest {
	if options == nil {
		return request
	}

	if options.Workflow != "" {
		request = request.Workflow(options.Workflow)
	}

	if options.Page > 1 {
		request = request.Page(options.Page)
	}

	if len(options.Status) > 0 {
		request = request.Status(options.Status)
	}

	return request
}
