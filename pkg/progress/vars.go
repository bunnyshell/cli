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
	PipelineUnknownState
)

const (
	StatusSuccess    = "success"
	StatusFailed     = "failed"
	StatusInProgress = "in_progress"
	StatusPending    = "pending"
)

const (
	defaultSpinnerUpdate = 150 * time.Millisecond
	defaultProgressSet   = 69 // ∙∙●
)

var statusMap = map[PipelineStatus]string{
	PipelineWorking:      color.New(color.FgCyan).Sprintf("»"),
	PipelineFinished:     color.New(color.FgGreen).Sprintf("✔"),
	PipelineFailed:       color.New(color.FgRed).Sprintf("✘"),
	PipelineUnknownState: color.New(color.FgYellow).Sprintf("?"),
}

var ErrPipeline = errors.New("pipeline has encountered an error")
