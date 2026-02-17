package workflow_job

import (
	"bunnyshell.com/cli/pkg/api"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

// PipelineLogs wraps logs from all jobs in a workflow
type PipelineLogs struct {
	WorkflowID string            `json:"workflowId"`
	Jobs       []WorkflowJobLogs `json:"jobs"`
}

// WorkflowJobList wraps a list of jobs for a pipeline
type WorkflowJobList struct {
	PipelineID string             `json:"pipelineId"`
	Jobs       []sdk.WorkflowJobItem `json:"jobs"`
}

// WorkflowJobLogs represents the structure of workflow job logs
type WorkflowJobLogs struct {
	WorkflowJobID string     `json:"workflowJobId"`
	JobName       string     `json:"jobName,omitempty"`
	Status        string     `json:"status"`
	Steps         []LogStep  `json:"steps"`
	Pagination    Pagination `json:"pagination"`
}

// LogStep represents a single step in the workflow
type LogStep struct {
	Name      string       `json:"name"`
	Status    string       `json:"status"`
	ExitCode  int          `json:"exitCode"`
	StartedAt string       `json:"startedAt"`
	Logs      []LogMessage `json:"logs"`
}

// LogMessage represents a single log message
type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Offset  int  `json:"offset"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	HasMore bool `json:"hasMore"`
}

// LogsOptions contains options for fetching workflow job logs
type LogsOptions struct {
	Profile config.Profile
	JobID   string
	Offset  int
	Limit   int
}

// GetJobInfo fetches metadata for a workflow job via the SDK
func GetJobInfo(profile config.Profile, jobID string) (*sdk.WorkflowJobItem, error) {
	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	model, resp, err := lib.GetAPIFromProfile(profile).WorkflowJobAPI.WorkflowJobView(ctx, jobID).Execute()
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return model, nil
}

// GetLogs fetches workflow job logs from the API via the SDK
func GetLogs(options *LogsOptions) (*WorkflowJobLogs, error) {
	ctx, cancel := lib.GetContextFromProfile(options.Profile)
	defer cancel()

	request := lib.GetAPIFromProfile(options.Profile).WorkflowJobAPI.WorkflowJobLogsList(ctx, options.JobID)
	request = request.Offset(int32(options.Offset))
	request = request.Limit(int32(options.Limit))

	sdkLogs, resp, err := request.Execute()
	if err != nil {
		return nil, api.ParseError(resp, err)
	}

	return fromSDKLogs(sdkLogs), nil
}

// FetchAllPages fetches all pages of logs automatically
func FetchAllPages(options *LogsOptions) (*WorkflowJobLogs, error) {
	var allLogs *WorkflowJobLogs
	var allSteps []LogStep

	offset := options.Offset
	limit := options.Limit

	for {
		logs, err := GetLogs(&LogsOptions{
			Profile: options.Profile,
			JobID:   options.JobID,
			Offset:  offset,
			Limit:   limit,
		})
		if err != nil {
			return nil, err
		}

		if allLogs == nil {
			allLogs = logs
			allSteps = logs.Steps
		} else {
			allSteps = mergeSteps(allSteps, logs.Steps)
		}

		if !logs.Pagination.HasMore {
			break
		}

		offset += limit
	}

	if allLogs != nil {
		allLogs.Steps = allSteps
		allLogs.Pagination.HasMore = false
	}

	return allLogs, nil
}

// fromSDKLogs converts SDK response to CLI types
func fromSDKLogs(sdkLogs *sdk.WorkflowJobLogsResponse) *WorkflowJobLogs {
	logs := &WorkflowJobLogs{
		WorkflowJobID: sdkLogs.WorkflowJobID,
		Status:        sdkLogs.Status,
		Pagination: Pagination{
			Offset:  sdkLogs.Pagination.Offset,
			Limit:   sdkLogs.Pagination.Limit,
			Total:   sdkLogs.Pagination.Total,
			HasMore: sdkLogs.Pagination.HasMore,
		},
	}

	for _, sdkStep := range sdkLogs.Steps {
		step := LogStep{
			Name:      sdkStep.Name,
			Status:    sdkStep.Status,
			ExitCode:  sdkStep.ExitCode,
			StartedAt: sdkStep.StartedAt,
		}
		for _, sdkLog := range sdkStep.Logs {
			step.Logs = append(step.Logs, LogMessage{
				Timestamp: sdkLog.Timestamp,
				Message:   sdkLog.Message,
			})
		}
		logs.Steps = append(logs.Steps, step)
	}

	return logs
}

// mergeSteps merges log steps, combining logs from the same step
func mergeSteps(existing, new []LogStep) []LogStep {
	if len(existing) > 0 && len(new) > 0 {
		lastExisting := &existing[len(existing)-1]
		firstNew := new[0]

		if lastExisting.Name == firstNew.Name {
			lastExisting.Logs = append(lastExisting.Logs, firstNew.Logs...)
			return append(existing, new[1:]...)
		}
	}

	return append(existing, new...)
}
