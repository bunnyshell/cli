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
		Use: "remove",

		ValidArgsFunction: cobra.NoFileCompletions,

		PersistentPreRunE: util.PersistentPreRunChain,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if errors.Is(config.MainManager.Error, config.ErrConfigLoad) {
				return config.MainManager.Error
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := removeProfileByName(profileName); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if err := config.MainManager.Save(); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"message": "Profile removed",
				"data":    profileName,
			})
		},
	}

	flags := command.Flags()

	flags.StringVar(&profileName, "profile", profileName, "Profile name to remove")
	util.MarkFlagRequiredWithHelp(flags.Lookup("profile"), "The local profile name to remove from available profiles")

	mainCmd.AddCommand(command)
}

func removeProfileByName(name string) error {
	return config.MainManager.RemoveProfile(name)
}
