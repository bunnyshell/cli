package action

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

	editSettingsOptions := environment.NewEditSettingsOptions("")

	command := &cobra.Command{
		Use: "update-settings",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			editSettingsOptions.ID = settings.Profile.Context.Environment

			environmentModel, err := environment.Get(&editSettingsOptions.ItemOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}
			editSettingsOptions.UpdateEditSettingsForType(environmentModel.GetType())

			model, err := environment.EditSettings(editSettingsOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Environment.GetFlag("id", util.FlagRequired))

	editSettingsOptions.UpdateCommandFlags(command)

	mainCmd.AddCommand(command)
}
