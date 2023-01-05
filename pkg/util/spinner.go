package util

import (
	"github.com/briandowns/spinner"
)

func MakeSpinner(suffix string) *spinner.Spinner {
	spinner := spinner.New(spinner.CharSets[9], defaultDuration)
	spinner.Suffix = suffix

	return spinner
}
