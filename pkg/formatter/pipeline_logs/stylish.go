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

// Format outputs logs in stylish format
func (f *StylishFormatter) Format(logs *workflow_job.WorkflowJobLogs, w io.Writer) error {
	// Print header
	fmt.Fprintf(w, "\nWorkflow Job: %s\n", logs.WorkflowJobID)
	fmt.Fprintf(w, "Status: %s\n\n", f.colorizeStatus(logs.Status))

	// Print each step
	for _, step := range logs.Steps {
		f.printStep(w, &step)
	}

	// Print summary
	f.printSummary(w, logs)

	return nil
}

// printStep prints a single step with its logs
func (f *StylishFormatter) printStep(w io.Writer, step *workflow_job.LogStep) {
	// Step header with separator
	separator := strings.Repeat("━", 70)
	fmt.Fprintf(w, "%s\n", color.New(color.Faint).Sprint(separator))

	// Step name with status indicator
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

	// Print logs
	for _, log := range step.Logs {
		f.printLogMessage(w, &log)
	}

	fmt.Fprintln(w)
}

// printLogMessage prints a single log message
func (f *StylishFormatter) printLogMessage(w io.Writer, log *workflow_job.LogMessage) {
	// Format timestamp (HH:MM:SS)
	timestamp := f.formatTimestamp(log.Timestamp)
	timestampStr := color.New(color.Faint).Sprintf("  %s", timestamp)

	message := fmt.Sprintf("  %s", log.Message)

	fmt.Fprintf(w, "%s  %s\n", timestampStr, message)
}

// printSummary prints summary information
func (f *StylishFormatter) printSummary(w io.Writer, logs *workflow_job.WorkflowJobLogs) {
	separator := strings.Repeat("━", 70)
	fmt.Fprintf(w, "%s\n\n", color.New(color.Faint).Sprint(separator))

	// Count total logs
	totalLogs := 0
	for _, step := range logs.Steps {
		totalLogs += len(step.Logs)
	}

	fmt.Fprintf(w, "Pipeline %s\n", f.colorizeStatus(logs.Status))
	fmt.Fprintf(w, "Total log lines: %d\n", totalLogs)

	if logs.Pagination.HasMore {
		fmt.Fprintf(w, "\n%s\n", color.YellowString("⚠ More logs available (showing %d of %d)",
			logs.Pagination.Offset+totalLogs, logs.Pagination.Total))
	}
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
		// Fallback to original if parsing fails
		return timestamp
	}

	return t.Format("15:04:05")
}
