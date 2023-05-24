package template

import (
	"bunnyshell.com/cli/pkg/api/template"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	itemOptions := template.NewItemOptions("")

	command := &cobra.Command{
		Use:     "show",
		GroupID: mainGroup.ID,

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := template.Get(itemOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(getIDOption(&itemOptions.ID).GetRequiredFlag("id"))

	mainCmd.AddCommand(command)
}
