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

	stopOptions := environment.NewStopOptions("")

	command := &cobra.Command{
		Use: "stop",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateActionOptions(&stopOptions.ActionOptions)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			stopOptions.ID = settings.Profile.Context.Environment

			event, err := environment.Stop(stopOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if stopOptions.WithoutPipeline {
				return lib.FormatCommandData(cmd, event)
			}

			printLogs := settings.IsStylish()

			if err = processEventPipeline(cmd, event, "stop", printLogs, stopOptions.Interval); err != nil {
				if printLogs {
					cmd.Printf("\nEnvironment %s stopping failed\n", stopOptions.ID)
				}

				return err
			}

			if printLogs {
				cmd.Printf("\nEnvironment %s successfully stopped\n", stopOptions.ID)
			}

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Environment.GetRequiredFlag("id"))

	stopOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
