package organization

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/sdk"
)

func init() {
	var page int32

	command := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := lib.GetContext()
			defer cancel()

			request := lib.GetAPI().OrganizationApi.OrganizationList(ctx)

			if page != 0 {
				request = request.Page(page)
			}

			return withStylishPagination(cmd, request)
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")

	mainCmd.AddCommand(command)
}

func withStylishPagination(cmd *cobra.Command, request sdk.ApiOrganizationListRequest) error {
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
