package pipeline_logs

import (
	"fmt"
	"io"

	"bunnyshell.com/cli/pkg/api/workflow_job"
)

// RawFormatter formats logs as plain text (messages only)
type RawFormatter struct{}

// NewRawFormatter creates a new raw formatter
func NewRawFormatter() *RawFormatter {
	return &RawFormatter{}
}

// Format outputs pipeline logs in raw format (just messages, no formatting)
func (f *RawFormatter) Format(logs *workflow_job.PipelineLogs, w io.Writer) error {
	for _, job := range logs.Jobs {
		for _, step := range job.Steps {
			for _, log := range step.Logs {
				fmt.Fprintln(w, log.Message)
			}
		}
	}

	return nil
}
