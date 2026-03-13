package pipeline

import (
	"fmt"
	"strings"

	"bunnyshell.com/cli/pkg/api/workflow_job"
	wfstatus "bunnyshell.com/cli/pkg/api/workflow_job/status"
	"bunnyshell.com/cli/pkg/formatter"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	var pipelineID string
	var jobs []string
	var jobStatuses []string
	var stepStatuses []string

	command := &cobra.Command{
		Use: "logs",

		Short: "View logs from pipeline jobs",
		Long:  "View logs from pipeline jobs and job steps with optional filtering by job and step status",
		Example: `  # Logs by explicit pipeline ID
  bns pipeline logs --id <PIPELINE_ID>

  # Logs for the latest pipeline in an environment
  bns pipeline logs --id "$(bns pipeline list --environment <ENV_ID> --sort=createdAt:desc -o json | jq -r '._embedded.item[0].id')"

  # Logs for the failed jobs and steps in a pipeline
  bns pipeline logs --id <PIPELINE_ID> --job-status failed --step-status failed
`,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			// If specific jobs are requested, use those directly; otherwise fetch all jobs in the pipeline
			var jobsToFetch []string
			if len(jobs) > 0 {
				jobsToFetch = jobs
			} else {
				listOptions := workflow_job.NewListOptions()
				listOptions.Workflow = pipelineID
				if len(jobStatuses) > 0 {
					listOptions.Status = jobStatuses
				}

				allJobs, err := workflow_job.AllJobs(listOptions)
				if err != nil {
					return fmt.Errorf("failed to list jobs: %w", err)
				}

				for _, job := range allJobs {
					if id, ok := job.GetIdOk(); ok && id != nil {
						jobsToFetch = append(jobsToFetch, *id)
					}
				}
			}

			if len(jobsToFetch) == 0 {
				return lib.FormatCommandData(cmd, []formatter.WorkflowJobLogsResult{})
			}

			// Fetch logs for each job
			var allLogs []formatter.WorkflowJobLogsResult
			for _, jobID := range jobsToFetch {
				logsOptions := workflow_job.NewLogsOptions(jobID)
				if len(stepStatuses) > 0 {
					logsOptions.StepStatus = stepStatuses
				}

				logs, err := workflow_job.Logs(logsOptions)
				if err != nil {
					return fmt.Errorf("failed to fetch logs for job %s: %w", jobID, err)
				}

				allLogs = append(allLogs, formatter.WorkflowJobLogsResult{
					JobID: jobID,
					Logs:  logs,
				})
			}

			// Display the logs
			return lib.FormatCommandData(cmd, allLogs)
		},
	}

	flags := command.Flags()

	flags.AddFlag(getIDOption(&pipelineID).GetRequiredFlag("id"))

	jobStatusValues := strings.Join([]string{
		wfstatus.JobPending, wfstatus.JobQueued, wfstatus.JobInProgress,
		wfstatus.JobFailed, wfstatus.JobAbortFailed, wfstatus.JobSuccess,
		wfstatus.JobSkipped, wfstatus.JobAborting, wfstatus.JobAborted,
	}, ", ")

	stepStatusValues := strings.Join([]string{
		wfstatus.StepFailed, wfstatus.StepSuccess,
	}, ", ")

	flags.StringArrayVar(&jobs, "job", jobs, "Filter to specific job ID(s) (repeatable)")
	flags.StringArrayVar(&jobStatuses, "jobStatus", jobStatuses, fmt.Sprintf("Filter by job status (repeatable); possible values: %s", jobStatusValues))
	flags.StringArrayVar(&stepStatuses, "stepStatus", stepStatuses, fmt.Sprintf("Filter by step status (repeatable); possible values: %s", stepStatusValues))

	mainCmd.AddCommand(command)
}
