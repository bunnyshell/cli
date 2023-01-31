package api

import (
	"fmt"
	"net"
	"net/http"

	"bunnyshell.com/sdk"
)

type Error struct {
	Title  string `json:"title" yaml:"title"`
	Detail string `json:"detail" yaml:"detail"`
}

func (pe Error) Error() string {
	return pe.Title + ": " + pe.Detail
}

func ParseError(resp *http.Response, err error) error {
	switch err := err.(type) {
	case net.Error:
		if err.Timeout() {
			return Error{
				Title:  "Operation timed out",
				Detail: err.Error(),
			}
		}
	case *sdk.GenericOpenAPIError:
		problem, isProblem := err.Model().(sdk.ProblemGeneric)
		if isProblem {
			return Error{
				Title:  *problem.Title,
				Detail: *problem.Detail,
			}
		}
	}

	if resp == nil {
		return err
	}

	return Error{
		Title:  fmt.Sprintf("Response Status: %d", resp.StatusCode),
		Detail: err.Error(),
	}
}
