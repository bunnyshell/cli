package component

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	component := &lib.CLIContext.Profile.Context.ServiceComponent

	command := &cobra.Command{
		Use: "show",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().ComponentApi.ComponentView(ctx, *component)

			resp, r, err := request.Execute()
			return lib.FormatRequestResult(cmd, resp, r, err)
		},
	}

	command.Flags().StringVar(component, "id", *component, "Component Id")
	command.MarkFlagRequired("id")

	mainCmd.AddCommand(command)
}
