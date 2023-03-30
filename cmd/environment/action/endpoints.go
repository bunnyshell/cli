package action

import (
	"bunnyshell.com/cli/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	command := &cobra.Command{
		Use:     "endpoints",
		Aliases: []string{"end"},

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			return showEnvironmentEndpoints(cmd, settings.Profile.Context.Environment)
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Environment.GetRequiredFlag("id"))

	mainCmd.AddCommand(command)
}
