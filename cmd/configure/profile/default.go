package profile

import (
	"errors"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	command := &cobra.Command{
		Use: "default",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			if errors.Is(config.MainManager.Error, config.ErrConfigLoad) {
				return config.MainManager.Error
			}

			if err := setDefaultProfile(&settings.Profile); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if err := config.MainManager.Save(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"message": "Profile set as default",
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

func setDefaultProfile(profile *config.Profile) error {
	return config.MainManager.SetDefaultProfile(profile.Name)
}
