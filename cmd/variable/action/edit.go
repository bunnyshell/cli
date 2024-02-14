package action

import (
	"io"
	"os"

	"bunnyshell.com/cli/pkg/api/variable"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	editOptions := variable.NewEditOptions("")

	command := &cobra.Command{
		Use: "edit",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			hasStdin, err := util.IsStdinPresent()
			if err != nil {
				return err
			}

			flags := cmd.Flags()
			if flags.Changed("value") && hasStdin {
				return errMultipleValueInputs
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			if flags.Changed("value") {
				editOptions.EnvironmentVariableEditAction.SetValue(flags.Lookup("value").Value.String())
			}

			hasStdin, err := util.IsStdinPresent()
			if err != nil {
				return err
			}

			if hasStdin {
				buf, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}

				editOptions.EnvironmentVariableEditAction.SetValue(string(buf))
			}

			model, err := variable.Edit(editOptions)
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
