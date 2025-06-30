package action

import (
	"bunnyshell.com/cli/pkg/api/environment"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	deleteOptions := environment.NewDeleteOptions("")

	command := &cobra.Command{
		Use: "delete",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateActionOptions(&deleteOptions.ActionOptions)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			deleteOptions.ID = settings.Profile.Context.Environment

			event, err := environment.Delete(deleteOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if deleteOptions.WithoutPipeline {
				return lib.FormatCommandData(cmd, event)
			}

			if err = processEventPipeline(cmd, event, "delete", settings.IsStylish(), deleteOptions.Interval); err != nil {
				cmd.Printf("\nEnvironment %s deletion failed\n", deleteOptions.ID)

				return err
			}

			cmd.Printf("\nEnvironment %s successfully deleted\n", deleteOptions.ID)

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Environment.GetRequiredFlag("id"))

	deleteOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
