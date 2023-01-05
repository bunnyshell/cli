package variable

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

var mainCmd = &cobra.Command{
	Use:     "variables",
	Aliases: []string{"var"},

	Short: "Bunnyshell Environment Variables",
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
