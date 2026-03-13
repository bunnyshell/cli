package progress

import (
	"errors"
	"time"

	"github.com/fatih/color"
)

type (
	UpdateStatus int

	PipelineStatus int
)

const (
	Failed UpdateStatus = iota
	Success

	Synced
)

const (
	PipelineWorking PipelineStatus = iota
	PipelineFinished
	PipelineFailed
	PipelineAborted
	PipelineUnknownState
)

const (
	defaultSpinnerUpdate = 2000 * time.Millisecond
	defaultProgressSet   = 69 // ∙∙●
)

var statusMap = map[PipelineStatus]string{
	PipelineWorking:      color.New(color.FgCyan).Sprintf("»"),
	PipelineFinished:     color.New(color.FgGreen).Sprintf("✔"),
	PipelineFailed:       color.New(color.FgRed).Sprintf("✘"),
	PipelineAborted:      color.New(color.FgRed, color.Bold).Sprintf("⊘"),
	PipelineUnknownState: color.New(color.FgYellow).Sprintf("?"),
}

var ErrPipeline = errors.New("pipeline has encountered an error")
var ErrPipelineAborted = errors.New("pipeline was aborted")
