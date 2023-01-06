package organization

import (
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	organization := &lib.CLIContext.Profile.Context.Organization

	command := &cobra.Command{
		Use: "show",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().OrganizationApi.OrganizationView(ctx, *organization)

			model, resp, err := request.Execute()

			return lib.FormatRequestResult(cmd, model, resp, err)
		},
	}

	command.Flags().StringVar(organization, "id", *organization, "Organization Id")
	command.MarkFlagRequired("id")

	mainCmd.AddCommand(command)
}
