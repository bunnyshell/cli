package pipeline_logs

import (
	"fmt"
	"io"
	"strings"
	"time"

	"bunnyshell.com/cli/pkg/api/workflow_job"
	"github.com/fatih/color"
)

// StylishFormatter formats logs with colors and visual hierarchy
type StylishFormatter struct {
	colorEnabled bool
}

// NewStylishFormatter creates a new stylish formatter
func NewStylishFormatter() *StylishFormatter {
	return &StylishFormatter{
		colorEnabled: true,
	}
}

// Format outputs pipeline logs in stylish format
func (f *StylishFormatter) Format(logs *workflow_job.PipelineLogs, w io.Writer) error {
	if logs.WorkflowID != "" {
		fmt.Fprintf(w, "\nWorkflow: %s\n", logs.WorkflowID)
	}

	for i, job := range logs.Jobs {
		if i > 0 {
			fmt.Fprintln(w)
		}
		f.printJobHeader(w, &job)

		for _, step := range job.Steps {
			f.printStep(w, &step)
		}
	}

	f.printSummary(w, logs)

	return nil
}

// printJobHeader prints a header for each workflow job
func (f *StylishFormatter) printJobHeader(w io.Writer, job *workflow_job.WorkflowJobLogs) {
	separator := strings.Repeat("═", 70)
	fmt.Fprintf(w, "\n%s\n", color.New(color.FgCyan, color.Bold).Sprint(separator))

	statusIcon := f.getStatusIcon(job.Status)
	jobLabel := job.WorkflowJobID
	if job.JobName != "" {
		jobLabel = job.JobName
	}

	header := fmt.Sprintf("%s Job: %s", statusIcon, jobLabel)

	switch job.Status {
	case "success":
		fmt.Fprintln(w, color.New(color.FgCyan, color.Bold).Sprint(header))
	case "failed":
		fmt.Fprintln(w, color.New(color.FgRed, color.Bold).Sprint(header))
	case "running":
		fmt.Fprintln(w, color.New(color.FgYellow, color.Bold).Sprint(header))
	default:
		fmt.Fprintln(w, color.New(color.Bold).Sprint(header))
	}

	fmt.Fprintf(w, "%s  %s\n", color.New(color.Faint).Sprintf("  ID: %s", job.WorkflowJobID),
		color.New(color.Faint).Sprintf("Status: %s", f.colorizeStatus(job.Status)))
	fmt.Fprintf(w, "%s\n", color.New(color.FgCyan, color.Bold).Sprint(separator))
}

// printStep prints a single step with its logs
func (f *StylishFormatter) printStep(w io.Writer, step *workflow_job.LogStep) {
	separator := strings.Repeat("━", 70)
	fmt.Fprintf(w, "%s\n", color.New(color.Faint).Sprint(separator))

	statusIcon := f.getStatusIcon(step.Status)
	stepHeader := fmt.Sprintf("%s Step: %s", statusIcon, step.Name)

	if step.Status == "success" {
		fmt.Fprintln(w, color.GreenString(stepHeader))
	} else if step.Status == "failed" {
		fmt.Fprintln(w, color.RedString(stepHeader))
	} else if step.Status == "running" {
		fmt.Fprintln(w, color.YellowString(stepHeader))
	} else {
		fmt.Fprintln(w, stepHeader)
	}

	fmt.Fprintf(w, "%s\n\n", color.New(color.Faint).Sprint(separator))

	for _, log := range step.Logs {
		f.printLogMessage(w, &log)
	}

	fmt.Fprintln(w)
}

// printLogMessage prints a single log message
func (f *StylishFormatter) printLogMessage(w io.Writer, log *workflow_job.LogMessage) {
	timestamp := f.formatTimestamp(log.Timestamp)
	timestampStr := color.New(color.Faint).Sprintf("  %s", timestamp)

	message := fmt.Sprintf("  %s", log.Message)

	fmt.Fprintf(w, "%s  %s\n", timestampStr, message)
}

// printSummary prints summary information
func (f *StylishFormatter) printSummary(w io.Writer, logs *workflow_job.PipelineLogs) {
	separator := strings.Repeat("━", 70)
	fmt.Fprintf(w, "%s\n\n", color.New(color.Faint).Sprint(separator))

	totalLogs := 0
	totalJobs := len(logs.Jobs)
	failedJobs := 0
	for _, job := range logs.Jobs {
		if job.Status == "failed" {
			failedJobs++
		}
		for _, step := range job.Steps {
			totalLogs += len(step.Logs)
		}
	}

	fmt.Fprintf(w, "Jobs: %d", totalJobs)
	if failedJobs > 0 {
		fmt.Fprintf(w, " (%s)", color.RedString("%d failed", failedJobs))
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Total log lines: %d\n", totalLogs)
}

// getStatusIcon returns an icon for the status
func (f *StylishFormatter) getStatusIcon(status string) string {
	switch status {
	case "success":
		return "✓"
	case "failed":
		return "✗"
	case "running":
		return "⟳"
	case "pending":
		return "○"
	default:
		return "•"
	}
}

// colorizeStatus returns a colorized status string
func (f *StylishFormatter) colorizeStatus(status string) string {
	switch status {
	case "success", "completed":
		return color.GreenString(status)
	case "failed":
		return color.RedString(status)
	case "running":
		return color.YellowString(status)
	default:
		return status
	}
}

// formatTimestamp formats ISO timestamp to HH:MM:SS
func (f *StylishFormatter) formatTimestamp(timestamp string) string {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return timestamp
	}

	return t.Format("15:04:05")
}
