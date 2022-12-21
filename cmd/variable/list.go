package variable

import (
	"net/http"

	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	var page int32
	var name string

	organization := &lib.CLIContext.Profile.Context.Organization
	environment := &lib.CLIContext.Profile.Context.Environment

	command := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return lib.ShowCollection(cmd, page, func(page int32) (lib.ModelWithPagination, *http.Response, error) {
				ctx, cancel := lib.GetContext()
				defer cancel()

				request := lib.GetAPI().EnvironmentVariableApi.EnvironmentVariableList(ctx)

				if page != 0 {
					request = request.Page(page)
				}

				if *organization != "" {
					request = request.Organization(*organization)
				}

				if *environment != "" {
					request = request.Environment(*environment)
				}

				if name != "" {
					request = request.Name(name)
				}

				return request.Execute()
			})
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")
	command.Flags().StringVar(&name, "name", name, "Filter by Name")
	command.Flags().StringVar(organization, "organization", *organization, "Filter by Organization")
	command.Flags().StringVar(environment, "environment", *environment, "Filter by Environment")

	mainCmd.AddCommand(command)
}
