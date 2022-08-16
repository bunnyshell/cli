package environment

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/cmd/environment/action"
	"bunnyshell.com/cli/pkg/lib"
)

var mainCmd = &cobra.Command{
	Use:   "environments",
	Short: "Bunnyshell Environments",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		lib.LoadViperConfigIntoContext()
	},
}

func init() {
	lib.CLIContext.RequireTokenOnCommand(mainCmd)

	for _, command := range action.GetMainCommand().Commands() {
		mainCmd.AddCommand(command)
	}
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
