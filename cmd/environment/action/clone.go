package action

import (
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
	"github.com/spf13/cobra"
)

func init() {
	cloneName := ""
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

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	command.Flags().StringVar(environment, "id", *environment, "Environment Id")
	command.MarkFlagRequired("id")

	command.Flags().StringVar(&cloneName, "name", cloneName, "Environment Clone Name")
	command.MarkFlagRequired("name")

	mainCmd.AddCommand(command)
}
