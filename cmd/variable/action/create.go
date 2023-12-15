package action

import (
	"bunnyshell.com/cli/pkg/api/variable"
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	createOptions := variable.NewCreateOptions()

	command := &cobra.Command{
		Use: "create",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			createOptions.Environment = settings.Profile.Context.Environment

			model, err := variable.Create(createOptions)
			if err != nil {
				return lib.FormatCommandError(cmd, err)
			}

			return lib.FormatCommandData(cmd, model)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Project.AddFlagWithExtraHelp(
		"environment",
		"Environment for the variable",
		"Environments contain multiple variables",
		util.FlagRequired,
	))

	createOptions.UpdateFlagSet(flags)

	mainCmd.AddCommand(command)
}
