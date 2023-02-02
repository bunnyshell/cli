package progress

import (
	"time"

	"github.com/fatih/color"
)

type (
	UpdateStatus int

	PipelineStatus int

	InProgress bool
)

const (
	Failed UpdateStatus = iota
	Success

	Synced
)

const (
	StatusWorking PipelineStatus = iota
	StatusFinished
	StatusFailed
	StatusUnknown
)

const (
	defaultSpinnerUpdate = 150 * time.Millisecond
	defaultProgressSet   = 69 // ∙∙●
)

var statusMap = map[PipelineStatus]string{
	StatusWorking:  color.New(color.FgCyan).Sprintf("»"),
	StatusFinished: color.New(color.FgGreen).Sprintf("✔"),
	StatusFailed:   color.New(color.FgRed).Sprintf("✘"),
	StatusUnknown:  color.New(color.FgYellow).Sprintf("?"),
}
