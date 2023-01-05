package organization

import (
	"net/http"

	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	var page int32

	command := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			return lib.ShowCollection(cmd, page, func(page int32) (lib.ModelWithPagination, *http.Response, error) {
				ctx, cancel := lib.GetContext()
				defer cancel()

				request := lib.GetAPI().OrganizationApi.OrganizationList(ctx)

				if page != 0 {
					request = request.Page(page)
				}

				return request.Execute()
			})
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")

	mainCmd.AddCommand(command)
}
