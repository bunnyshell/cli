package action

import (
	"io"
	"os"

	"bunnyshell.com/cli/pkg/api/project_variable"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	createOptions := project_variable.NewCreateOptions()

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
			createOptions.Project = settings.Profile.Context.Project

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

			model, err := project_variable.Create(createOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Project.AddFlagWithExtraHelp(
		"project",
		"Project for the variable",
		"Projects contain multiple variables and build settings",
		util.FlagRequired,
	))

	createOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
