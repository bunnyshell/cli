package event

import (
	"net/http"

	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	var (
		page      int32
		eventType string
		status    string
	)

	organization := &lib.CLIContext.Profile.Context.Organization
	environment := &lib.CLIContext.Profile.Context.Environment

	command := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			return lib.ShowCollection(cmd, page, func(page int32) (lib.ModelWithPagination, *http.Response, error) {
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

				if eventType != "" {
					request = request.Type_(eventType)
				}

				if status != "" {
					request = request.Status(status)
				}

				return request.Execute()
			})
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")
	command.Flags().StringVar(organization, "organization", *organization, "Filter by Organization")
	command.Flags().StringVar(environment, "environment", *environment, "Filter by Environment")
	command.Flags().StringVar(&eventType, "type", eventType, "Filter by Type")
	command.Flags().StringVar(&status, "status", status, "Filter by Status")

	mainCmd.AddCommand(command)
}
