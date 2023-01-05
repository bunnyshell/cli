package variable

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

func init() {
	var id string
	var value string

	command := &cobra.Command{
		Use: "edit",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			var api = lib.GetAPI().EnvironmentVariableApi

			ctx, cancel := lib.GetContext()
			defer cancel()

			edit := *sdk.NewEnvironmentVariableEdit()
			edit.SetValue(value)

			request := api.EnvironmentVariableEdit(ctx, id).EnvironmentVariableEdit(edit)

			resp, r, err := request.Execute()
			return lib.FormatRequestResult(cmd, resp, r, err)
		},
	}

	command.Flags().StringVar(&id, "id", id, "Environment Variable Id")
	command.MarkFlagRequired("id")

	command.Flags().StringVar(&value, "value", value, "Environment Variable Value")
	command.MarkFlagRequired("value")

	mainCmd.AddCommand(command)
}
