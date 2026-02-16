package pipeline

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"bunnyshell.com/cli/pkg/api/workflow_job"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/formatter/pipeline_logs"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	options := NewLogsOptions()

	command := &cobra.Command{
		Use:     "logs [ENVIRONMENT_ID]",
		Aliases: []string{"log"},

		Short: "View pipeline logs for an environment",
		Long: `View and stream logs from pipeline executions (build jobs, deployment steps).

This command fetches logs from workflow jobs and displays them in a structured format.
Use --follow to stream logs in real-time for active pipelines.

Examples:
  # View latest pipeline logs
  bns pipeline logs my-env

  # Follow active pipeline logs
  bns pipeline logs my-env --follow

  # View only specific step
  bns pipeline logs my-env --step build

  # Show last 50 lines
  bns pipeline logs my-env --tail 50

  # JSON output for parsing
  bns pipeline logs my-env --output json`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Get environment ID from args or context
			if len(args) > 0 {
				options.EnvironmentID = args[0]
			} else if ctx := config.GetSettings().Profile.Context; ctx.Environment != "" {
				options.EnvironmentID = ctx.Environment
			}

			if options.EnvironmentID == "" {
				return fmt.Errorf("environment required: provide ID/name or set context with 'bns configure set-context --environment ID'")
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogs(options)
		},
	}

	// Add flags
	command.Flags().BoolVarP(&options.Follow, "follow", "f", false, "Follow log output (stream in real-time)")
	command.Flags().IntVar(&options.Tail, "tail", 0, "Show last N lines")
	command.Flags().StringVar(&options.Step, "step", "", "Filter logs by step name")
	command.Flags().StringVar(&options.JobID, "job", "", "Specific workflow job ID (defaults to latest)")
	command.Flags().StringVarP(&options.OutputFormat, "output", "o", "stylish", "Output format: stylish, json, yaml, raw")

	// Add global options
	config.MainManager.CommandWithGlobalOptions(command)

	mainCmd.AddCommand(command)
}

type LogsOptions struct {
	EnvironmentID string
	JobID         string
	Follow        bool
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

func runLogs(options *LogsOptions) error {
	options.Profile = config.GetSettings().Profile

	// If no explicit job ID, find the latest workflow job for the environment
	if options.JobID == "" {
		jobID, err := getLatestWorkflowJobForEnvironment(options.EnvironmentID, options.Profile)
		if err != nil {
			return fmt.Errorf("failed to find workflow job: %w", err)
		}
		options.JobID = jobID
	}

	// Fetch logs
	var logs *workflow_job.WorkflowJobLogs
	var err error

	if options.Follow {
		// Follow mode: stream logs with polling
		logs, err = followLogs(options)
	} else {
		// One-shot: fetch all logs
		logs, err = fetchLogs(options)
	}

	if err != nil {
		return err
	}

	// Apply filters
	if options.Step != "" {
		logs = filterByStep(logs, options.Step)
	}

	if options.Tail > 0 {
		logs = tailLogs(logs, options.Tail)
	}

	// Format and output
	return outputLogs(logs, options.OutputFormat)
}

// getLatestWorkflowJobForEnvironment finds the latest workflow job for an environment
func getLatestWorkflowJobForEnvironment(environmentID string, profile config.Profile) (string, error) {
	ctx, cancel := lib.GetContextFromProfile(profile)
	defer cancel()

	// Build API URL to get workflows for environment
	scheme := profile.Scheme
	if scheme == "" {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, profile.Host)
	apiURL := fmt.Sprintf("%s/v1/workflows?environment=%s&page=1", baseURL, environmentID)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", profile.Token)
	req.Header.Set("Accept", "application/hal+json")

	if config.GetSettings().Debug {
		fmt.Fprintf(os.Stderr, "GET %s\n", apiURL)
	}

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch workflows (HTTP %d)", resp.StatusCode)
	}

	// Parse collection response (jobs not included in collection, only in item view)
	var workflowsResp struct {
		Embedded struct {
			Items []struct {
				ID string `json:"id"`
			} `json:"item"`
		} `json:"_embedded"`
	}

	if err := json.Unmarshal(body, &workflowsResp); err != nil {
		return "", fmt.Errorf("failed to parse workflows response: %w", err)
	}

	if len(workflowsResp.Embedded.Items) == 0 {
		return "", fmt.Errorf("no workflows found for environment %s", environmentID)
	}

	// Fetch individual workflow to get its jobs
	workflowID := workflowsResp.Embedded.Items[0].ID
	workflowURL := fmt.Sprintf("%s/v1/workflows/%s", baseURL, workflowID)

	wReq, err := http.NewRequestWithContext(ctx, http.MethodGet, workflowURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create workflow request: %w", err)
	}
	wReq.Header.Set("X-Auth-Token", profile.Token)
	wReq.Header.Set("Accept", "application/hal+json")

	if config.GetSettings().Debug {
		fmt.Fprintf(os.Stderr, "GET %s\n", workflowURL)
	}

	wResp, err := client.Do(wReq)
	if err != nil {
		return "", fmt.Errorf("failed to fetch workflow: %w", err)
	}
	defer wResp.Body.Close()

	wBody, err := io.ReadAll(wResp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read workflow response: %w", err)
	}

	if wResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch workflow %s (HTTP %d)", workflowID, wResp.StatusCode)
	}

	var workflowResp struct {
		ID   string   `json:"id"`
		Jobs []string `json:"jobs"`
	}

	if err := json.Unmarshal(wBody, &workflowResp); err != nil {
		return "", fmt.Errorf("failed to parse workflow response: %w", err)
	}

	if len(workflowResp.Jobs) == 0 {
		return "", fmt.Errorf("workflow %s has no jobs", workflowID)
	}

	// Return last job (most recent)
	return workflowResp.Jobs[len(workflowResp.Jobs)-1], nil
}

// fetchLogs fetches all pages of logs
func fetchLogs(options *LogsOptions) (*workflow_job.WorkflowJobLogs, error) {
	spinner := util.MakeSpinner(" Fetching pipeline logs...")
	spinner.Start()
	defer spinner.Stop()

	logs, err := workflow_job.FetchAllPages(&workflow_job.LogsOptions{
		Profile: options.Profile,
		JobID:   options.JobID,
		Offset:  0,
		Limit:   1000,
	})

	if err != nil {
		return nil, err
	}

	return logs, nil
}

// followLogs streams logs with polling
func followLogs(options *LogsOptions) (*workflow_job.WorkflowJobLogs, error) {
	// TODO: Implement follow mode with polling
	// For now, just fetch once
	fmt.Fprintln(os.Stderr, "⚠ Follow mode not yet fully implemented, showing current logs...")
	return fetchLogs(options)
}

// filterByStep filters logs to only show specific step
func filterByStep(logs *workflow_job.WorkflowJobLogs, stepName string) *workflow_job.WorkflowJobLogs {
	filtered := &workflow_job.WorkflowJobLogs{
		WorkflowJobID: logs.WorkflowJobID,
		Status:        logs.Status,
		Steps:         []workflow_job.LogStep{},
		Pagination:    logs.Pagination,
	}

	for _, step := range logs.Steps {
		if step.Name == stepName {
			filtered.Steps = append(filtered.Steps, step)
			return filtered
		}
	}

	// Step not found
	fmt.Fprintf(os.Stderr, "⚠ Step '%s' not found. Available steps:\n", stepName)
	for _, step := range logs.Steps {
		fmt.Fprintf(os.Stderr, "  - %s\n", step.Name)
	}

	return filtered
}

// tailLogs limits output to last N lines
func tailLogs(logs *workflow_job.WorkflowJobLogs, n int) *workflow_job.WorkflowJobLogs {
	// Count total logs
	totalLogs := 0
	for _, step := range logs.Steps {
		totalLogs += len(step.Logs)
	}

	if totalLogs <= n {
		return logs // No need to tail
	}

	// Calculate how many to skip
	toSkip := totalLogs - n

	tailed := &workflow_job.WorkflowJobLogs{
		WorkflowJobID: logs.WorkflowJobID,
		Status:        logs.Status,
		Steps:         []workflow_job.LogStep{},
		Pagination:    logs.Pagination,
	}

	skipped := 0
	for _, step := range logs.Steps {
		if skipped+len(step.Logs) <= toSkip {
			// Skip entire step
			skipped += len(step.Logs)
			continue
		}

		// Partial step
		newStep := step
		startIdx := toSkip - skipped
		if startIdx < 0 {
			startIdx = 0
		}
		newStep.Logs = step.Logs[startIdx:]
		tailed.Steps = append(tailed.Steps, newStep)

		skipped += len(step.Logs)
	}

	return tailed
}

// outputLogs formats and outputs logs based on format
func outputLogs(logs *workflow_job.WorkflowJobLogs, format string) error {
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
