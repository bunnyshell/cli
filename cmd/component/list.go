package component

import (
	"net/http"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	var (
		page            int32
		clusterStatus   string
		operationStatus string
	)

	options := config.GetOptions()
	settings := config.GetSettings()

	command := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			return lib.ShowCollection(cmd, page, func(page int32) (lib.ModelWithPagination, *http.Response, error) {
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

				if settings.Profile.Context.Organization != "" {
					request = request.Organization(settings.Profile.Context.Organization)
				}

				if settings.Profile.Context.Environment != "" {
					request = request.Environment(settings.Profile.Context.Environment)
				}

				if settings.Profile.Context.Project != "" {
					request = request.Project(settings.Profile.Context.Project)
				}

				return request.Execute()
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))
	flags.AddFlag(options.Environment.GetFlag("environment"))

	flags.Int32Var(&page, "page", page, "Listing Page")
	flags.StringVar(&clusterStatus, "clusterStatus", clusterStatus, "Filter by ClusterStatus")
	flags.StringVar(&operationStatus, "operationStatus", operationStatus, "Filter by OperationStatus")

	mainCmd.AddCommand(command)
}
