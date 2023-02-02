package project

import (
	"bunnyshell.com/cli/pkg/api/project"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := project.NewListOptions()

	command := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Organization = settings.Profile.Context.Organization

			return lib.ShowCollectionNoResponse(cmd, listOptions.Page, func(page int32) (lib.ModelWithPagination, error) {
				listOptions.Page = page

				return project.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))

	listOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
