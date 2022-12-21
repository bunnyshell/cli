package environment

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
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

			return withStylishPagination(cmd, request)
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")
	command.Flags().StringVar(organization, "organization", *organization, "Filter by Organization")
	command.Flags().StringVar(project, "project", *project, "Filter by Project")
	command.Flags().StringVar(&type_, "type", type_, "Filter by Type")
	command.Flags().StringVar(&operationStatus, "operationStatus", operationStatus, "Filter by Operation Status")

	mainCmd.AddCommand(command)
}

func withStylishPagination(cmd *cobra.Command, request sdk.ApiEnvironmentListRequest) error {
	for {
		model, resp, err := request.Execute()
		if err = lib.FormatRequestResult(cmd, model, resp, err); err != nil {
			return err
		}

		page, err := lib.ProcessPagination(cmd, model)
		if err != nil {
			return err
		}

		if page == lib.PAGINATION_QUIT {
			return nil
		}

		request = request.Page(page)
	}
}
