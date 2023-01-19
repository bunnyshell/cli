package environment

import (
	"net/http"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	var (
		page            int32
		environmentType string
		operationStatus string
	)

	command := &cobra.Command{
		Use:     "list",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		PersistentPreRunE: util.PersistentPreRunChain,

		RunE: func(cmd *cobra.Command, args []string) error {
			return lib.ShowCollection(cmd, page, func(page int32) (lib.ModelWithPagination, *http.Response, error) {
				ctx, cancel := lib.GetContext()
				defer cancel()

				request := lib.GetAPI().EnvironmentApi.EnvironmentList(ctx)

				if page != 0 {
					request = request.Page(page)
				}

				if settings.Profile.Context.Organization != "" {
					request = request.Organization(settings.Profile.Context.Organization)
				}

				if settings.Profile.Context.Project != "" {
					request = request.Project(settings.Profile.Context.Project)
				}

				if environmentType != "" {
					request = request.Type_(environmentType)
				}

				if operationStatus != "" {
					request = request.OperationStatus(operationStatus)
				}

				return request.Execute()
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))

	flags.Int32Var(&page, "page", page, "Listing Page")
	flags.StringVar(&environmentType, "type", environmentType, "Filter by Type")
	flags.StringVar(&operationStatus, "operationStatus", operationStatus, "Filter by Operation Status")

	mainCmd.AddCommand(command)
}
