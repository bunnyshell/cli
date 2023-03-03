package repository

import (
	"bunnyshell.com/cli/pkg/api/template/repository"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := repository.NewListOptions()

	command := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Organization = settings.Profile.Context.Organization

			return lib.ShowCollection(cmd, listOptions, func() (lib.ModelWithPagination, error) {
				return repository.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))

	listOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
