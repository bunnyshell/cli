package action

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	environment := &lib.CLIContext.Profile.Context.Environment

	command := &cobra.Command{
		Use: "start",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().EnvironmentApi.EnvironmentStart(ctx, *environment)

			resp, r, err := request.Execute()
			return lib.FormatRequestResult(cmd, resp, r, err)
		},
	}

	command.Flags().StringVar(environment, "id", *environment, "Environment Id")
	command.MarkFlagRequired("id")

	mainCmd.AddCommand(command)
}
