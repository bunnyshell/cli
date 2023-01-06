package profile

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	command := &cobra.Command{
		Use: "remove",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.MainManager.RemoveProfile(settings.Profile.Name); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"message": "Profile removed",
				"data":    settings.Profile.Name,
			})
		},
	}

	flags := command.Flags()

	profileNameFlag := options.ProfileName.CloneMainFlag()
	flags.AddFlag(profileNameFlag)
	_ = command.MarkFlagRequired(profileNameFlag.Name)

	mainCmd.AddCommand(command)
}
