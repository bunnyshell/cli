package variable

import (
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

func init() {
	var (
		variableID string
		value      string
	)

	command := &cobra.Command{
		Use: "edit",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().EnvironmentVariableApi.EnvironmentVariableEdit(
				ctx,
				variableID,
			).EnvironmentVariableEdit(
				*toVariableEdit(value),
			)

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	command.Flags().StringVar(&variableID, "id", variableID, "Environment Variable Id")
	command.MarkFlagRequired("id")

	command.Flags().StringVar(&value, "value", value, "Environment Variable Value")
	command.MarkFlagRequired("value")

	mainCmd.AddCommand(command)
}

func toVariableEdit(value string) *sdk.EnvironmentVariableEdit {
	edit := sdk.NewEnvironmentVariableEdit()
	edit.SetValue(value)

	return edit
}
