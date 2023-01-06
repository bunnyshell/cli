package variable

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
		page int32
		name string
	)

	command := &cobra.Command{
		Use:     "list",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			return lib.ShowCollection(cmd, page, func(page int32) (lib.ModelWithPagination, *http.Response, error) {
				ctx, cancel := lib.GetContext()
				defer cancel()

				request := lib.GetAPI().EnvironmentVariableApi.EnvironmentVariableList(ctx)

				if page != 0 {
					request = request.Page(page)
				}

				if settings.Profile.Context.Organization != "" {
					request = request.Organization(settings.Profile.Context.Organization)
				}

				if settings.Profile.Context.Environment != "" {
					request = request.Environment(settings.Profile.Context.Environment)
				}

				if name != "" {
					request = request.Name(name)
				}

				return request.Execute()
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Environment.GetFlag("environment"))

	flags.Int32Var(&page, "page", page, "Listing Page")
	flags.StringVar(&name, "name", name, "Filter by Name")

	mainCmd.AddCommand(command)
}
