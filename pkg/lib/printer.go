package lib

import (
	"fmt"
	"net"
	"net/http"

	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/formatter"
	"bunnyshell.com/sdk"
)

func FormatCommandError(cmd *cobra.Command, err error) error {
	FormatCommandData(cmd, map[string]interface{}{
		"error": err.Error(),
	})

	return err
}

func FormatCommandData(cmd *cobra.Command, data interface{}) error {
	result, err := formatter.Formatter(data, CLIContext.OutputFormat)
	if err != nil {
		cmd.PrintErrln(err)
		return err
	}

	cmd.Println(string(result))
	return nil
}

func FormatRequestResult(cmd *cobra.Command, data interface{}, r *http.Response, err error) error {
	if err != nil {
		switch err := err.(type) {
		case net.Error:
			if err.Timeout() {
				return printTimeout(cmd, err)
			}
		case *sdk.GenericOpenAPIError:
			return printOpenAPIError(cmd, err, r)
		}

		return printGenericError(cmd, err, r)
	}

	return FormatCommandData(cmd, data)
}

func printGenericError(cmd *cobra.Command, err error, r *http.Response) error {
	return FormatCommandData(cmd, map[string]interface{}{
		"status": r.StatusCode,
		"error":  err.Error(),
	})
}

func printOpenAPIError(cmd *cobra.Command, err *sdk.GenericOpenAPIError, r *http.Response) error {
	switch model := err.Model().(type) {
	case sdk.ProblemGeneric:
		return FormatCommandData(cmd, &model)
	}

	data := sdk.NewProblemGeneric()
	data.SetTitle(fmt.Sprintf("Response status: %d", r.StatusCode))
	data.SetDetail(err.Error())

	return FormatCommandData(cmd, data)
}

func printTimeout(cmd *cobra.Command, err error) error {
	data := sdk.NewProblemGeneric()
	data.SetTitle("Operation timed out")
	data.SetDetail(err.Error())
	return FormatCommandData(cmd, data)
}
