package profile

import (
	"errors"

	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	profileName := ""

	command := &cobra.Command{
		Use: "default",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if errors.Is(config.MainManager.Error, config.ErrConfigLoad) {
				return config.MainManager.Error
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.MainManager.SetDefaultProfile(profileName); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if err := config.MainManager.Save(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"message": "Profile set as default",
				"data":    profileName,
			})
		},
	}

	flags := command.Flags()

	flags.StringVar(&profileName, "profile", profileName, "Profile name to set as default")
	util.MarkFlagRequiredWithHelp(flags.Lookup("profile"), "The local profile name to set as the default profile")

	mainCmd.AddCommand(command)
}

func setDefaultProfile(profile *config.Profile) error {
	return setDefaultProfileByName(profile.Name)
}

func setDefaultProfileByName(name string) error {
	return config.MainManager.SetDefaultProfile(name)
}
