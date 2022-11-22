package port_forward

import (
	"github.com/spf13/cobra"

	"bunnyshell.com/cli/pkg/lib"
)

var mainCmd = &cobra.Command{
	Use:   "port-forward",
	Short: "Port Forward",
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
