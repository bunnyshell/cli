package pipeline_logs

import (
	"io"

	"bunnyshell.com/cli/pkg/api/workflow_job"
	"gopkg.in/yaml.v3"
)

// YAMLFormatter formats logs as YAML
type YAMLFormatter struct{}

// NewYAMLFormatter creates a new YAML formatter
func NewYAMLFormatter() *YAMLFormatter {
	return &YAMLFormatter{}
}

// Format outputs pipeline logs in YAML format
func (f *YAMLFormatter) Format(logs *workflow_job.PipelineLogs, w io.Writer) error {
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	defer encoder.Close()

	return encoder.Encode(logs)
}
