package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"bunnyshell.com/cli/pkg/api/workflow_job"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

func init() {
	var pipelineID string
	var outputFormat string

	command := &cobra.Command{
		Use: "jobs",

		Short: "List jobs in a pipeline",
		Long: `List all jobs within a pipeline, showing their ID, name, type, status, and duration.

Use the job IDs from this output with 'bns pipeline logs --job JOB_ID' to view logs for a specific job.

Examples:
  # List jobs in a pipeline
  bns pipeline jobs --id PIPELINE_ID

  # JSON output
  bns pipeline jobs --id PIPELINE_ID --output json`,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			profile := config.GetSettings().Profile

			spinner := util.MakeSpinner(" Fetching pipeline jobs...")
			spinner.Start()

			jobIDs, err := getWorkflowJobs(pipelineID, profile)
			if err != nil {
				spinner.Stop()
				return fmt.Errorf("failed to get jobs for pipeline %s: %w", pipelineID, err)
			}

			var jobs []sdk.WorkflowJobItem
			for _, jobID := range jobIDs {
				info, err := workflow_job.GetJobInfo(profile, jobID)
				if err != nil {
					spinner.Stop()
					return fmt.Errorf("failed to get info for job %s: %w", jobID, err)
				}
				jobs = append(jobs, *info)
			}

			spinner.Stop()

			model := &workflow_job.WorkflowJobList{
				PipelineID: pipelineID,
				Jobs:       jobs,
			}

			return outputJobList(model, outputFormat)
		},
	}

	flags := command.Flags()
	flags.AddFlag(getIDOption(&pipelineID).GetRequiredFlag("id"))
	flags.StringVarP(&outputFormat, "output", "o", "stylish", "Output format: stylish, json")

	config.MainManager.CommandWithGlobalOptions(command)

	mainCmd.AddCommand(command)
}

func outputJobList(data *workflow_job.WorkflowJobList, format string) error {
	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	case "stylish":
		writer := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)

		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", "JobID", "Name", "Type", "Status", "Duration")

		for _, job := range data.Jobs {
			duration := "-"
			if job.HasDuration() {
				duration = (time.Duration(job.GetDuration()) * time.Second).String()
			}

			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n",
				job.GetId(),
				job.GetName(),
				job.GetType(),
				job.GetStatus(),
				duration,
			)
		}

		return writer.Flush()
	default:
		return fmt.Errorf("unknown output format: %s (use: stylish, json)", format)
	}
}
