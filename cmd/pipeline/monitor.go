package pipeline

import (
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/progress"
	"github.com/spf13/cobra"
)

func init() {
	var pipelineID string

	progressOptions := progress.NewOptions()

	command := &cobra.Command{
		Use: "monitor",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: lib.OnlyStylish,

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := progress.Pipeline(pipelineID, progressOptions); err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(getIDOption(&pipelineID).GetRequiredFlag("id"))

	flags.DurationVar(&progressOptions.Interval, "interval", progressOptions.Interval, "Pipeline check interval")

	mainCmd.AddCommand(command)
}
