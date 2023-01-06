package component

import (
	"bunnyshell.com/cli/pkg/lib"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:     "components",
	Aliases: []string{"comp"},

	Short: "Bunnyshell Components",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		lib.LoadViperConfigIntoContext()
	},
}

func init() {
	lib.CLIContext.RequireTokenOnCommand(mainCmd)
}

func GetMainCommand() *cobra.Command {
	return mainCmd
}
