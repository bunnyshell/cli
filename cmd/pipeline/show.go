package pipeline

import (
	"bunnyshell.com/cli/pkg/api/pipeline"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	var pipelineID string

	itemOptions := pipeline.NewItemOptions("")

	command := &cobra.Command{
		Use: "show",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			itemOptions.ID = pipelineID

			model, err := pipeline.Get(itemOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	idFlagName := "id"
	flags.StringVar(&pipelineID, idFlagName, pipelineID, "Pipeline Id")
	_ = command.MarkFlagRequired(idFlagName)

	mainCmd.AddCommand(command)
}
