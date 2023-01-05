package action

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

func init() {
	var cloneName string

	environment := &lib.CLIContext.Profile.Context.Environment

	command := &cobra.Command{
		Use: "clone",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().EnvironmentApi.EnvironmentClone(
				ctx,
				*environment,
			).EnvironmentCloneAction(
				*sdk.NewEnvironmentCloneAction(cloneName),
			)

			resp, r, err := request.Execute()
			return lib.FormatRequestResult(cmd, resp, r, err)
		},
	}

	command.Flags().StringVar(environment, "id", *environment, "Environment Id")
	command.MarkFlagRequired("id")

	command.Flags().StringVar(&cloneName, "name", cloneName, "Environment Clone Name")
	command.MarkFlagRequired("name")

	mainCmd.AddCommand(command)
}
