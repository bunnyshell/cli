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

			if deleteOptions.WithPipeline {
				return processEventPipeline(cmd, event, "delete")
			}

			return lib.FormatCommandData(cmd, event)
		},
	}

	flags := command.Flags()

	idFlag := options.Environment.GetFlag("id")
	flags.AddFlag(idFlag)
	_ = command.MarkFlagRequired(idFlag.Name)

	deleteOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
