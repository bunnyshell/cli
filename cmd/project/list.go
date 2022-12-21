package project

import (
	"net/http"

	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	var page int32
	organization := &lib.CLIContext.Profile.Context.Organization

	command := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return lib.ShowCollection(cmd, page, func(page int32) (lib.ModelWithPagination, *http.Response, error) {
				ctx, cancel := lib.GetContext()
				defer cancel()

				request := lib.GetAPI().ProjectApi.ProjectList(ctx)

				if page != 0 {
					request = request.Page(page)
				}

				if *organization != "" {
					request = request.Organization(*organization)
				}

				return request.Execute()
			})
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")
	command.Flags().StringVar(organization, "organization", *organization, "Filter by organization")

	mainCmd.AddCommand(command)
}
