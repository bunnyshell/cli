package progress

import (
	"time"

	"github.com/fatih/color"
)

type UpdateStatus int
type InProgress bool

const (
	Failed UpdateStatus = iota
	Success

	Synced
)

const (
	defaultUpdate = 100 * time.Millisecond

	prefix = "»"
)

var (
	prefixWait = color.New(color.FgCyan).Sprintf(prefix)
	prefixDone = color.New(color.FgGreen).Sprintf(prefix)
	prefixErr  = color.New(color.FgRed).Sprintf(prefix)
	prefixUnk  = color.New(color.FgYellow).Sprintf(prefix)
)
