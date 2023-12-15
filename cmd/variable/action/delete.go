package action

import (
	"bunnyshell.com/cli/pkg/api/variable"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	settings := config.GetSettings()

	deleteOptions := variable.NewDeleteOptions()

	command := &cobra.Command{
		Use: "delete",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			deleteOptions.ID = settings.Profile.Context.Project

			err := variable.Delete(deleteOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			cmd.Printf("\nEnvironment variable %s successfully deleted\n", deleteOptions.ID)

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(GetIDOption(&deleteOptions.ID).GetRequiredFlag("id"))

	mainCmd.AddCommand(command)
}
