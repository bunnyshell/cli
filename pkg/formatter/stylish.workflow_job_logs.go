package formatter

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	wfstatus "bunnyshell.com/cli/pkg/api/workflow_job/status"
	"bunnyshell.com/sdk"
	"github.com/fatih/color"
)

// WorkflowJobLogsResult represents logs for a single job
type WorkflowJobLogsResult struct {
	JobID string
	Logs  *sdk.WorkflowJobWorkflowJobLogsOutputItem
}

const separatorWidth = 80

func statusIcon(status string) string {
	switch status {
	case wfstatus.JobSuccess: // same value as StepSuccess
		return color.New(color.FgGreen).Sprint("✓")
	case wfstatus.JobFailed, wfstatus.JobAbortFailed: // JobFailed same value as StepFailed
		return color.New(color.FgRed).Sprint("✘")
	case wfstatus.JobAborting, wfstatus.JobAborted:
		return color.New(color.FgRed, color.Bold).Sprint("⊘")
	case wfstatus.JobInProgress:
		return color.New(color.FgCyan).Sprint("▶︎")
	case wfstatus.JobPending, wfstatus.JobQueued:
		return color.New(color.FgWhite).Sprint("⋯")
	case wfstatus.JobSkipped:
		return color.New(color.FgHiBlack).Sprint("»")
	default:
		return color.New(color.FgWhite).Sprint("?")
	}
}

func statusNameColor(status string) *color.Color {
	switch status {
	case wfstatus.JobSuccess: // same value as StepSuccess
		return color.New(color.FgGreen, color.Bold)
	case wfstatus.JobFailed, wfstatus.JobAbortFailed: // JobFailed same value as StepFailed
		return color.New(color.FgRed, color.Bold)
	case wfstatus.JobAborting, wfstatus.JobAborted:
		return color.New(color.FgRed)
	case wfstatus.JobInProgress:
		return color.New(color.FgCyan, color.Bold)
	case wfstatus.JobPending, wfstatus.JobQueued:
		return color.New(color.FgWhite, color.Bold)
	case wfstatus.JobSkipped:
		return color.New(color.FgHiBlack, color.Bold)
	default:
		return color.New(color.FgWhite, color.Bold)
	}
}

func formatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05Z07:00")
}

func isStepFailed(status string) bool {
	return status == wfstatus.StepFailed || status == wfstatus.JobAbortFailed
}

func tabulateWorkflowJobLogs(writer *tabwriter.Writer, data []WorkflowJobLogsResult) {
	dim := color.New(color.FgHiBlack)

	for _, result := range data {
		if result.Logs == nil {
			continue
		}

		// Job header
		if job, ok := result.Logs.GetWorkflowJobOk(); ok && job != nil {
			jobName := getJobDisplayName(job, result.JobID)
			status := ""
			if s, ok := job.GetStatusOk(); ok && s != nil {
				status = *s
			}

			jobSep := strings.Repeat("═", separatorWidth+4) // 4 for the step indentation below
			fmt.Fprintf(writer, "\n%s\n", jobSep)
			fmt.Fprintf(writer, "  %s  Job: %s\n", statusIcon(status), statusNameColor(status).Sprint(jobName))

			// Row 1: Status + JobId
			var row1 []string
			if status != "" {
				row1 = append(row1, fmt.Sprintf("Status: %s", status))
			}
			if id, ok := job.GetIdOk(); ok && id != nil {
				row1 = append(row1, fmt.Sprintf("JobId: %s", *id))
			}
			if len(row1) > 0 {
				fmt.Fprintf(writer, "     %s\n", dim.Sprint(strings.Join(row1, "  ")))
			}

			// Row 2: Type + AllowedToFail
			var row2 []string
			if jobType, ok := job.GetTypeOk(); ok && jobType != nil {
				row2 = append(row2, fmt.Sprintf("Type: %s", *jobType))
			}
			if allowedToFail, ok := job.GetAllowedToFailOk(); ok && allowedToFail != nil {
				row2 = append(row2, fmt.Sprintf("AllowedToFail: %v", *allowedToFail))
			}
			if len(row2) > 0 {
				fmt.Fprintf(writer, "     %s\n", dim.Sprint(strings.Join(row2, "  ")))
			}

			// Row 3: StartedAt + Duration
			var row3 []string
			if startedAt, ok := job.GetStartedAtOk(); ok && startedAt != nil {
				row3 = append(row3, fmt.Sprintf("StartedAt: %s", formatDateTime(*startedAt)))
			}
			if duration, ok := job.GetDurationOk(); ok && duration != nil {
				row3 = append(row3, fmt.Sprintf("Duration: %ds", *duration))
			}
			if len(row3) > 0 {
				fmt.Fprintf(writer, "     %s\n", dim.Sprint(strings.Join(row3, "  ")))
			}

			fmt.Fprintf(writer, "%s\n", jobSep)
		}

		// Steps (indented inside the job)
		if steps, ok := result.Logs.GetStepsOk(); ok && steps != nil {
			for _, step := range steps {
				stepName := "Unknown"
				if name, ok := step.GetNameOk(); ok && name != nil {
					stepName = *name
				}

				stepStatus := ""
				if s, ok := step.GetStatusOk(); ok && s != nil {
					stepStatus = *s
				}

				// Step header
				stepSep := strings.Repeat("━", separatorWidth)
				fmt.Fprintf(writer, "\n    %s\n", stepSep)

				var stepHeader strings.Builder
				stepHeader.WriteString(fmt.Sprintf("    %s  Step: %s", statusIcon(stepStatus), statusNameColor(stepStatus).Sprint(stepName)))
				if stepStatus != wfstatus.StepSuccess {
					stepHeader.WriteString(fmt.Sprintf("  %s", stepStatus))
					if exitCode, ok := step.GetExitCodeOk(); ok && exitCode != nil {
						stepHeader.WriteString(fmt.Sprintf("  (exit: %d)", *exitCode))
					}
				}
				fmt.Fprintln(writer, stepHeader.String())

				fmt.Fprintf(writer, "    %s\n", stepSep)

				// Log lines: timestamp (white or red for failed) + two spaces + message
				tsColor := color.New(color.FgWhite, color.Bold)
				if isStepFailed(stepStatus) {
					tsColor = color.New(color.FgRed)
				}

				// Continuation indent aligns with the message column:
				// 4 (step indent) + 12 (timestamp "15:04:05.000") + 2 (spaces) = 18
				const logIndent = "    "
				const logContinuation = "                  " // 18 spaces

				if logs, ok := step.GetLogsOk(); ok && logs != nil && len(logs) > 0 {
					for _, logEntry := range logs {
						message := ""
						if msg, ok := logEntry.GetLogOk(); ok && msg != nil {
							message = *msg
						}

						if t, ok := logEntry.GetTimeOk(); ok && t != nil {
							message = strings.ReplaceAll(message, "\n", "\n"+logContinuation)
							fmt.Fprintf(writer, "%s%s  %s\n", logIndent, tsColor.Sprint(t.Format("15:04:05.000")), message)
						} else {
							message = strings.ReplaceAll(message, "\n", "\n"+logIndent)
							fmt.Fprintf(writer, "%s%s\n", logIndent, message)
						}
					}
				} else {
					fmt.Fprintf(writer, "%s%s\n", logIndent, dim.Sprint("(no logs)"))
				}
			}
		} else {
			fmt.Fprintf(writer, "\n  %s\n", dim.Sprint("No steps found"))
		}

		fmt.Fprintf(writer, "\n")
	}
}

func getJobDisplayName(job interface{}, fallbackID string) string {
	if j, ok := job.(interface{ GetName() string }); ok {
		if name := j.GetName(); name != "" {
			return name
		}
	}
	if j, ok := job.(interface{ GetId() string }); ok {
		if id := j.GetId(); id != "" {
			return id
		}
	}
	return fallbackID
}
