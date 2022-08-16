package lib

import (
	"encoding/json"
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
				return printTimeout(cmd)
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
	result := map[string]interface{}{
		"status": r.StatusCode,
		"error":  err.Error(),
	}

	var extra map[string]interface{}
	json.Unmarshal(err.Body(), &extra)
	if extra["detail"] != nil {
		result["extra"] = extra["detail"]
	}

	return FormatCommandData(cmd, result)
}

func printTimeout(cmd *cobra.Command) error {
	return FormatCommandData(cmd, map[string]interface{}{
		"status": 0,
		"error":  "Operation timed out",
	})
}
