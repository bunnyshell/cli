package progress

import "github.com/briandowns/spinner"

const defaultProgressSet = 36 // [===>         ]

func newSpinner() *spinner.Spinner {
	return spinner.New(spinner.CharSets[defaultProgressSet], defaultUpdate)
}
