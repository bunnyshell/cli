package profile

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	var profileName string

	defaultProfileCommand := &cobra.Command{
		Use: "default",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := lib.SetDefaultProfile(profileName); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"message": "Profile set as default",
				"data":    profileName,
			})
		},
	}

	defaultProfileCommand.Flags().StringVar(&profileName, "name", profileName, "Default profile for future api calls")
	defaultProfileCommand.MarkFlagRequired("name")

	mainCmd.AddCommand(defaultProfileCommand)
}
