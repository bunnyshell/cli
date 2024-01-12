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

	cloneOptions := environment.NewCloneOptions("")

	command := &cobra.Command{
		Use: "clone",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			cloneOptions.ID = settings.Profile.Context.Environment

			model, err := environment.Clone(cloneOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			printLogs := settings.IsStylish()

			if printLogs {
				cmd.Printf("\nEnvironment %s successfully cloned:\n\n", cloneOptions.ID)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Environment.GetRequiredFlag("id"))

	cloneOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
