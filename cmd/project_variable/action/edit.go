package action

import (
	"bunnyshell.com/cli/pkg/api/project_variable"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	editOptions := project_variable.NewEditOptions("")

	command := &cobra.Command{
		Use: "edit",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := project_variable.Edit(editOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(GetIDOption(&editOptions.ID).GetRequiredFlag("id"))

	editOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
