package component

import (
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	component := &lib.CLIContext.Profile.Context.ServiceComponent

	command := &cobra.Command{
		Use: "show",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().ComponentApi.ComponentView(ctx, *component)

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	command.Flags().StringVar(component, "id", *component, "Component Id")
	command.MarkFlagRequired("id")

	mainCmd.AddCommand(command)
}
