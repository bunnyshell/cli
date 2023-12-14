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

	editBuildSettingsOptions := project.NewEditBuildSettingsOptions("")

	command := &cobra.Command{
		Use: "update-build-settings",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			editBuildSettingsOptions.ID = settings.Profile.Context.Project

			model, err := project.EditBuildSettings(editBuildSettingsOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Project.GetFlag("id", util.FlagRequired))

	editBuildSettingsOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
