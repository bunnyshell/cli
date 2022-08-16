package organization

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	organization := &lib.CLIContext.Profile.Context.Organization

	command := &cobra.Command{
		Use: "show",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().OrganizationApi.OrganizationView(ctx, *organization)

			resp, r, err := request.Execute()
			return lib.FormatRequestResult(cmd, resp, r, err)
		},
	}

	command.Flags().StringVar(organization, "id", *organization, "Organization Id")
	command.MarkFlagRequired("id")

	mainCmd.AddCommand(command)
}
