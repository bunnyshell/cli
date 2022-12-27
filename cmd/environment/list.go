package environment

import (
	"net/http"

	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	var page int32
	var type_ string
	var operationStatus string

	organization := &lib.CLIContext.Profile.Context.Organization
	project := &lib.CLIContext.Profile.Context.Project

	command := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return lib.ShowCollection(cmd, page, func(page int32) (lib.ModelWithPagination, *http.Response, error) {
				ctx, cancel := lib.GetContext()
				defer cancel()

				request := lib.GetAPI().EnvironmentApi.EnvironmentList(ctx)

				if page != 0 {
					request = request.Page(page)
				}

				if *organization != "" {
					request = request.Organization(*organization)
				}

				if *project != "" {
					request = request.Project(*project)
				}

				if type_ != "" {
					request = request.Type_(type_)
				}

				if operationStatus != "" {
					request = request.OperationStatus(operationStatus)
				}

				return request.Execute()
			})
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")
	command.Flags().StringVar(organization, "organization", *organization, "Filter by Organization")
	command.Flags().StringVar(project, "project", *project, "Filter by Project")
	command.Flags().StringVar(&type_, "type", type_, "Filter by Type")
	command.Flags().StringVar(&operationStatus, "operationStatus", operationStatus, "Filter by Operation Status")

	mainCmd.AddCommand(command)
}
