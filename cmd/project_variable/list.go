package project_variable

import (
	"bunnyshell.com/cli/pkg/api/project_variable"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := project_variable.NewListOptions()

	command := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Organization = settings.Profile.Context.Organization
			listOptions.Project = settings.Profile.Context.Project

			return lib.ShowCollection(cmd, listOptions, func() (lib.ModelWithPagination, error) {
				return project_variable.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))

	listOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
