package environment

import (
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := environment.NewListOptions()

	command := &cobra.Command{
		Use:     "list",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		PersistentPreRunE: util.PersistentPreRunChain,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Organization = settings.Profile.Context.Organization
			listOptions.Project = settings.Profile.Context.Project

			return lib.ShowCollection(cmd, listOptions, func() (lib.ModelWithPagination, error) {
				return environment.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))

	listOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
