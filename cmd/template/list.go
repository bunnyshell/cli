package template

import (
	"bunnyshell.com/cli/pkg/api/template"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	listOptions := template.NewListOptions()

	command := &cobra.Command{
		Use:     "list",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			listOptions.Organization = settings.Profile.Context.Organization

			if listOptions.Source == "public" {
				listOptions.Organization = ""
			}

			return lib.ShowCollection(cmd, listOptions, func() (lib.ModelWithPagination, error) {
				return template.List(listOptions)
			})
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))

	listOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
