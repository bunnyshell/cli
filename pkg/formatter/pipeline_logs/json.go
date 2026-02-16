package pipeline_logs

import (
	"encoding/json"
	"io"

	"bunnyshell.com/cli/pkg/api/workflow_job"
)

// JSONFormatter formats logs as JSON
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format outputs pipeline logs in JSON format
func (f *JSONFormatter) Format(logs *workflow_job.PipelineLogs, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(logs)
}
