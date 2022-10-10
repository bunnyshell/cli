package variable

import (
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
			var api = lib.GetAPI().EnvironmentVariableApi

			ctx, cancel := lib.GetContext()
			defer cancel()

			request := api.EnvironmentVariableList(ctx)

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

			resp, r, err := request.Execute()
			return lib.FormatRequestResult(cmd, resp, r, err)
		},
	}

	command.Flags().Int32Var(&page, "page", page, "Listing Page")
	command.Flags().StringVar(&name, "name", name, "Filter by Name")
	command.Flags().StringVar(organization, "organization", *organization, "Filter by Organization")
	command.Flags().StringVar(environment, "environment", *environment, "Filter by Environment")

	mainCmd.AddCommand(command)
}
