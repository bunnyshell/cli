package variable

import (
	"bunnyshell.com/cli/pkg/api/component_variable"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := component_variable.NewListOptions()

	command := &cobra.Command{
		Use:     "list",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Organization = settings.Profile.Context.Organization
			listOptions.Project = settings.Profile.Context.Project
			listOptions.Environment = settings.Profile.Context.Environment
			listOptions.Component = settings.Profile.Context.ServiceComponent

			return lib.ShowCollection(cmd, listOptions, func() (lib.ModelWithPagination, error) {
				return component_variable.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))
	flags.AddFlag(options.Environment.GetFlag("environment"))
	flags.AddFlag(options.ServiceComponent.GetFlag("component"))

	listOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
