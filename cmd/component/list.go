package component

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
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

			return withStylishPagination(cmd, request)
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

func withStylishPagination(cmd *cobra.Command, request sdk.ApiComponentListRequest) error {
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
