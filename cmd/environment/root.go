package environment

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/cmd/environment/action"
	"bunnyshell.com/cli/pkg/lib"
	"bunnyshell.com/cli/pkg/util"
)

var mainCmd = &cobra.Command{
	Use:     "environments",
	Aliases: []string{"env"},

	Short: "Bunnyshell Environments",

	ValidArgsFunction: cobra.NoFileCompletions,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		lib.LoadViperConfigIntoContext()
	},
}

func init() {
	lib.CLIContext.RequireTokenOnCommand(mainCmd)

	util.AddGroupedCommands(
		mainCmd,
		cobra.Group{
			ID:    "actions",
			Title: "Environment Actions",
		},
		action.GetMainCommand().Commands(),
	)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
