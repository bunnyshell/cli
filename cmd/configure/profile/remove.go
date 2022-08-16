package profile

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

func init() {
	var profileName string

	removeProfileCommand := &cobra.Command{
		Use: "remove",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := lib.RemoveProfile(profileName); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, map[string]interface{}{
				"message": "Profile removed",
				"data":    profileName,
			})
		},
	}

	removeProfileCommand.Flags().StringVar(&profileName, "name", profileName, "Profile name to remove")
	removeProfileCommand.MarkFlagRequired("name")

	mainCmd.AddCommand(removeProfileCommand)
}
