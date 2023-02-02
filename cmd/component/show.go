package component

import (
	"bunnyshell.com/cli/pkg/api/component"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	itemOptions := component.NewItemOptions("")

	command := &cobra.Command{
		Use: "show",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			itemOptions.ID = settings.Profile.Context.ServiceComponent

			model, err := component.Get(itemOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	idFlag := options.ServiceComponent.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	mainCmd.AddCommand(command)
}
