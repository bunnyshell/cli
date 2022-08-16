package organization

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
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

			resp, r, err := request.Execute()
			return lib.FormatRequestResult(cmd, resp, r, err)
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")

	mainCmd.AddCommand(command)
}
