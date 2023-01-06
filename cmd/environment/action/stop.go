package action

import (
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	environment := &lib.CLIContext.Profile.Context.Environment

	command := &cobra.Command{
		Use: "stop",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().EnvironmentApi.EnvironmentStop(ctx, *environment)

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	command.Flags().StringVar(environment, "id", *environment, "Environment Id")
	command.MarkFlagRequired("id")

	mainCmd.AddCommand(command)
}
