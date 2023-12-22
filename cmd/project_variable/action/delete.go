package action

import (
	"bunnyshell.com/cli/pkg/api/project_variable"
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

func init() {
	deleteOptions := project_variable.NewDeleteOptions()

	command := &cobra.Command{
		Use: "delete",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			err := project_variable.Delete(deleteOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			cmd.Printf("\nProject %s successfully deleted\n", deleteOptions.ID)

			return nil
		},
	}

	flags := command.Flags()

	flags.AddFlag(GetIDOption(&deleteOptions.ID).GetRequiredFlag("id"))

	mainCmd.AddCommand(command)
}
