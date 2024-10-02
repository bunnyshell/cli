package action

import (
	"io"
	"os"

	"bunnyshell.com/cli/pkg/api/variable_group"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	createOptions := variable_group.NewCreateOptions()

	command := &cobra.Command{
		Use: "create",

		ValidArgsFunction: cobra.NoFileCompletions,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			hasStdin, err := util.IsStdinPresent()
			if err != nil {
				return err
			}

			flags := cmd.Flags()
			if !flags.Changed("value") && !hasStdin {
				return errMissingValue
			}

			if flags.Changed("value") && hasStdin {
				return errMultipleValueInputs
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			createOptions.Environment = settings.Profile.Context.Environment

			hasStdin, err := util.IsStdinPresent()
			if err != nil {
				return err
			}

			if hasStdin {
				buf, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}

				createOptions.Value = string(buf)
			}

			model, err := variable_group.Create(createOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Environment.AddFlagWithExtraHelp(
		"environment",
		"Environment for the variable",
		"Environments contain multiple variables",
		util.FlagRequired,
	))

	createOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
