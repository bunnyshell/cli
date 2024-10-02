package variable_group

import (
	"bunnyshell.com/cli/pkg/api/variable_group"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := variable_group.NewListOptions()

	command := &cobra.Command{
		Use: "list",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, _ []string) error {
			listOptions.Organization = settings.Profile.Context.Organization
			listOptions.Environment = settings.Profile.Context.Environment

			return lib.ShowCollection(cmd, listOptions, func() (lib.ModelWithPagination, error) {
				return variable_group.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Environment.GetFlag("environment"))

	listOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
