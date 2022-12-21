package component

import (
	"net/http"

	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	var page int32
	var clusterStatus string
	var operationStatus string

	organization := &lib.CLIContext.Profile.Context.Organization
	environment := &lib.CLIContext.Profile.Context.Environment
	project := &lib.CLIContext.Profile.Context.Project

	command := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return lib.ShowCollection(cmd, page, func(page int32) (lib.ModelWithPagination, *http.Response, error) {
				ctx, cancel := lib.GetContext()
				defer cancel()

				request := lib.GetAPI().ComponentApi.ComponentList(ctx)

				if page != 0 {
					request = request.Page(page)
				}

				if clusterStatus != "" {
					request = request.ClusterStatus(clusterStatus)
				}

				if operationStatus != "" {
					request = request.OperationStatus(operationStatus)
				}

				if *organization != "" {
					request = request.Organization(*organization)
				}

				if *environment != "" {
					request = request.Environment(*environment)
				}

				if *project != "" {
					request = request.Project(*project)
				}

				return request.Execute()
			})
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")
	command.Flags().StringVar(&clusterStatus, "clusterStatus", clusterStatus, "Filter by ClusterStatus")
	command.Flags().StringVar(&operationStatus, "operationStatus", operationStatus, "Filter by OperationStatus")
	command.Flags().StringVar(environment, "environment", *environment, "Filter by Environment")
	command.Flags().StringVar(project, "project", *project, "Filter by Project")
	command.Flags().StringVar(organization, "organization", *organization, "Filter by Organization")

	mainCmd.AddCommand(command)
}
