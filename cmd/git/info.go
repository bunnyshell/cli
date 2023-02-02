package git

import (
	"bunnyshell.com/cli/pkg/api/component/git"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := git.NewListOptions()

	command := &cobra.Command{
		Use: "info",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Organization = settings.Profile.Context.Organization
			listOptions.Project = settings.Profile.Context.Project
			listOptions.Environment = settings.Profile.Context.Environment

			return lib.ShowCollection(cmd, listOptions, func() (lib.ModelWithPagination, error) {
				return git.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Environment.GetFlag("environment"))
	flags.AddFlag(options.Project.GetFlag("project"))

	listOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
