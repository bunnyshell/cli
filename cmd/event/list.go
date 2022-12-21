package event

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

func init() {
	var page int32
	var type_ string
	var status string

	organization := &lib.CLIContext.Profile.Context.Organization
	environment := &lib.CLIContext.Profile.Context.Environment

	command := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().EventApi.EventList(ctx)

			if page != 0 {
				request = request.Page(page)
			}

			if *organization != "" {
				request = request.Organization(*organization)
			}

			if *environment != "" {
				request = request.Environment(*environment)
			}

			if type_ != "" {
				request = request.Type_(type_)
			}

			if status != "" {
				request = request.Status(status)
			}

			return withStylishPagination(cmd, request)
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")
	command.Flags().StringVar(organization, "organization", *organization, "Filter by Organization")
	command.Flags().StringVar(environment, "environment", *environment, "Filter by Environment")
	command.Flags().StringVar(&type_, "type", type_, "Filter by Type")
	command.Flags().StringVar(&status, "status", status, "Filter by Status")

	mainCmd.AddCommand(command)
}

func withStylishPagination(cmd *cobra.Command, request sdk.ApiEventListRequest) error {
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
