package project_variable

import (
	"bunnyshell.com/cli/cmd/project_variable/action"
	"bunnyshell.com/cli/pkg/api/project_variable"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	itemOptions := project_variable.NewItemOptions("")

	command := &cobra.Command{
		Use: "show",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := project_variable.Get(itemOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(action.GetIDOption(&itemOptions.ID).GetRequiredFlag("id"))

	mainCmd.AddCommand(command)
}
