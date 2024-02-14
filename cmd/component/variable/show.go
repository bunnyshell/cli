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

	itemOptions := component_variable.NewItemOptions("")

	command := &cobra.Command{
		Use:     "show",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			itemOptions.ID = settings.Profile.Context.ServiceComponent

			model, err := component_variable.Get(itemOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.ServiceComponent.GetRequiredFlag("id"))

	mainCmd.AddCommand(command)
}
