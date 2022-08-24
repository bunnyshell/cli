package util

import (
	"time"

	"github.com/briandowns/spinner"
)

func MakeSpinner(suffix string) *spinner.Spinner {
	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spinner.Suffix = suffix
	return spinner
}
