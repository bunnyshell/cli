package template

import (
	"bunnyshell.com/cli/pkg/api/template"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	var templateID string

	itemOptions := template.NewItemOptions("")

	command := &cobra.Command{
		Use:     "show",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			itemOptions.ID = templateID

			model, err := template.Get(itemOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	idFlagName := "id"
	flags.StringVar(&templateID, idFlagName, templateID, "Template Id")
	_ = command.MarkFlagRequired(idFlagName)

	mainCmd.AddCommand(command)
}
