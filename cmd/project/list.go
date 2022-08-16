package project

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	var page int32
	organization := &lib.CLIContext.Profile.Context.Organization

	command := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().ProjectApi.ProjectList(ctx)

			if page != 0 {
				request = request.Page(page)
			}

			if *organization != "" {
				request = request.Organization(*organization)
			}

			resp, r, err := request.Execute()
			return lib.FormatRequestResult(cmd, resp, r, err)
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")
	command.Flags().StringVar(organization, "organization", *organization, "Filter by organization")

	mainCmd.AddCommand(command)
}
