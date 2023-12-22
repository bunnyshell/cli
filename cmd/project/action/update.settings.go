package action

import (
	"bunnyshell.com/cli/pkg/api/project"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	editSettingsOptions := project.NewEditSettingsOptions("")

	command := &cobra.Command{
		Use: "update-settings",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			editSettingsOptions.ID = settings.Profile.Context.Project

			model, err := project.EditSettings(editSettingsOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Project.GetFlag("id", util.FlagRequired))

	editSettingsOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
