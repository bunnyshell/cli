package workflow_job

import (
	"net/http"

	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/api/common"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

type LogsOptions struct {
	common.Options

	JobID      string
	StepStatus []string
}

func NewLogsOptions(jobID string) *LogsOptions {
	return &LogsOptions{
		Options: *common.NewOptions(),
		JobID:   jobID,
	}
}

func Logs(options *LogsOptions) (*sdk.WorkflowJobWorkflowJobLogsOutputItem, error) {
	model, resp, err := LogsRaw(options)
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

func LogsRaw(options *LogsOptions) (*sdk.WorkflowJobWorkflowJobLogsOutputItem, *http.Response, error) {
	profile := options.GetProfile()

	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	request := lib.GetAPIFromProfile(profile).WorkflowJobAPI.WorkflowJobLogs(ctx, options.JobID)

	return applyLogsOptions(request, options).Execute()
}

func applyLogsOptions(request sdk.ApiWorkflowJobLogsRequest, options *LogsOptions) sdk.ApiWorkflowJobLogsRequest {
	if options == nil {
		return request
	}

	if len(options.StepStatus) > 0 {
		request = request.StepStatus(options.StepStatus)
	}

	return request
}
