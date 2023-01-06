package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	var (
		page            int32
		environmentType string
		operationStatus string
	)

	organization := &lib.CLIContext.Profile.Context.Organization
	project := &lib.CLIContext.Profile.Context.Project

	command := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

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

				if environmentType != "" {
					request = request.Type_(environmentType)
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
	command.Flags().StringVar(&environmentType, "type", environmentType, "Filter by Type")
	command.Flags().StringVar(&operationStatus, "operationStatus", operationStatus, "Filter by Operation Status")

	mainCmd.AddCommand(command)
}
