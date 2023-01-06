package event

import (
	"net/http"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	var (
		page      int32
		eventType string
		status    string
	)

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

				if settings.Profile.Context.Organization != "" {
					request = request.Organization(settings.Profile.Context.Organization)
				}

				if settings.Profile.Context.Environment != "" {
					request = request.Environment(settings.Profile.Context.Environment)
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

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Environment.GetFlag("environment"))

	flags.Int32Var(&page, "page", page, "Listing Page")
	flags.StringVar(&eventType, "type", eventType, "Filter by Type")
	flags.StringVar(&status, "status", status, "Filter by Status")

	mainCmd.AddCommand(command)
}
