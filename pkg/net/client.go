package net

import (
	"net/http"

	"github.com/briandowns/spinner"
)

type SpinnerTransport struct {
	Disabled bool

	Proxied http.RoundTripper
}

var DefaultSpinnerTransport = SpinnerTransport{
	Disabled: false,
	Proxied:  http.DefaultTransport,
}

func (st SpinnerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !st.Disabled {
		spinner := MakeSpinner()

		spinner.Start()

		defer spinner.Stop()
	}

	return st.Proxied.RoundTrip(req)
}

func GetCLIClient() *http.Client {
	return &http.Client{
		Transport: DefaultSpinnerTransport,
	}
}

func PauseSpinner() func() {
	prev := DefaultSpinnerTransport.Disabled
	DefaultSpinnerTransport.Disabled = true

	return func() {
		DefaultSpinnerTransport.Disabled = prev
	}
}

func MakeSpinner() *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], defaultDuration)
	s.Suffix = " Fetching API data..."

	return s
}
