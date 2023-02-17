package lib

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/formatter"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

var errGeneric = errors.New("oops! Something went wrong")

func FormatCommandError(cmd *cobra.Command, err error) error {
	_ = FormatCommandData(cmd, err)

	return errGeneric
}

func FormatCommandData(cmd *cobra.Command, data interface{}) error {
	result, err := formatter.Formatter(data, config.GetSettings().OutputFormat)
	if err != nil {
		cmd.PrintErrln(err)

		return err
	}

	cmd.Println(string(result))

	return nil
}

func FormatRequestResult(cmd *cobra.Command, data interface{}, resp *http.Response, err error) error {
	if err != nil {
		switch err := err.(type) {
		case net.Error:
			if err.Timeout() {
				return printTimeout(cmd, err)
			}
		case *sdk.GenericOpenAPIError:
			return printOpenAPIError(cmd, err, resp)
		}

		return printGenericError(cmd, err, resp)
	}

	return FormatCommandData(cmd, data)
}

func printGenericError(cmd *cobra.Command, err error, resp *http.Response) error {
	result := map[string]interface{}{
		"error": err.Error(),
	}

	if resp != nil {
		result["status"] = resp.StatusCode
	}

	return FormatCommandData(cmd, result)
}

func printOpenAPIError(cmd *cobra.Command, err *sdk.GenericOpenAPIError, resp *http.Response) error {
	switch model := err.Model().(type) {
	case sdk.ProblemGeneric:
		return FormatCommandData(cmd, &model)
	}

	data := sdk.NewProblemGeneric()
	data.SetTitle(fmt.Sprintf("Response status: %d", resp.StatusCode))
	data.SetDetail(err.Error())

	return FormatCommandData(cmd, data)
}

func printTimeout(cmd *cobra.Command, err error) error {
	data := sdk.NewProblemGeneric()
	data.SetTitle("Operation timed out")
	data.SetDetail(err.Error())

	return FormatCommandData(cmd, data)
}
