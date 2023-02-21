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

	startOptions := environment.NewStartOptions("")

	command := &cobra.Command{
		Use: "start",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateActionOptions(&startOptions.ActionOptions)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			startOptions.ID = settings.Profile.Context.Environment

			event, err := environment.Start(startOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			if startOptions.WithoutPipeline {
				return lib.FormatCommandData(cmd, event)
			}

			if err = processEventPipeline(cmd, event, "start"); err != nil {
				cmd.Printf("\nEnvironment %s starting failed\n", startOptions.ID)

				return err
			}

			cmd.Printf("\nEnvironment %s successfully started\n", startOptions.ID)

			return showEnvironmentEndpoints(cmd, startOptions.ID)
		},
	}

	flags := command.Flags()

	idFlag := options.Environment.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	startOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
