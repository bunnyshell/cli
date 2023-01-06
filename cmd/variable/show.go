package variable

import (
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	var variableID string

	command := &cobra.Command{
		Use: "show",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().EnvironmentVariableApi.EnvironmentVariableView(ctx, variableID)

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	command.Flags().StringVar(&variableID, "id", variableID, "Environment Variable Id")
	command.MarkFlagRequired("id")

	mainCmd.AddCommand(command)
}
