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

			if err = processEventPipeline(cmd, event, "stop"); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			cmd.Printf("\nEnvironment %s successfully stopped\n", stopOptions.ID)

			return nil
		},
	}

	flags := command.Flags()

	idFlag := options.Environment.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	stopOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
