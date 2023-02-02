package variable

import (
	"bunnyshell.com/cli/pkg/api/variable"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := variable.NewListOptions()

	command := &cobra.Command{
		Use:     "list",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Organization = settings.Profile.Context.Organization
			listOptions.Environment = settings.Profile.Context.Environment

			return lib.ShowCollection(cmd, listOptions, func() (lib.ModelWithPagination, error) {
				return variable.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Environment.GetFlag("environment"))

	listOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
