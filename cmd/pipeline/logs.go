package pipeline

import (
	"fmt"
	"os"

	"bunnyshell.com/cli/pkg/api/workflow"
	"bunnyshell.com/cli/pkg/api/workflow_job"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/formatter/pipeline_logs"
	"bunnyshell.com/cli/pkg/util"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

func init() {
	options := NewLogsOptions()

	command := &cobra.Command{
		Use:     "logs",
		Aliases: []string{"log"},

		Short: "View pipeline logs for an environment",
		Long: `View and stream logs from pipeline executions (build jobs, deployment steps).

This command fetches logs from all jobs and displays them in a structured format.
Use --follow to stream logs in real-time for active pipelines.

Examples:
  # View latest pipeline logs (all jobs)
  bns pipeline logs --environment ENV_ID

  # View only failed jobs
  bns pipeline logs --environment ENV_ID --failed

  # View a specific pipeline (use 'bns pipeline list' to find IDs)
  bns pipeline logs --id PIPELINE_ID

  # View a specific job within a pipeline
  bns pipeline logs --job JOB_ID

  # Follow active pipeline logs
  bns pipeline logs --environment ENV_ID --follow

  # View only specific step
  bns pipeline logs --environment ENV_ID --step build

  # Show last 50 lines
  bns pipeline logs --environment ENV_ID --tail 50

  # JSON output for parsing
  bns pipeline logs --environment ENV_ID --output json`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Fall back to context if --environment not provided
			if options.EnvironmentID == "" {
				if ctx := config.GetSettings().Profile.Context; ctx.Environment != "" {
					options.EnvironmentID = ctx.Environment
				}
			}

			// Environment ID not required when --job or --id is specified directly
			if options.EnvironmentID == "" && options.JobID == "" && options.WorkflowID == "" {
				return fmt.Errorf("environment required: provide --environment or set context with 'bns configure set-context --environment ID'")
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogs(options)
		},
	}

	// Add flags
	command.Flags().StringVar(&options.EnvironmentID, "environment", "", "Environment ID")
	command.Flags().BoolVarP(&options.Follow, "follow", "f", false, "Follow log output (stream in real-time)")
	command.Flags().BoolVar(&options.Failed, "failed", false, "Show only failed jobs")
	command.Flags().IntVar(&options.Tail, "tail", 0, "Show last N lines")
	command.Flags().StringVar(&options.Step, "step", "", "Filter logs by step name")
	command.Flags().StringVar(&options.JobID, "job", "", "Specific job ID (shows only that job)")
	command.Flags().StringVar(&options.WorkflowID, "id", "", "Pipeline ID (use 'bns pipeline list' to find IDs)")
	command.Flags().StringVarP(&options.OutputFormat, "output", "o", "stylish", "Output format: stylish, json, yaml, raw")

	// Add global options
	config.MainManager.CommandWithGlobalOptions(command)

	mainCmd.AddCommand(command)
}

type LogsOptions struct {
	EnvironmentID string
	WorkflowID    string
	JobID         string
	Follow        bool
	Failed        bool
	Tail          int
	Step          string
	OutputFormat  string

	Profile config.Profile
}

func NewLogsOptions() *LogsOptions {
	return &LogsOptions{
		OutputFormat: "stylish",
	}
}

type workflowInfo struct {
	WorkflowID string
	JobIDs     []string
}

func runLogs(options *LogsOptions) error {
	options.Profile = config.GetSettings().Profile

	// If a specific job ID is given, fetch only that job's logs (legacy single-job mode)
	if options.JobID != "" {
		return runSingleJobLogs(options)
	}

	// Resolve which workflow to use
	wfInfo, err := resolveWorkflow(options)
	if err != nil {
		return err
	}

	// Fetch job metadata for all jobs in the workflow
	spinner := util.MakeSpinner(" Fetching pipeline info...")
	spinner.Start()

	var jobs []*sdk.WorkflowJobItem
	for _, jobID := range wfInfo.JobIDs {
		info, err := workflow_job.GetJobInfo(options.Profile, jobID)
		if err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to get info for job %s: %w", jobID, err)
		}
		jobs = append(jobs, info)
	}
	spinner.Stop()

	// Filter by --failed
	if options.Failed {
		var filtered []*sdk.WorkflowJobItem
		for _, job := range jobs {
			if job.GetStatus() == "failed" {
				filtered = append(filtered, job)
			}
		}
		if len(filtered) == 0 {
			fmt.Fprintln(os.Stderr, "No failed jobs in this pipeline.")
			fmt.Fprintf(os.Stderr, "Jobs: ")
			for i, job := range jobs {
				if i > 0 {
					fmt.Fprintf(os.Stderr, ", ")
				}
				fmt.Fprintf(os.Stderr, "%s [%s]", job.GetName(), job.GetStatus())
			}
			fmt.Fprintln(os.Stderr)
			return nil
		}
		jobs = filtered
	}

	// Fetch logs for each job
	spinner = util.MakeSpinner(fmt.Sprintf(" Fetching logs for %d job(s)...", len(jobs)))
	spinner.Start()

	pipelineLogs := &workflow_job.PipelineLogs{
		WorkflowID: wfInfo.WorkflowID,
	}

	for _, job := range jobs {
		logs, err := fetchJobLogs(options, job.GetId())
		if err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to fetch logs for job %s (%s): %w", job.GetId(), job.GetName(), err)
		}
		logs.JobName = job.GetName()

		// Apply step filter per job
		if options.Step != "" {
			logs = filterByStep(logs, options.Step)
		}

		pipelineLogs.Jobs = append(pipelineLogs.Jobs, *logs)
	}

	spinner.Stop()

	// Apply tail across all jobs
	if options.Tail > 0 {
		pipelineLogs = tailPipelineLogs(pipelineLogs, options.Tail)
	}

	return outputPipelineLogs(pipelineLogs, options.OutputFormat)
}

// runSingleJobLogs handles the --job flag (single job mode)
func runSingleJobLogs(options *LogsOptions) error {
	// Get job info for the name
	info, err := workflow_job.GetJobInfo(options.Profile, options.JobID)
	jobName := options.JobID
	if err == nil {
		jobName = info.GetName()
	}

	logs, err := fetchJobLogs(options, options.JobID)
	if err != nil {
		return err
	}
	logs.JobName = jobName

	if options.Step != "" {
		logs = filterByStep(logs, options.Step)
	}

	pipelineLogs := &workflow_job.PipelineLogs{
		Jobs: []workflow_job.WorkflowJobLogs{*logs},
	}

	if options.Tail > 0 {
		pipelineLogs = tailPipelineLogs(pipelineLogs, options.Tail)
	}

	return outputPipelineLogs(pipelineLogs, options.OutputFormat)
}

// resolveWorkflow determines which workflow to use and returns all its job IDs
func resolveWorkflow(options *LogsOptions) (*workflowInfo, error) {
	if options.WorkflowID != "" {
		// Use explicitly provided workflow ID
		jobIDs, err := getWorkflowJobs(options.WorkflowID, options.Profile)
		if err != nil {
			return nil, fmt.Errorf("failed to get jobs for pipeline %s: %w", options.WorkflowID, err)
		}
		return &workflowInfo{WorkflowID: options.WorkflowID, JobIDs: jobIDs}, nil
	}

	// Auto-detect: find latest workflow for environment
	return getLatestWorkflowForEnvironment(options.EnvironmentID, options.Profile)
}

// getLatestWorkflowForEnvironment finds the latest workflow and returns all its job IDs
func getLatestWorkflowForEnvironment(environmentID string, profile config.Profile) (*workflowInfo, error) {
	collection, err := workflow.List(&workflow.ListOptions{
		Environment: environmentID,
		Profile:     profile,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pipelines: %w", err)
	}

	if collection.Embedded == nil || len(collection.Embedded.Item) == 0 {
		return nil, fmt.Errorf("no pipelines found for environment %s", environmentID)
	}

	workflowID := collection.Embedded.Item[0].GetId()
	jobIDs, err := getWorkflowJobs(workflowID, profile)
	if err != nil {
		return nil, err
	}

	return &workflowInfo{WorkflowID: workflowID, JobIDs: jobIDs}, nil
}

// getWorkflowJobs fetches a workflow by ID and returns its job IDs
func getWorkflowJobs(workflowID string, profile config.Profile) ([]string, error) {
	wf, err := workflow.Get(profile, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pipeline: %w", err)
	}

	jobs := wf.GetJobs()
	if len(jobs) == 0 {
		return nil, fmt.Errorf("pipeline %s has no jobs", workflowID)
	}

	return jobs, nil
}

// fetchJobLogs fetches all pages of logs for a single job
func fetchJobLogs(options *LogsOptions, jobID string) (*workflow_job.WorkflowJobLogs, error) {
	if options.Follow {
		fmt.Fprintln(os.Stderr, "Warning: Follow mode not yet fully implemented, showing current logs...")
	}

	return workflow_job.FetchAllPages(&workflow_job.LogsOptions{
		Profile: options.Profile,
		JobID:   jobID,
		Offset:  0,
		Limit:   1000,
	})
}

// filterByStep filters logs to only show specific step
func filterByStep(logs *workflow_job.WorkflowJobLogs, stepName string) *workflow_job.WorkflowJobLogs {
	filtered := &workflow_job.WorkflowJobLogs{
		WorkflowJobID: logs.WorkflowJobID,
		JobName:       logs.JobName,
		Status:        logs.Status,
		Steps:         []workflow_job.LogStep{},
		Pagination:    logs.Pagination,
	}

	for _, step := range logs.Steps {
		if step.Name == stepName {
			filtered.Steps = append(filtered.Steps, step)
		}
	}

	return filtered
}

// tailPipelineLogs limits output to last N lines across all jobs
func tailPipelineLogs(pl *workflow_job.PipelineLogs, n int) *workflow_job.PipelineLogs {
	// Count total logs across all jobs
	totalLogs := 0
	for _, job := range pl.Jobs {
		for _, step := range job.Steps {
			totalLogs += len(step.Logs)
		}
	}

	if totalLogs <= n {
		return pl
	}

	toSkip := totalLogs - n

	result := &workflow_job.PipelineLogs{
		WorkflowID: pl.WorkflowID,
	}

	skipped := 0
	for _, job := range pl.Jobs {
		newJob := workflow_job.WorkflowJobLogs{
			WorkflowJobID: job.WorkflowJobID,
			JobName:       job.JobName,
			Status:        job.Status,
			Pagination:    job.Pagination,
		}

		for _, step := range job.Steps {
			if skipped+len(step.Logs) <= toSkip {
				skipped += len(step.Logs)
				continue
			}

			newStep := step
			startIdx := toSkip - skipped
			if startIdx < 0 {
				startIdx = 0
			}
			newStep.Logs = step.Logs[startIdx:]
			newJob.Steps = append(newJob.Steps, newStep)
			skipped += len(step.Logs)
		}

		if len(newJob.Steps) > 0 {
			result.Jobs = append(result.Jobs, newJob)
		}
	}

	return result
}

// outputPipelineLogs formats and outputs pipeline logs
func outputPipelineLogs(logs *workflow_job.PipelineLogs, format string) error {
	switch format {
	case "stylish":
		return pipeline_logs.NewStylishFormatter().Format(logs, os.Stdout)
	case "json":
		return pipeline_logs.NewJSONFormatter().Format(logs, os.Stdout)
	case "yaml":
		return pipeline_logs.NewYAMLFormatter().Format(logs, os.Stdout)
	case "raw":
		return pipeline_logs.NewRawFormatter().Format(logs, os.Stdout)
	default:
		return fmt.Errorf("unknown output format: %s (use: stylish, json, yaml, raw)", format)
	}
}
